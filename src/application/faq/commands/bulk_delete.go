package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/events"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type BulkDeleteFAQCommand struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

type BulkDeleteFAQCommandHandler struct {
	repo                repositories.FAQRepository
	notificationService events.NotificationService
}

func NewBulkDeleteFAQCommandHandler(repo repositories.FAQRepository, notificationService events.NotificationService) *BulkDeleteFAQCommandHandler {
	return &BulkDeleteFAQCommandHandler{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (h *BulkDeleteFAQCommandHandler) HandleBulkDeleteFAQ(ctx context.Context, cmd BulkDeleteFAQCommand) (*dtos.BatchCommandResult, error) {
	result, err := h.repo.DeleteBatch(ctx, cmd.IDs)
	if err != nil {
		return &dtos.BatchCommandResult{
			SuccessCount: 0,
			FailureCount: len(cmd.IDs),
			Errors:       []string{fmt.Sprintf("failed to delete FAQs: %v", err)},
		}, err
	}

	// Отправляем уведомление о массовом удалении FAQ
	if h.notificationService != nil {
		h.notificationService.NotifyFAQBatchDeleted(ctx, cmd.IDs)
	}

	return &dtos.BatchCommandResult{
		SuccessCount: result.SuccessCount,
		FailureCount: result.FailureCount,
		Errors:       make([]string, len(result.Errors)),
	}, nil
}
