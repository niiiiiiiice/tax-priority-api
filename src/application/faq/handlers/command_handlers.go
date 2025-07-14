package handlers

import (
	"tax-priority-api/src/application/faq/commands"
	"tax-priority-api/src/application/repositories"
)

type FAQCommandHandlers struct {
	Activate       *commands.ActivateFAQCommandHandler
	BulkDelete     *commands.BulkDeleteFAQCommandHandler
	Create         *commands.CreateFAQCommandHandler
	CreateBatch    *commands.CreateFAQBatchCommandHandler
	Deactivate     *commands.DeactivateFAQCommandHandler
	Delete         *commands.DeleteFAQCommandHandler
	Update         *commands.UpdateFAQCommandHandler
	UpdateCategory *commands.UpdateFAQCategoryCommandHandler
	UpdatePriority *commands.UpdateFAQPriorityCommandHandler
}

func NewFAQCommandHandlers(repo repositories.FAQRepository) *FAQCommandHandlers {
	return &FAQCommandHandlers{
		Activate:       commands.NewActivateFAQCommandHandler(repo),
		BulkDelete:     commands.NewBulkDeleteFAQCommandHandler(repo),
		Create:         commands.NewCreateFAQCommandHandler(repo),
		CreateBatch:    commands.NewCreateFAQBatchCommandHandler(repo),
		Deactivate:     commands.NewDeactivateFAQCommandHandler(repo),
		Delete:         commands.NewDeleteFAQCommandHandler(repo),
		Update:         commands.NewUpdateFAQCommandHandler(repo),
		UpdateCategory: commands.NewUpdateFAQCategoryCommandHandler(repo),
		UpdatePriority: commands.NewUpdateFAQPriorityCommandHandler(repo),
	}
}
