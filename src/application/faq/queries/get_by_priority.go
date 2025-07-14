package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/shared/models"
	"time"
)

type GetFAQsByPriorityQuery struct {
	MinPriority int    `json:"minPriority" validate:"min=0,max=100"`
	Limit       int    `json:"limit" validate:"min=1,max=100"`
	Offset      int    `json:"offset" validate:"min=0"`
	SortBy      string `json:"sortBy"`
	SortOrder   string `json:"sortOrder" validate:"oneof=asc desc"`
	ActiveOnly  bool   `json:"activeOnly"`
}

type GetFAQsByPriorityQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func NewGetFAQsByPriorityQueryHandler(repo repositories.FAQRepository) *GetFAQsByPriorityQueryHandler {
	return &GetFAQsByPriorityQueryHandler{faqRepo: repo}
}

func (h *GetFAQsByPriorityQueryHandler) HandleGetFAQsByPriority(ctx context.Context, query GetFAQsByPriorityQuery) (*dtos.QueryResult, error) {
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

	if query.ActiveOnly {
		opts.Filters["isActive"] = true
	}

	faqs, err := h.faqRepo.FindByPriority(ctx, query.MinPriority, opts)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to find FAQs by priority: %v", err),
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
