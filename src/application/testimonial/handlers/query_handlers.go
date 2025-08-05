package handlers

import (
	"context"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"tax-priority-api/src/application/testimonial/queries"
)

type TestimonialQueryHandlers struct {
	GetManyHandler *queries.GetTestimonialsQueryHandler
	GetByIDHandler *queries.GetTestimonialByIDQueryHandler
}

func NewTestimonialQueryHandlers(repo repositories.TestimonialRepository) *TestimonialQueryHandlers {
	return &TestimonialQueryHandlers{
		GetManyHandler: queries.NewGetTestimonialsQueryHandler(repo),
		GetByIDHandler: queries.NewGetTestimonialByIDQueryHandler(repo),
	}
}

// Методы для выполнения запросов
func (h *TestimonialQueryHandlers) GetTestimonials(ctx context.Context, query dtos.GetTestimonialsQuery) (*dtos.QueryResult, error) {
	return h.GetManyHandler.Handle(ctx, query)
}

func (h *TestimonialQueryHandlers) GetTestimonialByID(ctx context.Context, query dtos.GetTestimonialByIDQuery) (*dtos.QueryResult, error) {
	return h.GetByIDHandler.Handle(ctx, query)
}
