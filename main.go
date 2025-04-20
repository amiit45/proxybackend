package main

import (
	"log"
	"myproject/internal/api"
)

func main() {
	// Initialize router
	router := api.SetupRouter()

	// Start server
	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
