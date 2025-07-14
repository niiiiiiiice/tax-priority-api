package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"time"
)

type GetFAQCountQuery struct {
	Filters map[string]interface{} `json:"filters"`
}

type GetFAQCountQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func NewGetFAQCountQueryHandler(repo repositories.FAQRepository) *GetFAQCountQueryHandler {
	return &GetFAQCountQueryHandler{faqRepo: repo}
}

func (h *GetFAQCountQueryHandler) HandleGetFAQCount(ctx context.Context, query GetFAQCountQuery) (*dtos.QueryResult, error) {
	count, err := h.faqRepo.Count(ctx, query.Filters)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to count FAQs: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		Count:     count,
		Success:   true,
		Message:   "FAQ count retrieved successfully",
		Timestamp: time.Now(),
	}, nil
}
