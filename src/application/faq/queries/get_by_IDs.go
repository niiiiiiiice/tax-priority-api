package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"time"
)

type GetFAQsByIDsQuery struct {
	IDs []string `json:"ids" validate:"required,min=1"`
}

type GetFAQsByIDsQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func NewGetFAQsByIDsQueryHandler(repo repositories.FAQRepository) *GetFAQsByIDsQueryHandler {
	return &GetFAQsByIDsQueryHandler{faqRepo: repo}
}

func (h *GetFAQsByIDsQueryHandler) HandleGetFAQsByIDs(ctx context.Context, query GetFAQsByIDsQuery) (*dtos.QueryResult, error) {
	faqs, err := h.faqRepo.FindByIDs(ctx, query.IDs)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to find FAQs by IDs: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		FAQs:      faqs,
		Success:   true,
		Message:   "FAQs retrieved successfully",
		Timestamp: time.Now(),
	}, nil
}
