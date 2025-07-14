package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type DeactivateFAQCommand struct {
	ID string `json:"id" validate:"required"`
}

type DeactivateFAQCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func NewDeactivateFAQCommandHandler(repo repositories.FAQRepository) *DeactivateFAQCommandHandler {
	return &DeactivateFAQCommandHandler{faqRepo: repo}
}

func (h *DeactivateFAQCommandHandler) HandleDeactivateFAQ(ctx context.Context, cmd DeactivateFAQCommand) (*dtos.CommandResult, error) {

	faq, err := h.faqRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to find FAQ: %v", err),
		}, err
	}

	faq.Deactivate()

	if err := h.faqRepo.Update(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to deactivate FAQ: %v", err),
		}, err
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ deactivated successfully",
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
