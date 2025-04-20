package services

import (
	"math"

	"myproject/internal/models"
)

// LocationServiceInterface defines the interface for location operations
type LocationServiceInterface interface {
	UpdateUserLocation(userId string, lat, lng float64) error
	GetUserLocation(userId string) (float64, float64, error)
	FindNearbyUsers(userId string, lat, lng, radiusKm float64) ([]models.NearbyUser, error)
}

// locationService implements LocationServiceInterface
type locationService struct{}

// LocationService is the singleton instance
var LocationService LocationServiceInterface = &locationService{}

// UpdateUserLocation updates a user's location
func (s *locationService) UpdateUserLocation(userId string, lat, lng float64) error {
	// Implementation to update user location in database
	return nil
}

// GetUserLocation gets a user's current location
func (s *locationService) GetUserLocation(userId string) (float64, float64, error) {
	// Implementation to get user location from database
	return 0.0, 0.0, nil
}

// FindNearbyUsers finds users within a given radius
func (s *locationService) FindNearbyUsers(userId string, lat, lng, radiusKm float64) ([]models.NearbyUser, error) {
	// Implementation using Haversine formula to find nearby users
	return []models.NearbyUser{}, nil
}

// HaversineDistance calculates the distance between two points in kilometers
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Earth's radius in kilometers
	const R = 6371.0

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Differences
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}
