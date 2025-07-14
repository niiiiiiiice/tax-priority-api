package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/shared/models"
	"time"
)

type GetActiveFAQsQuery struct {
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	Offset    int    `json:"offset" validate:"min=0"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder" validate:"oneof=asc desc"`
	Category  string `json:"category"`
}

type GetActiveFAQsQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func NewGetActiveFAQsQueryHandler(repo repositories.FAQRepository) *GetActiveFAQsQueryHandler {
	return &GetActiveFAQsQueryHandler{faqRepo: repo}
}

func (h *GetActiveFAQsQueryHandler) HandleGetActiveFAQs(ctx context.Context, query GetActiveFAQsQuery) (*dtos.QueryResult, error) {
	if query.Limit == 0 {
		query.Limit = 10
	}
	if query.SortBy == "" {
		query.SortBy = "priority"
	}
	if query.SortOrder == "" {
		query.SortOrder = "desc"
	}

	opts := &models.QueryOptions{
		Pagination: &models.PaginationParams{
			Offset: query.Offset,
			Limit:  query.Limit,
		},
		SortBy: []models.SortBy{
			{
				Field: query.SortBy,
				Order: models.SortOrder(query.SortOrder),
			},
		},
		Filters: make(map[string]interface{}),
	}

	if query.Category != "" {
		opts.Filters["category"] = query.Category
	}

	faqs, err := h.faqRepo.FindActive(ctx, opts)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to find active FAQs: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		FAQs:      faqs,
		Success:   true,
		Message:   "Active FAQs retrieved successfully",
		Timestamp: time.Now(),
	}, nil
}
