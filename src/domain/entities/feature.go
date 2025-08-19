package entities

import (
	"errors"
	"strings"
	"time"
)

type Feature struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (f *Feature) GetID() string {
	return f.ID
}

// SetID - устанавливает ID
func (f *Feature) SetID(id string) {
	f.ID = id
}

// GetCreatedAt - возвращает время создания
func (f *Feature) GetCreatedAt() time.Time {
	return f.CreatedAt
}

// SetCreatedAt - устанавливает время создания
func (f *Feature) SetCreatedAt(t time.Time) {
	f.CreatedAt = t
}

// GetUpdatedAt - возвращает время обновления
func (f *Feature) GetUpdatedAt() time.Time {
	return f.UpdatedAt
}

// SetUpdatedAt - устанавливает время обновления
func (f *Feature) SetUpdatedAt(t time.Time) {
	f.UpdatedAt = t
}

// Бизнес-логика

// NewFeature - создает новую Feature сущность
func NewFeature(name string) (*Feature, error) {
	faq := &Feature{
		Name:      strings.TrimSpace(name),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := faq.Validate(); err != nil {
		return nil, err
	}

	return faq, nil
}

func (f *Feature) SetName(name string) error {
	f.Name = strings.TrimSpace(name)

	return f.ValidateName()
}

func (f *Feature) ValidateName() error {
	if f.Name == "" {
		return errors.New("name cannot be empty")
	}

	if len(f.Name) < 10 {
		return errors.New("question must be at least 10 characters long")
	}

	if len(f.Name) > 500 {
		return errors.New("question cannot exceed 500 characters")
	}

	return nil
}

func (f *Feature) Validate() error {
	if err := f.ValidateName(); err != nil {
		return err
	}

	return nil
}
