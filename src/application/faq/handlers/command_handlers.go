package handlers

import (
	"tax-priority-api/src/application/events"
	"tax-priority-api/src/application/faq/commands"
	"tax-priority-api/src/application/repositories"
)

type FAQCommandHandlers struct {
	Activate       *commands.ActivateFAQCommandHandler
	BulkDelete     *commands.BulkDeleteFAQCommandHandler
	Deactivate     *commands.DeactivateFAQCommandHandler
	Delete         *commands.DeleteFAQCommandHandler
	Create         *commands.CreateFAQCommandHandler
	Update         *commands.UpdateFAQCommandHandler
	UpdateCategory *commands.UpdateFAQCategoryCommandHandler
	UpdatePriority *commands.UpdateFAQPriorityCommandHandler
}

func NewFAQCommandHandlers(repo repositories.CachedFAQRepository, notificationService events.NotificationService) *FAQCommandHandlers {
	return &FAQCommandHandlers{
		Activate:       commands.NewActivateFAQCommandHandler(repo, notificationService),
		BulkDelete:     commands.NewBulkDeleteFAQCommandHandler(repo, notificationService),
		Deactivate:     commands.NewDeactivateFAQCommandHandler(repo, notificationService),
		Delete:         commands.NewDeleteFAQCommandHandler(repo, notificationService),
		Create:         commands.NewCreateFAQCommandHandler(repo, notificationService),
		Update:         commands.NewUpdateFAQCommandHandler(repo, notificationService),
		UpdateCategory: commands.NewUpdateFAQCategoryCommandHandler(repo, notificationService),
		UpdatePriority: commands.NewUpdateFAQPriorityCommandHandler(repo, notificationService),
	}
}
