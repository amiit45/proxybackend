package main

import (
	"log"
	"myproject/internal/api"
	"myproject/internal/db"
)

func main() {
	// Initialize database
	db.Init()

	// Initialize router
	router := api.SetupRouter()

	// Start server
	log.Println("Starting server on 0.0.0.x:8080")
	if err := router.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
