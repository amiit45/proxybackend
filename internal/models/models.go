package models

import (
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// LocationUpdate represents a location update request
type LocationUpdate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// ConnectionRequest represents a connection request
type ConnectionRequest struct {
	ID         string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	SenderId   string         `gorm:"index;not null" json:"sender_id"`
	ReceiverId string         `gorm:"index;not null" json:"receiver_id"`
	Status     string         `gorm:"not null" json:"status"` // e.g., "pending", "accepted", "rejected"
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
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

// User represents a user in the database
type User struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Location represents a user's location in the database
type Location struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"index;not null" json:"user_id"`
	Latitude  float64        `gorm:"not null" json:"latitude"`
	Longitude float64        `gorm:"not null" json:"longitude"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
