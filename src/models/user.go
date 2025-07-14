package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null;default:'user'"`
	Status    string    `json:"status" gorm:"not null;default:'active'"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate hook для установки UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Валидация роли
func (u *User) IsValidRole() bool {
	validRoles := []string{"admin", "user", "moderator"}
	for _, role := range validRoles {
		if u.Role == role {
			return true
		}
	}
	return false
}

// Валидация статуса
func (u *User) IsValidStatus() bool {
	validStatuses := []string{"active", "inactive", "banned"}
	for _, status := range validStatuses {
		if u.Status == status {
			return true
		}
	}
	return false
}
