package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"time"
)

type GetTestimonialByIDQueryHandler struct {
	testimonialRepo repositories.TestimonialRepository
}

func NewGetTestimonialByIDQueryHandler(repo repositories.TestimonialRepository) *GetTestimonialByIDQueryHandler {
	return &GetTestimonialByIDQueryHandler{
		testimonialRepo: repo,
	}
}

func (h *GetTestimonialByIDQueryHandler) Handle(ctx context.Context, query dtos.GetTestimonialByIDQuery) (*dtos.QueryResult, error) {
	testimonial, err := h.testimonialRepo.FindByID(ctx, query.ID)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("testimonial not found: %v", err),
			Timestamp: time.Now(),
		}, err
	}

	return &dtos.QueryResult{
		Success:   true,
		Message:   "Testimonial retrieved successfully",
		Data:      testimonial,
		Timestamp: time.Now(),
	}, nil
}
