package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"time"
)

type BulkTestimonialCommandHandler struct {
	testimonialRepo repositories.TestimonialRepository
}

func NewBulkTestimonialCommandHandler(repo repositories.TestimonialRepository) *BulkTestimonialCommandHandler {
	return &BulkTestimonialCommandHandler{
		testimonialRepo: repo,
	}
}

func (h *BulkTestimonialCommandHandler) HandleBulkApprove(ctx context.Context, cmd dtos.BulkApproveTestimonialsCommand) (*dtos.CommandResult, error) {
	result, err := h.testimonialRepo.ApproveMany(ctx, cmd.IDs, cmd.ApprovedBy)
	if err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to approve testimonials: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.CommandResult{
		Success:   true,
		Message:   fmt.Sprintf("Successfully approved %d testimonials, %d failed", result.SuccessCount, result.FailureCount),
		Data:      result,
		Timestamp: time.Now(),
	}, nil
}

func (h *BulkTestimonialCommandHandler) HandleBulkDeactivate(ctx context.Context, cmd dtos.BulkDeactivateTestimonialsCommand) (*dtos.CommandResult, error) {
	result, err := h.testimonialRepo.DeactivateMany(ctx, cmd.IDs)
	if err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to deactivate testimonials: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.CommandResult{
		Success:   true,
		Message:   fmt.Sprintf("Successfully deactivated %d testimonials, %d failed", result.SuccessCount, result.FailureCount),
		Data:      result,
		Timestamp: time.Now(),
	}, nil
}

func (h *BulkTestimonialCommandHandler) HandleBulkActivate(ctx context.Context, cmd dtos.BulkActivateTestimonialsCommand) (*dtos.CommandResult, error) {
	result, err := h.testimonialRepo.ActivateMany(ctx, cmd.IDs)
	if err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to activate testimonials: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.CommandResult{
		Success:   true,
		Message:   fmt.Sprintf("Successfully activated %d testimonials, %d failed", result.SuccessCount, result.FailureCount),
		Data:      result,
		Timestamp: time.Now(),
	}, nil
}

func (h *BulkTestimonialCommandHandler) HandleBulkDelete(ctx context.Context, cmd dtos.BulkDeleteTestimonialsCommand) (*dtos.CommandResult, error) {
	result, err := h.testimonialRepo.DeleteMany(ctx, cmd.IDs)
	if err != nil {
		return &dtos.CommandResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to delete testimonials: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.CommandResult{
		Success:   true,
		Message:   fmt.Sprintf("Successfully deleted %d testimonials, %d failed", result.SuccessCount, result.FailureCount),
		Data:      result,
		Timestamp: time.Now(),
	}, nil
}
