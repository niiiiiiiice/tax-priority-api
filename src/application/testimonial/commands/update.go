package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"time"
)

type UpdateTestimonialCommandHandler struct {
	testimonialRepo repositories.TestimonialRepository
}

func NewUpdateTestimonialCommandHandler(repo repositories.TestimonialRepository) *UpdateTestimonialCommandHandler {
	return &UpdateTestimonialCommandHandler{
		testimonialRepo: repo,
	}
}

func (h *UpdateTestimonialCommandHandler) Handle(ctx context.Context, cmd dtos.UpdateTestimonialCommand) (*dtos.CommandResult, error) {
	// Получаем существующий отзыв
	testimonial, err := h.testimonialRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("testimonial not found: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	// Обновляем поля если они предоставлены
	if cmd.Content != "" {
		testimonial.UpdateContent(cmd.Content)
	}
	if cmd.Author != "" {
		testimonial.Author = cmd.Author
	}
	if cmd.AuthorEmail != "" {
		testimonial.AuthorEmail = cmd.AuthorEmail
	}
	if cmd.Rating > 0 {
		testimonial.UpdateRating(cmd.Rating)
	}
	if cmd.Company != "" {
		testimonial.Company = cmd.Company
	}
	if cmd.Position != "" {
		testimonial.Position = cmd.Position
	}

	// Сохраняем изменения
	if err := h.testimonialRepo.Update(ctx, testimonial); err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to update testimonial: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.CommandResult{
		Success:   true,
		Message:   "Testimonial updated successfully",
		Data:      testimonial,
		Timestamp: time.Now(),
	}, nil
}
