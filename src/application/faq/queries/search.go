package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"time"
)

type SearchFAQsQuery struct {
	Query      string `json:"query" validate:"required,min=3"`
	Category   string `json:"category"`
	Limit      int    `json:"limit" validate:"min=1,max=100"`
	Offset     int    `json:"offset" validate:"min=0"`
	SortBy     string `json:"sortBy"`
	SortOrder  string `json:"sortOrder" validate:"oneof=asc desc"`
	ActiveOnly bool   `json:"activeOnly"`
}

type SearchFAQsQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func NewSearchFAQsQueryHandler(repo repositories.FAQRepository) *SearchFAQsQueryHandler {
	return &SearchFAQsQueryHandler{faqRepo: repo}
}

func (h *SearchFAQsQueryHandler) HandleSearchFAQs(ctx context.Context, query SearchFAQsQuery) (*dtos.QueryResult, error) {
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

	if query.Category != "" {
		faqs, err = h.faqRepo.SearchByCategory(ctx, query.Query, query.Category, opts)
	} else {
		faqs, err = h.faqRepo.Search(ctx, query.Query, opts)
	}

	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to search FAQs: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		FAQs:      faqs,
		Success:   true,
		Message:   "FAQs search completed successfully",
		Timestamp: time.Now(),
	}, nil
}
