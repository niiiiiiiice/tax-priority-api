package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/domain/repositories"

	"github.com/google/uuid"
)

type CreateFAQCommand struct {
	Question string `json:"question" validate:"required,min=10,max=500"`
	Answer   string `json:"answer" validate:"required,min=10,max=2000"`
	Category string `json:"category" validate:"required,max=100"`
	Priority int    `json:"priority" validate:"min=0,max=100"`
}

type CreateFAQCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func NewCreateFAQCommandHandler(repo repositories.FAQRepository) *CreateFAQCommandHandler {
	return &CreateFAQCommandHandler{faqRepo: repo}
}

func (h *CreateFAQCommandHandler) HandleCreateFAQ(ctx context.Context, cmd CreateFAQCommand) (*dtos.CommandResult, error) {

	exists, err := h.faqRepo.ExistsByQuestion(ctx, cmd.Question)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to check if FAQ exists: %v", err),
		}, err
	}

	if exists {
		return &dtos.CommandResult{
			Success: false,
			Error:   "FAQ with this question already exists",
		}, fmt.Errorf("FAQ with question '%s' already exists", cmd.Question)
	}

	faq, err := entities.NewFAQ(cmd.Question, cmd.Answer, cmd.Category)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create FAQ entity: %v", err),
		}, err
	}

	faq.SetID(uuid.New().String())
	if cmd.Priority > 0 {
		if err := faq.SetPriority(cmd.Priority); err != nil {
			return &dtos.CommandResult{
				Success: false,
				Error:   fmt.Sprintf("failed to set priority: %v", err),
			}, err
		}
	}

	if err := h.faqRepo.Create(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create FAQ: %v", err),
		}, err
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ created successfully",
		CreatedAt: faq.CreatedAt,
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
