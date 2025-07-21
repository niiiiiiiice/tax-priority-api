package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"time"
)

type ApproveTestimonialCommandHandler struct {
	testimonialRepo repositories.TestimonialRepository
}

func NewApproveTestimonialCommandHandler(repo repositories.TestimonialRepository) *ApproveTestimonialCommandHandler {
	return &ApproveTestimonialCommandHandler{
		testimonialRepo: repo,
	}
}

func (h *ApproveTestimonialCommandHandler) Handle(ctx context.Context, cmd dtos.ApproveTestimonialCommand) (*dtos.CommandResult, error) {
	// Получаем отзыв
	testimonial, err := h.testimonialRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("testimonial not found: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	// Одобряем отзыв
	testimonial.Approve(cmd.ApprovedBy)

	// Сохраняем изменения
	if err := h.testimonialRepo.Update(ctx, testimonial); err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to approve testimonial: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.CommandResult{
		Success:   true,
		Message:   "Testimonial approved successfully",
		Data:      testimonial,
		Timestamp: time.Now(),
	}, nil
}
