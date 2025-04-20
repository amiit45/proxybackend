package handlers

import (
	"net/http"

	"myproject/internal/models"
	"myproject/internal/services"

	"github.com/gin-gonic/gin"
)

// UpdateLocation updates a user's location
func UpdateLocation(c *gin.Context) {
	userId := c.GetString("userId")
	var location models.LocationUpdate

	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.LocationService.UpdateUserLocation(userId, location.Latitude, location.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "location updated"})
}

// GetNearbyUsers finds users near the current user
func GetNearbyUsers(c *gin.Context) {
	userId := c.GetString("userId")

	// Get user's current location
	lat, lng, err := services.LocationService.GetUserLocation(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user location"})
		return
	}

	// Find nearby users
	nearbyUsers, err := services.LocationService.FindNearbyUsers(userId, lat, lng, 1.0) // 1km radius
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find nearby users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nearby_users": nearbyUsers})
}
