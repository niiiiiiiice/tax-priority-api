package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestUser_IsValidRole(t *testing.T) {
	tests := []struct {
		name string
		role string
		want bool
	}{
		{"Valid admin role", "admin", true},
		{"Valid user role", "user", true},
		{"Valid moderator role", "moderator", true},
		{"Invalid role", "invalid", false},
		{"Empty role", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Role: tt.role}
			if got := u.IsValidRole(); got != tt.want {
				t.Errorf("User.IsValidRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsValidStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   bool
	}{
		{"Valid active status", "active", true},
		{"Valid inactive status", "inactive", true},
		{"Valid banned status", "banned", true},
		{"Invalid status", "invalid", false},
		{"Empty status", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{Status: tt.status}
			if got := u.IsValidStatus(); got != tt.want {
				t.Errorf("User.IsValidStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_BeforeCreate(t *testing.T) {
	user := &User{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Проверяем, что ID не установлен
	if user.ID != uuid.Nil {
		t.Error("User ID should be nil before BeforeCreate")
	}

	// Вызываем BeforeCreate (в реальном приложении это делает GORM)
	err := user.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	// Проверяем, что ID установлен
	if user.ID == uuid.Nil {
		t.Error("User ID should be set after BeforeCreate")
	}
}
