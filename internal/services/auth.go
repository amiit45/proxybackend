package services

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthServiceInterface defines the interface for auth operations
type AuthServiceInterface interface {
	RegisterUser(deviceId string) (string, string, error)
	LoginUser(deviceId string) (string, error)
	ValidateToken(token string) (string, error)
}

// authService implements AuthServiceInterface
type authService struct {
	secretKey []byte
}

// AuthService is the singleton instance
var AuthService AuthServiceInterface = &authService{
	secretKey: []byte("your-secret-key"), // In production, use environment variables
}

// RegisterUser registers a new user
func (s *authService) RegisterUser(deviceId string) (string, string, error) {
	// Implementation for user registration
	// In a real application, you would store this in a database
	userId := generateUniqueId()

	// Generate JWT token
	token, err := s.generateToken(userId)
	if err != nil {
		return "", "", err
	}

	return userId, token, nil
}

// LoginUser authenticates a user
func (s *authService) LoginUser(deviceId string) (string, error) {
	// Implementation for user login
	// In a real application, you would verify from a database
	userId := "user-id" // Replace with lookup based on deviceId

	// Generate JWT token
	token, err := s.generateToken(userId)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (s *authService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"].(string)
		return userId, nil
	}

	return "", errors.New("invalid token")
}

// generateToken creates a new JWT token
func (s *authService) generateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// generateUniqueId generates a unique ID for a user
func generateUniqueId() string {
	// Implementation to generate a unique ID
	return "user-id"
}
