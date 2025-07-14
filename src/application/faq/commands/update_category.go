package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type UpdateFAQCategoryCommand struct {
	ID       string `json:"id" validate:"required"`
	Category string `json:"category" validate:"required,max=100"`
}

type UpdateFAQCategoryCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func NewUpdateFAQCategoryCommandHandler(repo repositories.FAQRepository) *UpdateFAQCategoryCommandHandler {
	return &UpdateFAQCategoryCommandHandler{faqRepo: repo}
}

func (h *UpdateFAQCategoryCommandHandler) HandleUpdateFAQCategory(ctx context.Context, cmd UpdateFAQCategoryCommand) (*dtos.CommandResult, error) {

	faq, err := h.faqRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to find FAQ: %v", err),
		}, err
	}

	if err := faq.UpdateCategory(cmd.Category); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to update category: %v", err),
		}, err
	}

	if err := h.faqRepo.Update(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to update FAQ category: %v", err),
		}, err
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ category updated successfully",
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
