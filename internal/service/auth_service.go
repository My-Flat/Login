package service

import (
	"context"
	"fmt"
	"time"

	"my-flat-login/internal/model" // Replace with your actual module path
	"my-flat-login/internal/repository"

	"firebase.google.com/go/auth" // Use this import path

	"github.com/golang-jwt/jwt/v4"
)

// AuthService interface defines the methods for authentication
type AuthService interface {
	Login(ctx context.Context, idToken string) (*model.User, string, error)
}

type authService struct {
	firebaseAuth   *auth.Client
	userRepository repository.UserRepository
}

// NewAuthService creates a new instance of authService
func NewAuthService(firebaseAuth *auth.Client, userRepository repository.UserRepository) *authService {
	return &authService{
		firebaseAuth:   firebaseAuth,
		userRepository: userRepository,
	}
}

// JWT claims struct
type Claims struct {
	FirebaseID string `json:"firebase_id"`
	jwt.RegisteredClaims
}

// Login handles the login process using Firebase Authentication and generates a JWT
func (s *authService) Login(ctx context.Context, idToken string) (*model.User, string, error) {
	// 1. Verify Firebase ID token
	token, err := s.firebaseAuth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to verify ID token: %w", err)
	}

	firebaseID := token.UID

	// 2. Find or create user in the database
	user, err := s.userRepository.FindByFirebaseID(ctx, firebaseID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to find user by Firebase ID: %w", err)
	}

	if user == nil {
		// User not found, create a new user
		newUser := &model.User{
			FirebaseID: firebaseID,
			Email:      token.Claims["email"].(string),
		}
		err = s.userRepository.Create(ctx, newUser)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}
		user = newUser
	}

	// 3. Generate JWT
	// Replace with your actual JWT secret
	jwtSecret := []byte("your-jwt-secret")

	expirationTime := time.Now().Add(24 * time.Hour) // Example: 24 hour expiration
	claims := &Claims{
		FirebaseID: firebaseID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString,
		err := jwtToken.SignedString(jwtSecret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return user, tokenString, nil
}
