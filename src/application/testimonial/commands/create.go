package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"tax-priority-api/src/domain/entities"
	"time"

	"github.com/google/uuid"
)

type CreateTestimonialCommandHandler struct {
	testimonialRepo repositories.TestimonialRepository
}

func NewCreateTestimonialCommandHandler(repo repositories.TestimonialRepository) *CreateTestimonialCommandHandler {
	return &CreateTestimonialCommandHandler{
		testimonialRepo: repo,
	}
}

func (h *CreateTestimonialCommandHandler) Handle(ctx context.Context, cmd dtos.CreateTestimonialCommand) (*dtos.CommandResult, error) {
	// Создаем новый отзыв
	testimonial := entities.NewTestimonial(
		cmd.Content,
		cmd.Author,
		cmd.AuthorEmail,
		cmd.Rating,
	)

	// Генерируем ID
	testimonial.SetID(uuid.New().String())

	// Устанавливаем дополнительные поля
	if cmd.Company != "" {
		testimonial.Company = cmd.Company
	}
	if cmd.Position != "" {
		testimonial.Position = cmd.Position
	}

	// Сохраняем в репозиторий
	if err := h.testimonialRepo.Create(ctx, testimonial); err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to create testimonial: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.CommandResult{
		Success:   true,
		Message:   "Testimonial created successfully",
		Data:      testimonial,
		Timestamp: time.Now(),
	}, nil
}
