package repository

import (
	"context"
	"errors"
	"fmt"

	"my-flat-login/internal/model" // Replace with your actual module path

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByFirebaseID(ctx context.Context, firebaseID string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByFirebaseID(ctx context.Context, firebaseID string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("firebase_id = ?", firebaseID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("error finding user by firebase ID: %w", result.Error)
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("error creating user: %w", result.Error)
	}
	return nil
}
