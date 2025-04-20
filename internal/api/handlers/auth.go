package handlers

import (
	"myproject/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUser handles user registration
func RegisterUser(c *gin.Context) {
	var input struct {
		DeviceId string `json:"device_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, token, err := services.AuthService.RegisterUser(input.DeviceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userId": userId, "token": token})
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var input struct {
		DeviceId string `json:"device_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.AuthService.LoginUser(input.DeviceId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
