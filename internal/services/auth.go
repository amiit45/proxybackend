package services

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"myproject/internal/db"
	"myproject/internal/models"

	"github.com/dgrijalva/jwt-go"
)

// AuthServiceInterface defines the interface for auth operations
type AuthServiceInterface interface {
	RegisterUser(username string, password string) (string, string, error)
	LoginUser(username string, password string) (string, string, error)
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
func (s *authService) RegisterUser(username string, password string) (string, string, error) {
	log.Printf("RegisterUser called: username=%s", username)

	// Check if user already exists
	var existingUser models.User
	err := db.DB.Where("username = ?", username).First(&existingUser).Error
	if err == nil {
		return "", "", errors.New("username already exists")
	} else if err != gorm.ErrRecordNotFound {
		return "", "", err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	// Create user
	user := models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return "", "", err
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return user.ID, token, nil
}

// LoginUser authenticates a user
func (s *authService) LoginUser(username string, password string) (string, string, error) {
	log.Printf("LoginUser called: username=%s", username)

	if username == "" || password == "" {
		return "", "", errors.New("username and password required")
	}

	// Find user
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", "", errors.New("invalid username or password")
		}
		return "", "", err
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return user.ID, token, nil
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
// func generateUniqueId() string {
// 	// Implementation to generate a unique ID
// 	return "user-id"
// }
