package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/events"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type DeactivateFAQCommand struct {
	ID string `json:"id" validate:"required"`
}

type DeactivateFAQCommandHandler struct {
	repo                repositories.FAQRepository
	notificationService events.NotificationService
}

func NewDeactivateFAQCommandHandler(repo repositories.FAQRepository, notificationService events.NotificationService) *DeactivateFAQCommandHandler {
	return &DeactivateFAQCommandHandler{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (h *DeactivateFAQCommandHandler) HandleDeactivateFAQ(ctx context.Context, cmd DeactivateFAQCommand) (*dtos.CommandResult, error) {

	faq, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to find FAQ: %v", err),
		}, err
	}

	faq.Deactivate()

	if err := h.repo.Update(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to deactivate FAQ: %v", err),
		}, err
	}

	// Отправляем уведомление о деактивации FAQ
	if h.notificationService != nil {
		h.notificationService.NotifyFAQDeactivated(ctx, faq)
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ deactivated successfully",
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
