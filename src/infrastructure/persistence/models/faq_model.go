package models

import (
	"time"

	"tax-priority-api/src/domain/entities"

	"gorm.io/gorm"
)

// FAQModel GORM модель для FAQ
type FAQModel struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)"`
	Question  string         `gorm:"type:text;not null"`
	Answer    string         `gorm:"type:text;not null"`
	Category  string         `gorm:"type:varchar(100);not null;index"`
	IsActive  bool           `gorm:"default:true;index"`
	Priority  int            `gorm:"default:0;index"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName возвращает имя таблицы для GORM
func (*FAQModel) TableName() string {
	return "faqs"
}

// ToEntity преобразует GORM модель в domain entity
func (m *FAQModel) ToEntity() *entities.FAQ {
	return &entities.FAQ{
		ID:        m.ID,
		Question:  m.Question,
		Answer:    m.Answer,
		Category:  m.Category,
		IsActive:  m.IsActive,
		Priority:  m.Priority,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity заполняет GORM модель из domain entity
func (m *FAQModel) FromEntity(faq *entities.FAQ) {
	m.ID = faq.ID
	m.Question = faq.Question
	m.Answer = faq.Answer
	m.Category = faq.Category
	m.IsActive = faq.IsActive
	m.Priority = faq.Priority
	m.CreatedAt = faq.CreatedAt
	m.UpdatedAt = faq.UpdatedAt
}

// NewFAQModelFromEntity создает новую GORM модель из domain entity
func NewFAQModelFromEntity(faq *entities.FAQ) *FAQModel {
	model := &FAQModel{}
	model.FromEntity(faq)
	return model
}
