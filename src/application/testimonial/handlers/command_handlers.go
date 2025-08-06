package handlers

import (
	"context"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/commands"
	"tax-priority-api/src/application/testimonial/dtos"
)

type TestimonialCommandHandlers struct {
	CreateHandler  *commands.CreateTestimonialCommandHandler
	UpdateHandler  *commands.UpdateTestimonialCommandHandler
	DeleteHandler  *commands.DeleteTestimonialCommandHandler
	ApproveHandler *commands.ApproveTestimonialCommandHandler
}

func NewTestimonialCommandHandlers(repo repositories.CachedTestimonialRepository) *TestimonialCommandHandlers {
	return &TestimonialCommandHandlers{
		CreateHandler:  commands.NewCreateTestimonialCommandHandler(repo),
		UpdateHandler:  commands.NewUpdateTestimonialCommandHandler(repo),
		DeleteHandler:  commands.NewDeleteTestimonialCommandHandler(repo),
		ApproveHandler: commands.NewApproveTestimonialCommandHandler(repo),
	}
}

// Методы для выполнения команд

// CreateTestimonial - создание отзыва
func (h *TestimonialCommandHandlers) CreateTestimonial(ctx context.Context, cmd dtos.CreateTestimonialCommand) (*dtos.CommandResult, error) {
	return h.CreateHandler.Handle(ctx, cmd)
}

// UpdateTestimonial - обновление отзыва
func (h *TestimonialCommandHandlers) UpdateTestimonial(ctx context.Context, cmd dtos.UpdateTestimonialCommand) (*dtos.CommandResult, error) {
	return h.UpdateHandler.Handle(ctx, cmd)
}

// DeleteTestimonial - удаление отзыва
func (h *TestimonialCommandHandlers) DeleteTestimonial(ctx context.Context, cmd dtos.DeleteTestimonialCommand) (*dtos.CommandResult, error) {
	return h.DeleteHandler.Handle(ctx, cmd)
}

// ApproveTestimonial - одобрение отзыва
func (h *TestimonialCommandHandlers) ApproveTestimonial(ctx context.Context, cmd dtos.ApproveTestimonialCommand) (*dtos.CommandResult, error) {
	return h.ApproveHandler.Handle(ctx, cmd)
}
