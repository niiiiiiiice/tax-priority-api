package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/events"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"

	"github.com/google/uuid"
)

type CreateFAQCommand struct {
	Question string `json:"question" validate:"required,min=10,max=500"`
	Answer   string `json:"answer" validate:"required,min=10,max=2000"`
	Category string `json:"category" validate:"required,max=100"`
	Priority int    `json:"priority" validate:"min=0,max=100"`
}

type CreateFAQCommandHandler struct {
	repo                repositories.FAQRepository
	notificationService events.NotificationService
}

func NewCreateFAQCommandHandler(repo repositories.FAQRepository, notificationService events.NotificationService) *CreateFAQCommandHandler {
	return &CreateFAQCommandHandler{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (h *CreateFAQCommandHandler) HandleCreateFAQ(ctx context.Context, cmd CreateFAQCommand) (*dtos.CommandResult, error) {
	faq, err := entities.NewFAQ(
		cmd.Question,
		cmd.Answer,
		cmd.Category,
	)

	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create FAQ: %v", err),
		}, err
	}

	faq.SetID(uuid.New().String())

	if err = h.repo.Create(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create FAQ: %v", err),
		}, err
	}

	if h.notificationService != nil {
		h.notificationService.NotifyFAQCreated(ctx, faq)
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ created successfully",
		CreatedAt: faq.CreatedAt,
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
