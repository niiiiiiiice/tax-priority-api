package dtos

import (
	"time"
)

// CreateTestimonialCommand для создания нового отзыва
type CreateTestimonialCommand struct {
	Content     string `json:"content" validate:"required,min=10,max=1000"`
	Author      string `json:"author" validate:"required,min=2,max=100"`
	AuthorEmail string `json:"authorEmail" validate:"required,email"`
	Rating      int    `json:"rating" validate:"required,min=1,max=5"`
	Company     string `json:"company,omitempty" validate:"max=255"`
	Position    string `json:"position,omitempty" validate:"max=255"`
}

// UpdateTestimonialCommand для обновления отзыва
type UpdateTestimonialCommand struct {
	ID          string `json:"id" validate:"required"`
	Content     string `json:"content,omitempty" validate:"min=10,max=1000"`
	Author      string `json:"author,omitempty" validate:"min=2,max=100"`
	AuthorEmail string `json:"authorEmail,omitempty" validate:"email"`
	Rating      int    `json:"rating,omitempty" validate:"min=1,max=5"`
	Company     string `json:"company,omitempty" validate:"max=255"`
	Position    string `json:"position,omitempty" validate:"max=255"`
}

// DeleteTestimonialCommand для удаления отзыва
type DeleteTestimonialCommand struct {
	ID string `json:"id" validate:"required"`
}

// ApproveTestimonialCommand для одобрения отзыва
type ApproveTestimonialCommand struct {
	ID         string `json:"id" validate:"required"`
	ApprovedBy string `json:"approvedBy" validate:"required"`
}

// DeactivateTestimonialCommand для деактивации отзыва
type DeactivateTestimonialCommand struct {
	ID string `json:"id" validate:"required"`
}

// ActivateTestimonialCommand для активации отзыва
type ActivateTestimonialCommand struct {
	ID string `json:"id" validate:"required"`
}

// BulkApproveTestimonialsCommand для массового одобрения
type BulkApproveTestimonialsCommand struct {
	IDs        []string `json:"ids" validate:"required,min=1"`
	ApprovedBy string   `json:"approvedBy" validate:"required"`
}

// BulkDeactivateTestimonialsCommand для массовой деактивации
type BulkDeactivateTestimonialsCommand struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

// BulkActivateTestimonialsCommand для массовой активации
type BulkActivateTestimonialsCommand struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

// BulkDeleteTestimonialsCommand для массового удаления
type BulkDeleteTestimonialsCommand struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

// UploadTestimonialFileCommand для загрузки файла к отзыву
type UploadTestimonialFileCommand struct {
	ID       string `json:"id" validate:"required"`
	FilePath string `json:"filePath" validate:"required"`
	FileName string `json:"fileName" validate:"required"`
	FileType string `json:"fileType" validate:"required"`
	FileSize int64  `json:"fileSize" validate:"required,min=1"`
}

// CommandResult общий результат выполнения команды
type CommandResult struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Error     string      `json:"error,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}
