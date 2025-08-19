package commands

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/events"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
)

type CreateFeatureCommand struct {
	Name string
}

type CreateFeatureCommandHandler struct {
	repo                repositories.FeatureRepository
	notificationService events.NotificationService
}

func NewCreateFeatureCommandHandler(repo repositories.FeatureRepository, notificationService events.NotificationService) *CreateFeatureCommandHandler {
	return &CreateFeatureCommandHandler{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (h *CreateFeatureCommandHandler) HandleCreateFeature(ctx context.Context, cmd CreateFeatureCommand) (*dtos.CommandResult, error) {
	feature, err := entities.NewFeature(cmd.Name)

	if err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create Feature: %v", err),
		}, err
	}

	if err = h.repo.Create(ctx, feature); err != nil {
		return &dtos.CommandResult{
			Success: false,
			Error:   fmt.Sprintf("failed to create Feature: %v", err),
		}, err
	}

	return &dtos.CommandResult{
		ID:        feature.ID,
		Success:   true,
		Message:   "Feature created successfully",
		CreatedAt: feature.CreatedAt,
		UpdatedAt: feature.UpdatedAt,
	}, nil
}
