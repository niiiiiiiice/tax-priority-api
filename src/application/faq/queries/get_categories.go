package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"time"
)

type GetFAQCategoriesQuery struct {
	WithCounts bool `json:"withCounts"`
}

type GetFAQCategoriesQueryHandler struct {
	faqRepo repositories.FAQRepository
}

func NewGetFAQCategoriesQueryHandler(repo repositories.FAQRepository) *GetFAQCategoriesQueryHandler {
	return &GetFAQCategoriesQueryHandler{faqRepo: repo}
}

func (h *GetFAQCategoriesQueryHandler) HandleGetFAQCategories(ctx context.Context, query GetFAQCategoriesQuery) (*dtos.QueryResult, error) {
	categories, categoryCounts, err := h.faqRepo.GetCategories(ctx, query.WithCounts)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get categories: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	result := &dtos.QueryResult{
		Categories: categories,
		Success:    true,
		Message:    "Categories retrieved successfully",
		Timestamp:  time.Now(),
	}

	if query.WithCounts {
		result.CategoryCounts = categoryCounts
	}

	return result, nil
}
