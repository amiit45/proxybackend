package handlers

import (
	"net/http"

	"myproject/internal/models"
	"myproject/internal/services"

	"github.com/gin-gonic/gin"
)

// SendRequest sends a connection request to another user
func SendRequest(c *gin.Context) {
	senderId := c.GetString("userId")
	var req models.ConnectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create connection request
	requestId, err := services.RequestService.CreateConnectionRequest(senderId, req.ReceiverId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request_id": requestId})
}

// GetRequests gets all pending connection requests for a user
func GetRequests(c *gin.Context) {
	userId := c.GetString("userId")

	// Get pending requests for user
	requests, err := services.RequestService.GetPendingRequests(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get requests"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

// RespondToRequest handles a user's response to a connection request
func RespondToRequest(c *gin.Context) {
	userId := c.GetString("userId")
	requestId := c.Param("id")

	var response models.RequestResponse
	if err := c.ShouldBindJSON(&response); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update request status
	err := services.RequestService.UpdateRequestStatus(requestId, userId, response.Accept)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
