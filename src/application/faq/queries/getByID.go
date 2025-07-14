package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/domain/repositories"
	"time"
)

type GetFAQByIDQuery struct {
	ID string `json:"id" validate:"required"`
}

type GetFAQByIDQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func (h *GetFAQByIDQueryHandler) HandleGetFAQByID(ctx context.Context, query GetFAQByIDQuery) (*dtos.QueryResult, error) {
	faq, err := h.faqRepo.FindByID(ctx, query.ID)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to find FAQ: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		FAQ:       faq,
		Success:   true,
		Message:   "FAQ retrieved successfully",
		Timestamp: time.Now(),
	}, nil
}
