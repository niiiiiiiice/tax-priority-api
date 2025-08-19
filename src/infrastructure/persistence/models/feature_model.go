package models

import (
	"tax-priority-api/src/domain/entities"
	"time"

	"gorm.io/gorm"
)

// FeatureModel GORM модель для FAQ
type FeatureModel struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)"`
	Name      string         `gorm:"type:text;not null"`
	IsActive  bool           `gorm:"default:true;index"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName возвращает имя таблицы для GORM
func (*FeatureModel) TableName() string {
	return "features"
}

// ToEntity преобразует GORM модель в domain entity
func (m *FeatureModel) ToEntity() *entities.Feature {
	return &entities.Feature{
		ID:        m.ID,
		Name:      m.Name,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity заполняет GORM модель из domain entity
func (m *FeatureModel) FromEntity(feature *entities.Feature) {
	m.ID = feature.ID
	m.Name = feature.Name
	m.IsActive = feature.IsActive
	m.CreatedAt = feature.CreatedAt
	m.UpdatedAt = feature.UpdatedAt
}

// NewFeatureModelFromEntity создает новую GORM модель из domain entity
func NewFeatureModelFromEntity(feature *entities.Feature) *FeatureModel {
	model := &FeatureModel{}
	model.FromEntity(feature)
	return model
}
