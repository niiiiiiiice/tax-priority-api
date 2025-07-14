package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/domain/repositories"
	"time"
)

type GetFAQCategoriesQuery struct {
	WithCounts bool `json:"withCounts"`
}

type GetFAQCategoriesQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func (h *GetFAQCategoriesQueryHandler) HandleGetFAQCategories(ctx context.Context, query GetFAQCategoriesQuery) (*dtos.QueryResult, error) {
	if query.WithCounts {
		categoryCounts, err := h.faqRepo.GetCategoriesWithCounts(ctx)
		if err != nil {
			return &dtos.QueryResult{
				Success:   false,
				Error:     fmt.Sprintf("failed to get categories with counts: %v", err),
				Timestamp: time.Now(),
			}, err
		}

		return &dtos.QueryResult{
			CategoryCounts: categoryCounts,
			Success:        true,
			Message:        "Categories with counts retrieved successfully",
			Timestamp:      time.Now(),
		}, nil
	}

	categories, err := h.faqRepo.GetCategories(ctx)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get categories: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		Categories: categories,
		Success:    true,
		Message:    "Categories retrieved successfully",
		Timestamp:  time.Now(),
	}, nil
}
