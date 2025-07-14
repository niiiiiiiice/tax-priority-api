package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestProduct_IsValidCurrency(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		want     bool
	}{
		{"Valid USD currency", "USD", true},
		{"Valid EUR currency", "EUR", true},
		{"Valid GBP currency", "GBP", true},
		{"Invalid currency", "RUB", false},
		{"Empty currency", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Product{Currency: tt.currency}
			if got := p.IsValidCurrency(); got != tt.want {
				t.Errorf("Product.IsValidCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProduct_BeforeCreate(t *testing.T) {
	product := &Product{
		Name:     "Test Product",
		Price:    29.99,
		Currency: "USD",
		Category: "test",
	}

	// Проверяем, что ID не установлен
	if product.ID != uuid.Nil {
		t.Error("Product ID should be nil before BeforeCreate")
	}

	// Вызываем BeforeCreate (в реальном приложении это делает GORM)
	err := product.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate() error = %v", err)
	}

	// Проверяем, что ID установлен
	if product.ID == uuid.Nil {
		t.Error("Product ID should be set after BeforeCreate")
	}
}

func TestProduct_BeforeSave(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		want     bool
	}{
		{"Positive quantity", 10, true},
		{"Zero quantity", 0, false},
		{"Negative quantity", -5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Product{
				Name:     "Test Product",
				Quantity: tt.quantity,
			}

			// Вызываем BeforeSave (в реальном приложении это делает GORM)
			err := p.BeforeSave(nil)
			if err != nil {
				t.Errorf("BeforeSave() error = %v", err)
			}

			if p.InStock != tt.want {
				t.Errorf("Product.InStock = %v, want %v", p.InStock, tt.want)
			}
		})
	}
}
