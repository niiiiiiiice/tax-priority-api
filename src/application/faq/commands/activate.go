package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/events"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type ActivateFAQCommand struct {
	ID string `json:"id" validate:"required"`
}

type ActivateFAQCommandHandler struct {
	repo                repositories.FAQRepository
	notificationService events.NotificationService
}

func NewActivateFAQCommandHandler(repo repositories.FAQRepository, notificationService events.NotificationService) *ActivateFAQCommandHandler {
	return &ActivateFAQCommandHandler{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (h *ActivateFAQCommandHandler) HandleActivateFAQ(ctx context.Context, cmd ActivateFAQCommand) (*dtos.CommandResult, error) {

	faq, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to find FAQ: %v", err),
		}, err
	}

	faq.Activate()

	if err := h.repo.Update(ctx, faq); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to activate FAQ: %v", err),
		}, err
	}

	// Отправляем уведомление об активации FAQ
	if h.notificationService != nil {
		h.notificationService.NotifyFAQActivated(ctx, faq)
	}

	return &dtos.CommandResult{
		ID:        faq.ID,
		Success:   true,
		Message:   "FAQ activated successfully",
		UpdatedAt: faq.UpdatedAt,
	}, nil
}
