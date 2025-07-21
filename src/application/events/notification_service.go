package events

import (
	"context"

	"tax-priority-api/src/domain/entities"
)

// NotificationService интерфейс для отправки уведомлений
type NotificationService interface {
	// FAQ события
	NotifyFAQCreated(ctx context.Context, faq *entities.FAQ)
	NotifyFAQUpdated(ctx context.Context, faq *entities.FAQ)
	NotifyFAQDeleted(ctx context.Context, faqID string)
	NotifyFAQActivated(ctx context.Context, faq *entities.FAQ)
	NotifyFAQDeactivated(ctx context.Context, faq *entities.FAQ)
	NotifyFAQPriorityChanged(ctx context.Context, faq *entities.FAQ, oldPriority int)
	NotifyFAQCategoryChanged(ctx context.Context, faq *entities.FAQ, oldCategory string)
	NotifyFAQBatchCreated(ctx context.Context, faqs []*entities.FAQ)
	NotifyFAQBatchDeleted(ctx context.Context, faqIDs []string)

	// Системные события
	NotifySystemEvent(ctx context.Context, event string, data interface{})

	// Статистика и состояние
	GetStats() map[string]interface{}
	IsEnabled() bool
}
