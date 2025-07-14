package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/shared/models"
	"time"
)

type GetFAQsQuery struct {
	Limit     int                    `json:"limit" validate:"min=1,max=100"`
	Offset    int                    `json:"offset" validate:"min=0"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder" validate:"oneof=asc desc"`
	Filters   map[string]interface{} `json:"filters"`
}

type GetFAQsQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func NewGetFAQsQueryHandler(repo repositories.FAQRepository) *GetFAQsQueryHandler {
	return &GetFAQsQueryHandler{faqRepo: repo}
}

func (h *GetFAQsQueryHandler) HandleGetFAQs(ctx context.Context, query GetFAQsQuery) (*dtos.QueryResult, error) {
	if query.Limit == 0 {
		query.Limit = 10
	}
	if query.SortBy == "" {
		query.SortBy = "createdAt"
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
		Filters: query.Filters,
	}

	paginated, err := h.faqRepo.FindWithPagination(ctx, opts)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to find FAQs: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		Paginated: paginated,
		Success:   true,
		Message:   "FAQs retrieved successfully",
		Timestamp: time.Now(),
	}, nil
}
