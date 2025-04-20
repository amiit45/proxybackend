package handlers

import (
	"log"
	"net/http"

	"myproject/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// HandleChatConnection handles WebSocket connections for chat
func HandleChatConnection(c *gin.Context) {
	userId := c.GetString("userId")
	requestId := c.Param("requestId")

	// Verify user is part of this chat
	if !services.ChatService.IsUserInChat(userId, requestId) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized for this chat"})
		return
	}

	// Upgrade to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}

	// Handle the chat connection
	services.ChatService.HandleConnection(conn, userId, requestId)
}
