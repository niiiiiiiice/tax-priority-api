package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"

	"github.com/google/uuid"
)

type CreateFAQBatchCommand struct {
	FAQs []CreateFAQCommand `json:"faqs" validate:"required,min=1,max=100"`
}

type CreateFAQBatchCommandHandler struct {
	faqRepo repositories.FAQRepository
}

func NewCreateFAQBatchCommandHandler(repo repositories.FAQRepository) *CreateFAQBatchCommandHandler {
	return &CreateFAQBatchCommandHandler{faqRepo: repo}
}

func (h *CreateFAQBatchCommandHandler) HandleCreateFAQBatch(ctx context.Context, cmd CreateFAQBatchCommand) (*dtos.BatchCommandResult, error) {
	results := make([]dtos.CommandResult, 0, len(cmd.FAQs))
	faqs := make([]*entities.FAQ, 0, len(cmd.FAQs))
	successCount := 0
	failureCount := 0
	errors := make([]string, 0)

	for _, faqCmd := range cmd.FAQs {
		faq, err := entities.NewFAQ(faqCmd.Question, faqCmd.Answer, faqCmd.Category)
		if err != nil {
			failureCount++
			errors = append(errors, fmt.Sprintf("failed to create FAQ entity: %v", err))
			results = append(results, dtos.CommandResult{
				Success: false,
				Error:   fmt.Sprintf("failed to create FAQ entity: %v", err),
			})
			continue
		}

		faq.SetID(uuid.New().String())
		if faqCmd.Priority > 0 {
			if err := faq.SetPriority(faqCmd.Priority); err != nil {
				failureCount++
				errors = append(errors, fmt.Sprintf("failed to set priority: %v", err))
				results = append(results, dtos.CommandResult{
					Success: false,
					Error:   fmt.Sprintf("failed to set priority: %v", err),
				})
				continue
			}
		}

		faqs = append(faqs, faq)
		successCount++
		results = append(results, dtos.CommandResult{
			ID:        faq.ID,
			Success:   true,
			Message:   "FAQ created successfully",
			CreatedAt: faq.CreatedAt,
			UpdatedAt: faq.UpdatedAt,
		})
	}

	if len(faqs) > 0 {
		_, err := h.faqRepo.CreateBatch(ctx, faqs)
		if err != nil {
			for i := range results {
				if results[i].Success {
					results[i].Success = false
					results[i].Error = fmt.Sprintf("failed to save FAQ batch: %v", err)
					failureCount++
					successCount--
				}
			}
			errors = append(errors, fmt.Sprintf("failed to save FAQ batch: %v", err))
		}
	}

	return &dtos.BatchCommandResult{
		SuccessCount: successCount,
		FailureCount: failureCount,
		Results:      results,
		Errors:       errors,
	}, nil
}
