package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/domain/repositories"
)

type DeleteFAQCommand struct {
	ID string `json:"id" validate:"required"`
}

type DeleteFAQCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func (h *DeleteFAQCommandHandler) HandleDeleteFAQ(ctx context.Context, cmd DeleteFAQCommand) (*dtos.CommandResult, error) {

	exists, err := h.faqRepo.Exists(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to check if FAQ exists: %v", err),
		}, err
	}

	if !exists {
		return &dtos.CommandResult{
			Success: false,
			Error:   "FAQ not found",
		}, fmt.Errorf("FAQ with ID %s not found", cmd.ID)
	}

	if err := h.faqRepo.Delete(ctx, cmd.ID); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to delete FAQ: %v", err),
		}, err
	}

	return &dtos.CommandResult{
		ID:      cmd.ID,
		Success: true,
		Message: "FAQ deleted successfully",
	}, nil
}
