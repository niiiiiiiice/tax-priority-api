package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
)

type BulkDeleteFAQCommand struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

type BulkDeleteFAQCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func NewBulkDeleteFAQCommandHandler(repo repositories.FAQRepository) *BulkDeleteFAQCommandHandler {
	return &BulkDeleteFAQCommandHandler{faqRepo: repo}
}

func (h *BulkDeleteFAQCommandHandler) HandleBulkDeleteFAQ(ctx context.Context, cmd BulkDeleteFAQCommand) (*dtos.BatchCommandResult, error) {
	result, err := h.faqRepo.DeleteBatch(ctx, cmd.IDs)
	if err != nil {
		return &dtos.BatchCommandResult{
			SuccessCount: 0,
			FailureCount: len(cmd.IDs),
			Errors:       []string{fmt.Sprintf("failed to delete FAQs: %v", err)},
		}, err
	}

	return &dtos.BatchCommandResult{
		SuccessCount: result.SuccessCount,
		FailureCount: result.FailureCount,
		Errors:       make([]string, len(result.Errors)),
	}, nil
}
