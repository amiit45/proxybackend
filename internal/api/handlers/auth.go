package handlers

import (
	"log"
	"myproject/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var input struct {
		// DeviceId string `json:"device_id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("RegisterUser: invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("RegisterUser: input username=%s", input.Username)

	userId, token, err := services.AuthService.RegisterUser(input.Username, input.Password)
	if err != nil {
		log.Printf("RegisterUser: failed to register user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	log.Printf("RegisterUser: success userId=%s token=%s", userId, token)

	c.JSON(http.StatusOK, gin.H{"userId": userId, "token": token})
}

func LoginUser(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("LoginUser: invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("LoginUser: input username=%s", input.Username)

	userId, token, err := services.AuthService.LoginUser(input.Username, input.Password)
	if err != nil {
		log.Printf("LoginUser: invalid credentials: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	log.Printf("LoginUser: success userId=%s token=%s", userId, token)

	c.JSON(http.StatusOK, gin.H{
		"userId": userId,
		"token":  token,
	})
}
