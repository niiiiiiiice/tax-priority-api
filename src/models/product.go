package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Currency    string    `json:"currency" gorm:"not null;default:'USD'"`
	Category    string    `json:"category" gorm:"not null"`
	Tags        []string  `json:"tags" gorm:"type:text[]"`
	InStock     bool      `json:"inStock" gorm:"default:true"`
	Quantity    int       `json:"quantity" gorm:"default:0"`
	Images      []string  `json:"images" gorm:"type:text[]"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// BeforeCreate hook для установки UUID
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// BeforeSave hook для автоматического обновления inStock на основе quantity
func (p *Product) BeforeSave(tx *gorm.DB) error {
	p.InStock = p.Quantity > 0
	return nil
}

// Валидация валюты
func (p *Product) IsValidCurrency() bool {
	validCurrencies := []string{"USD", "EUR", "GBP"}
	for _, currency := range validCurrencies {
		if p.Currency == currency {
			return true
		}
	}
	return false
}
