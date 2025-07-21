package models

import (
	"tax-priority-api/src/domain/entities"
	"time"
)

type TestimonialModel struct {
	ID          string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	Author      string     `gorm:"type:varchar(100);not null" json:"author"`
	AuthorEmail string     `gorm:"type:varchar(255);not null" json:"authorEmail"`
	Rating      int        `gorm:"type:int;not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	FilePath    string     `gorm:"type:varchar(500)" json:"filePath"`
	FileName    string     `gorm:"type:varchar(255)" json:"fileName"`
	FileType    string     `gorm:"type:varchar(50)" json:"fileType"`
	FileSize    int64      `gorm:"type:bigint" json:"fileSize"`
	IsApproved  bool       `gorm:"type:boolean;default:false" json:"isApproved"`
	IsActive    bool       `gorm:"type:boolean;default:true" json:"isActive"`
	ApprovedAt  *time.Time `gorm:"type:timestamp" json:"approvedAt"`
	ApprovedBy  string     `gorm:"type:varchar(255)" json:"approvedBy"`
	Company     string     `gorm:"type:varchar(255)" json:"company"`
	Position    string     `gorm:"type:varchar(255)" json:"position"`
	CreatedAt   time.Time  `gorm:"type:timestamp;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"type:timestamp;autoUpdateTime" json:"updatedAt"`
}

func (TestimonialModel) TableName() string {
	return "testimonials"
}

func (m *TestimonialModel) ToEntity() *entities.Testimonial {
	return &entities.Testimonial{
		ID:          m.ID,
		Content:     m.Content,
		Author:      m.Author,
		AuthorEmail: m.AuthorEmail,
		Rating:      m.Rating,
		FilePath:    m.FilePath,
		FileName:    m.FileName,
		FileType:    m.FileType,
		FileSize:    m.FileSize,
		IsApproved:  m.IsApproved,
		IsActive:    m.IsActive,
		ApprovedAt:  m.ApprovedAt,
		ApprovedBy:  m.ApprovedBy,
		Company:     m.Company,
		Position:    m.Position,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func TestimonialFromEntity(entity *entities.Testimonial) *TestimonialModel {
	return &TestimonialModel{
		ID:          entity.ID,
		Content:     entity.Content,
		Author:      entity.Author,
		AuthorEmail: entity.AuthorEmail,
		Rating:      entity.Rating,
		FilePath:    entity.FilePath,
		FileName:    entity.FileName,
		FileType:    entity.FileType,
		FileSize:    entity.FileSize,
		IsApproved:  entity.IsApproved,
		IsActive:    entity.IsActive,
		ApprovedAt:  entity.ApprovedAt,
		ApprovedBy:  entity.ApprovedBy,
		Company:     entity.Company,
		Position:    entity.Position,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
