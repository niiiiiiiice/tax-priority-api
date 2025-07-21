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
	BulkHandler    *commands.BulkTestimonialCommandHandler
}

func NewTestimonialCommandHandlers(repo repositories.TestimonialRepository) *TestimonialCommandHandlers {
	return &TestimonialCommandHandlers{
		CreateHandler:  commands.NewCreateTestimonialCommandHandler(repo),
		UpdateHandler:  commands.NewUpdateTestimonialCommandHandler(repo),
		DeleteHandler:  commands.NewDeleteTestimonialCommandHandler(repo),
		ApproveHandler: commands.NewApproveTestimonialCommandHandler(repo),
		BulkHandler:    commands.NewBulkTestimonialCommandHandler(repo),
	}
}

// Методы для выполнения команд
func (h *TestimonialCommandHandlers) CreateTestimonial(ctx context.Context, cmd dtos.CreateTestimonialCommand) (*dtos.CommandResult, error) {
	return h.CreateHandler.Handle(ctx, cmd)
}

func (h *TestimonialCommandHandlers) UpdateTestimonial(ctx context.Context, cmd dtos.UpdateTestimonialCommand) (*dtos.CommandResult, error) {
	return h.UpdateHandler.Handle(ctx, cmd)
}

func (h *TestimonialCommandHandlers) DeleteTestimonial(ctx context.Context, cmd dtos.DeleteTestimonialCommand) (*dtos.CommandResult, error) {
	return h.DeleteHandler.Handle(ctx, cmd)
}

func (h *TestimonialCommandHandlers) ApproveTestimonial(ctx context.Context, cmd dtos.ApproveTestimonialCommand) (*dtos.CommandResult, error) {
	return h.ApproveHandler.Handle(ctx, cmd)
}

func (h *TestimonialCommandHandlers) BulkApproveTestimonials(ctx context.Context, cmd dtos.BulkApproveTestimonialsCommand) (*dtos.CommandResult, error) {
	return h.BulkHandler.HandleBulkApprove(ctx, cmd)
}

func (h *TestimonialCommandHandlers) BulkDeactivateTestimonials(ctx context.Context, cmd dtos.BulkDeactivateTestimonialsCommand) (*dtos.CommandResult, error) {
	return h.BulkHandler.HandleBulkDeactivate(ctx, cmd)
}

func (h *TestimonialCommandHandlers) BulkActivateTestimonials(ctx context.Context, cmd dtos.BulkActivateTestimonialsCommand) (*dtos.CommandResult, error) {
	return h.BulkHandler.HandleBulkActivate(ctx, cmd)
}

func (h *TestimonialCommandHandlers) BulkDeleteTestimonials(ctx context.Context, cmd dtos.BulkDeleteTestimonialsCommand) (*dtos.CommandResult, error) {
	return h.BulkHandler.HandleBulkDelete(ctx, cmd)
}
