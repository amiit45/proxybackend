package services

import (
	"log"
	"myproject/internal/models"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed
	maxMessageSize = 512
)

// ChatServiceInterface defines the interface for chat operations
type ChatServiceInterface interface {
	IsUserInChat(userId, chatId string) bool
	HandleConnection(conn *websocket.Conn, userId, chatId string)
}

// chatService implements ChatServiceInterface
type chatService struct {
	// clients map to track all connected clients
	clients map[*models.Client]bool
	// rooms map to track clients by chat room
	rooms map[string]map[*models.Client]bool
	// register channel for new clients
	register chan *models.Client
	// unregister channel for disconnected clients
	unregister chan *models.Client
	// broadcast channel for messages
	broadcast chan Message
}

type Message struct {
	chatId string
	data   []byte
}

// ChatService is the singleton instance
var ChatService ChatServiceInterface = newChatService()

// newChatService creates a new chat service
func newChatService() *chatService {
	service := &chatService{
		clients:    make(map[*models.Client]bool),
		rooms:      make(map[string]map[*models.Client]bool),
		register:   make(chan *models.Client),
		unregister: make(chan *models.Client),
		broadcast:  make(chan Message),
	}

	// Start the background goroutine
	go service.run()

	return service
}

// IsUserInChat checks if a user is part of a chat
func (s *chatService) IsUserInChat(userId, chatId string) bool {
	// Implementation to check if user is part of chat
	return true
}

// HandleConnection handles a new WebSocket connection
func (s *chatService) HandleConnection(conn *websocket.Conn, userId, chatId string) {
	client := &models.Client{
		Conn:   conn,
		UserId: userId,
		ChatId: chatId,
		Send:   make(chan []byte, 256),
	}

	// Register client
	s.register <- client

	// Handle WebSocket communication
	go s.writePump(client)
	go s.readPump(client)
}

// run processes the chat service's channels
func (s *chatService) run() {
	for {
		select {
		case client := <-s.register:
			s.registerClient(client)
		case client := <-s.unregister:
			s.unregisterClient(client)
		case message := <-s.broadcast:
			s.broadcastToRoom(message)
		}
	}
}

// registerClient registers a new client
func (s *chatService) registerClient(client *models.Client) {
	s.clients[client] = true

	// Create room if it doesn't exist
	if _, exists := s.rooms[client.ChatId]; !exists {
		s.rooms[client.ChatId] = make(map[*models.Client]bool)
	}

	// Add client to room
	s.rooms[client.ChatId][client] = true
}

// unregisterClient removes a client
func (s *chatService) unregisterClient(client *models.Client) {
	if _, exists := s.clients[client]; exists {
		delete(s.clients, client)
		close(client.Send)

		// Remove from room
		if room, exists := s.rooms[client.ChatId]; exists {
			delete(room, client)

			// Remove room if empty
			if len(room) == 0 {
				delete(s.rooms, client.ChatId)
			}
		}
	}
}

// broadcastToRoom sends a message to all clients in a room
func (s *chatService) broadcastToRoom(msg Message) {
	if room, exists := s.rooms[msg.chatId]; exists {
		for client := range room {
			select {
			case client.Send <- msg.data:
			default:
				s.unregisterClient(client)
			}
		}
	}
}

// readPump reads messages from the WebSocket
func (s *chatService) readPump(client *models.Client) {
	defer func() {
		s.unregister <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Broadcast message to room
		s.broadcast <- Message{chatId: client.ChatId, data: message}
	}
}

// writePump writes messages to the WebSocket
func (s *chatService) writePump(client *models.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Channel was closed
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
