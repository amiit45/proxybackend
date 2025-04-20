package models

import "github.com/gorilla/websocket"

//import "golang.org/x/net/websocket"

// LocationUpdate represents a location update request
type LocationUpdate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ConnectionRequest represents a connection request
type ConnectionRequest struct {
	ReceiverId string `json:"receiver_id"`
}

// RequestResponse represents a response to a connection request
type RequestResponse struct {
	Accept bool `json:"accept"`
}

// NearbyUser represents a user found nearby
type NearbyUser struct {
	ID       string  `json:"id"`
	Distance float64 `json:"distance"` // in kilometers
}

// Client represents a connected WebSocket client
type Client struct {
	Conn   *websocket.Conn
	UserId string
	ChatId string
	Send   chan []byte
}
