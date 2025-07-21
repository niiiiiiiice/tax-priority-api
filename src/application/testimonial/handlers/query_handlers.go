package handlers

import (
	"context"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"tax-priority-api/src/application/testimonial/queries"
)

type TestimonialQueryHandlers struct {
	GetManyHandler     *queries.GetTestimonialsQueryHandler
	GetByIDHandler     *queries.GetTestimonialByIDQueryHandler
	GetApprovedHandler *queries.GetApprovedTestimonialsQueryHandler
	GetStatsHandler    *queries.GetTestimonialStatsQueryHandler
}

func NewTestimonialQueryHandlers(repo repositories.TestimonialRepository) *TestimonialQueryHandlers {
	return &TestimonialQueryHandlers{
		GetManyHandler:     queries.NewGetTestimonialsQueryHandler(repo),
		GetByIDHandler:     queries.NewGetTestimonialByIDQueryHandler(repo),
		GetApprovedHandler: queries.NewGetApprovedTestimonialsQueryHandler(repo),
		GetStatsHandler:    queries.NewGetTestimonialStatsQueryHandler(repo),
	}
}

// Методы для выполнения запросов
func (h *TestimonialQueryHandlers) GetTestimonials(ctx context.Context, query dtos.GetTestimonialsQuery) (*dtos.QueryResult, error) {
	return h.GetManyHandler.Handle(ctx, query)
}

func (h *TestimonialQueryHandlers) GetTestimonialByID(ctx context.Context, query dtos.GetTestimonialByIDQuery) (*dtos.QueryResult, error) {
	return h.GetByIDHandler.Handle(ctx, query)
}

func (h *TestimonialQueryHandlers) GetApprovedTestimonials(ctx context.Context, query dtos.GetApprovedTestimonialsQuery) (*dtos.QueryResult, error) {
	return h.GetApprovedHandler.Handle(ctx, query)
}

func (h *TestimonialQueryHandlers) GetTestimonialStats(ctx context.Context, query dtos.GetTestimonialStatsQuery) (*dtos.QueryResult, error) {
	return h.GetStatsHandler.Handle(ctx, query)
}
