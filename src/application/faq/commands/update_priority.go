package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type UpdateFAQPriorityCommand struct {
	ID       string `json:"id" validate:"required"`
	Priority int    `json:"priority" validate:"min=0,max=100"`
}

type UpdateFAQPriorityCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func NewUpdateFAQPriorityCommandHandler(repo repositories.FAQRepository) *UpdateFAQPriorityCommandHandler {
	return &UpdateFAQPriorityCommandHandler{faqRepo: repo}
}

func (h *UpdateFAQPriorityCommandHandler) HandleUpdateFAQPriority(ctx context.Context, cmd UpdateFAQPriorityCommand) (*dtos.CommandResult, error) {

	faq, err := h.faqRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to find FAQ: %v", err),
		}, err
	}

	if err := faq.SetPriority(cmd.Priority); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to set priority: %v", err),
		}, err
	}

	if err := h.faqRepo.Update(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to update FAQ priority: %v", err),
		}, err
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ priority updated successfully",
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
