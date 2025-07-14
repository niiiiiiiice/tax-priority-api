package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/domain/repositories"
)

type UpdateFAQCommand struct {
	ID       string `json:"id" validate:"required"`
	Question string `json:"question" validate:"required,min=10,max=500"`
	Answer   string `json:"answer" validate:"required,min=10,max=2000"`
	Category string `json:"category" validate:"required,max=100"`
	Priority int    `json:"priority" validate:"min=0,max=100"`
}

type UpdateFAQCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func (h *UpdateFAQCommandHandler) HandleUpdateFAQ(ctx context.Context, cmd UpdateFAQCommand) (*dtos.CommandResult, error) {

	faq, err := h.faqRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to find FAQ: %v", err),
		}, err
	}

	if err := faq.UpdateQuestion(cmd.Question); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to update question: %v", err),
		}, err
	}

	if err := faq.UpdateAnswer(cmd.Answer); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to update answer: %v", err),
		}, err
	}

	if err := faq.UpdateCategory(cmd.Category); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to update category: %v", err),
		}, err
	}

	if err := faq.SetPriority(cmd.Priority); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to set priority: %v", err),
		}, err
	}

	// Сохраняем изменения
	if err := h.faqRepo.Update(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to update FAQ: %v", err),
		}, err
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ updated successfully",
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
