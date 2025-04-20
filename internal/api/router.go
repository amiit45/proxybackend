package api

import (
	"myproject/internal/api/handlers"
	"myproject/internal/api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures and returns the API router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Setup CORS middleware
	r.Use(middleware.CORS())

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterUser)
		auth.POST("/login", handlers.LoginUser)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.Auth())
	{
		api.POST("/location", handlers.UpdateLocation)
		api.GET("/nearby", handlers.GetNearbyUsers)
		api.POST("/request", handlers.SendRequest)
		api.GET("/requests", handlers.GetRequests)
		api.PUT("/request/:id", handlers.RespondToRequest)
		api.GET("/ws/chat/:requestId", handlers.HandleChatConnection)
	}

	return r
}
