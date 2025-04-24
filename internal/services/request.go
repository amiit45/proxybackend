package services

import (
	"time"

	"myproject/internal/db"
	"myproject/internal/models"

	"gorm.io/gorm"
)

// RequestServiceInterface defines the interface for request operations
type RequestServiceInterface interface {
	CreateConnectionRequest(senderId, receiverId string) (string, error)
	GetPendingRequests(userId string) ([]models.ConnectionRequest, error)
	UpdateRequestStatus(requestId, userId string, accept bool) error
}

// requestService implements RequestServiceInterface
type requestService struct{}

// RequestService is the singleton instance
var RequestService RequestServiceInterface = &requestService{}

// CreateConnectionRequest creates a new connection request
func (s *requestService) CreateConnectionRequest(senderId, receiverId string) (string, error) {
	request := models.ConnectionRequest{
		SenderId:   senderId,
		ReceiverId: receiverId,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.DB.Create(&request).Error; err != nil {
		return "", err
	}

	return request.ID, nil
}

// GetPendingRequests gets all pending requests for a user
func (s *requestService) GetPendingRequests(userId string) ([]models.ConnectionRequest, error) {
	var requests []models.ConnectionRequest
	if err := db.DB.Where("receiver_id = ? AND status = ?", userId, "pending").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// UpdateRequestStatus updates the status of a request
func (s *requestService) UpdateRequestStatus(requestId, userId string, accept bool) error {
	var request models.ConnectionRequest
	if err := db.DB.First(&request, "id = ? AND receiver_id = ?", requestId, userId).Error; err != nil {
		return err
	}

	if request.Status != "pending" {
		return gorm.ErrInvalidData
	}

	if accept {
		request.Status = "accepted"
	} else {
		request.Status = "rejected"
	}

	if err := db.DB.Save(&request).Error; err != nil {
		return err
	}

	if accept {
		chatId, err := ChatService.CreateChatSession(request.SenderId, request.ReceiverId)
		if err != nil {
			return err
		}

		// TODO: Save chatId association with users if needed
		_ = chatId
	}

	return nil
}
