package entity

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user entity with business rules
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;size:100" binding:"required"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:100" binding:"required,email"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Active    bool           `json:"active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// BeforeCreate hook to validate business rules
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Name == "" {
		return ErrInvalidUserName
	}
	if u.Email == "" {
		return ErrInvalidUserEmail
	}
	return nil
}

// IsValid checks if user entity is valid
func (u *User) IsValid() bool {
	return u.Name != "" && u.Email != ""
}

// Activate sets user as active
func (u *User) Activate() {
	u.Active = true
}

// Deactivate sets user as inactive
func (u *User) Deactivate() {
	u.Active = false
}
