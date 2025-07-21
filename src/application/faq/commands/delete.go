package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/events"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type DeleteFAQCommand struct {
	ID string `json:"id" validate:"required"`
}

type DeleteFAQCommandHandler struct {
	repo                repositories.FAQRepository
	notificationService events.NotificationService
}

func NewDeleteFAQCommandHandler(repo repositories.FAQRepository, notificationService events.NotificationService) *DeleteFAQCommandHandler {
	return &DeleteFAQCommandHandler{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (h *DeleteFAQCommandHandler) HandleDeleteFAQ(ctx context.Context, cmd DeleteFAQCommand) (*dtos.CommandResult, error) {

	exists, err := h.repo.Exists(ctx, cmd.ID)
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

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to delete FAQ: %v", err),
		}, err
	}

	// Отправляем уведомление об удалении FAQ
	if h.notificationService != nil {
		h.notificationService.NotifyFAQDeleted(ctx, cmd.ID)
	}

	return &dtos.CommandResult{
		ID:      cmd.ID,
		Success: true,
		Message: "FAQ deleted successfully",
	}, nil
}
