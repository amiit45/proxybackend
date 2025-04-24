package services

import (
	"log"
	"math"

	"myproject/internal/db"
	"myproject/internal/models"

	"gorm.io/gorm"
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
	log.Printf("UpdateUserLocation called: userId=%s lat=%f lng=%f", userId, lat, lng)
	var location models.Location
	result := db.DB.First(&location, "user_id = ?", userId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new location record
			location = models.Location{
				UserID:    userId,
				Latitude:  lat,
				Longitude: lng,
			}
			if err := db.DB.Create(&location).Error; err != nil {
				return err
			}
		} else {
			return result.Error
		}
	} else {
		// Update existing location record
		location.Latitude = lat
		location.Longitude = lng
		if err := db.DB.Save(&location).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetUserLocation gets a user's current location
func (s *locationService) GetUserLocation(userId string) (float64, float64, error) {
	log.Printf("GetUserLocation called: userId=%s", userId)
	var location models.Location
	result := db.DB.First(&location, "user_id = ?", userId)
	if result.Error != nil {
		return 0.0, 0.0, result.Error
	}
	return location.Latitude, location.Longitude, nil
}

// FindNearbyUsers finds users within a given radius
func (s *locationService) FindNearbyUsers(userId string, lat, lng, radiusKm float64) ([]models.NearbyUser, error) {
	log.Printf("FindNearbyUsers called: userId=%s lat=%f lng=%f radiusKm=%f", userId, lat, lng, radiusKm)
	var nearbyUsers []models.NearbyUser
	query := `
	SELECT user_id, (
		6371 * acos(
			cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) +
			sin(radians(?)) * sin(radians(latitude))
		)
	) AS distance
	FROM locations
	WHERE user_id != ?
	HAVING distance <= ?
	ORDER BY distance;
	`
	rows, err := db.DB.Raw(query, lat, lng, lat, userId, radiusKm).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.NearbyUser
		if err := rows.Scan(&user.ID, &user.Distance); err != nil {
			return nil, err
		}
		nearbyUsers = append(nearbyUsers, user)
	}
	return nearbyUsers, nil
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
