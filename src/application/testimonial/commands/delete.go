package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"time"
)

type DeleteTestimonialCommandHandler struct {
	testimonialRepo repositories.TestimonialRepository
}

func NewDeleteTestimonialCommandHandler(repo repositories.TestimonialRepository) *DeleteTestimonialCommandHandler {
	return &DeleteTestimonialCommandHandler{
		testimonialRepo: repo,
	}
}

func (h *DeleteTestimonialCommandHandler) Handle(ctx context.Context, cmd dtos.DeleteTestimonialCommand) (*dtos.CommandResult, error) {
	testimonial, err := h.testimonialRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("testimonial not found: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	// Удаляем отзыв
	if err := h.testimonialRepo.Delete(ctx, cmd.ID); err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to delete testimonial: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.CommandResult{
		Success:   true,
		Message:   "Testimonial deleted successfully",
		Data:      testimonial,
		Timestamp: time.Now(),
	}, nil
}
