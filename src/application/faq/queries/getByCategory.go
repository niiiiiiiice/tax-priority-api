package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/shared/models"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/domain/repositories"
	"time"
)

type GetFAQsByCategoryQuery struct {
	Category   string `json:"category" validate:"required"`
	Limit      int    `json:"limit" validate:"min=1,max=100"`
	Offset     int    `json:"offset" validate:"min=0"`
	SortBy     string `json:"sortBy"`
	SortOrder  string `json:"sortOrder" validate:"oneof=asc desc"`
	ActiveOnly bool   `json:"activeOnly"`
}

type GetFAQsByCategoryQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func (h *GetFAQsByCategoryQueryHandler) HandleGetFAQsByCategory(ctx context.Context, query GetFAQsByCategoryQuery) (*dtos.QueryResult, error) {
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

	var faqs []*entities.FAQ
	var err error

	if query.ActiveOnly {
		faqs, err = h.faqRepo.FindActive(ctx, opts)
	} else {
		faqs, err = h.faqRepo.FindByCategory(ctx, query.Category, opts)
	}

	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to find FAQs by category: %v", err),
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
