package entities

import (
	"time"
)

type Testimonial struct {
	ID          string     `json:"id"`
	Content     string     `json:"content" validate:"required,min=10,max=1000"`
	Author      string     `json:"author" validate:"required,min=2,max=100"`
	AuthorEmail string     `json:"authorEmail" validate:"required,email"`
	Rating      int        `json:"rating" validate:"required,min=1,max=5"`
	FilePath    string     `json:"filePath,omitempty"`
	FileName    string     `json:"fileName,omitempty"`
	FileType    string     `json:"fileType,omitempty"`
	FileSize    int64      `json:"fileSize,omitempty"`
	IsApproved  bool       `json:"isApproved" default:"false"`
	IsActive    bool       `json:"isActive" default:"true"`
	ApprovedAt  *time.Time `json:"approvedAt,omitempty"`
	ApprovedBy  string     `json:"approvedBy,omitempty"`
	Company     string     `json:"company,omitempty"`
	Position    string     `json:"position,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// Реализация интерфейса Entity
func (t *Testimonial) GetID() string {
	return t.ID
}

func (t *Testimonial) SetID(id string) {
	t.ID = id
}

func (t *Testimonial) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Testimonial) SetCreatedAt(time time.Time) {
	t.CreatedAt = time
}

func (t *Testimonial) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}

func (t *Testimonial) SetUpdatedAt(time time.Time) {
	t.UpdatedAt = time
}

func NewTestimonial(content, author, authorEmail string, rating int) *Testimonial {
	now := time.Now()
	return &Testimonial{
		Content:     content,
		Author:      author,
		AuthorEmail: authorEmail,
		Rating:      rating,
		IsApproved:  false,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (t *Testimonial) Approve(approvedBy string) {
	now := time.Now()
	t.IsApproved = true
	t.ApprovedAt = &now
	t.ApprovedBy = approvedBy
	t.UpdatedAt = now
}

func (t *Testimonial) Deactivate() {
	t.IsActive = false
	t.UpdatedAt = time.Now()
}

func (t *Testimonial) Activate() {
	t.IsActive = true
	t.UpdatedAt = time.Now()
}

func (t *Testimonial) SetFile(filePath, fileName, fileType string, fileSize int64) {
	t.FilePath = filePath
	t.FileName = fileName
	t.FileType = fileType
	t.FileSize = fileSize
	t.UpdatedAt = time.Now()
}

func (t *Testimonial) UpdateContent(content string) {
	t.Content = content
	t.UpdatedAt = time.Now()
}

func (t *Testimonial) UpdateRating(rating int) {
	t.Rating = rating
	t.UpdatedAt = time.Now()
}
