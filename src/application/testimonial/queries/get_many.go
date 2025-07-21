package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"time"
)

type GetTestimonialsQueryHandler struct {
	testimonialRepo repositories.TestimonialRepository
}

func NewGetTestimonialsQueryHandler(repo repositories.TestimonialRepository) *GetTestimonialsQueryHandler {
	return &GetTestimonialsQueryHandler{
		testimonialRepo: repo,
	}
}

func (h *GetTestimonialsQueryHandler) Handle(ctx context.Context, query dtos.GetTestimonialsQuery) (*dtos.QueryResult, error) {
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

	paginated, err := h.testimonialRepo.FindWithPagination(ctx, opts)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to find testimonials: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		Success:   true,
		Message:   "Testimonials retrieved successfully",
		Paginated: paginated,
		Timestamp: time.Now(),
	}, nil
}
