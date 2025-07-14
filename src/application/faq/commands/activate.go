package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type ActivateFAQCommand struct {
	ID string `json:"id" validate:"required"`
}

type ActivateFAQCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func NewActivateFAQCommandHandler(repo repositories.FAQRepository) *ActivateFAQCommandHandler {
	return &ActivateFAQCommandHandler{faqRepo: repo}
}

func (h *ActivateFAQCommandHandler) HandleActivateFAQ(ctx context.Context, cmd ActivateFAQCommand) (*dtos.CommandResult, error) {

	faq, err := h.faqRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to find FAQ: %v", err),
		}, err
	}

	faq.Activate()

	if err := h.faqRepo.Update(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to activate FAQ: %v", err),
		}, err
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ activated successfully",
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
