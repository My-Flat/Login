package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct represents the user entity in the database
type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	FirebaseID string    `gorm:"uniqueIndex;not null"` // Use Firebase ID for authentication
	Email      string    `gorm:"uniqueIndex;not null"`
	UserName   string    `gorm:"uniqueIndex;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// BeforeCreate hook to generate UUID before saving
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return
}
