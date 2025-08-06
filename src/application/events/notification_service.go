package events

import (
	"context"

	"tax-priority-api/src/domain/entities"
)

// NotificationService интерфейс для отправки уведомлений
type NotificationService interface {
	// FAQ события

	// NotifyFAQCreated - создание FAQ
	NotifyFAQCreated(ctx context.Context, faq *entities.FAQ)
	// NotifyFAQUpdated - обновление FAQ
	NotifyFAQUpdated(ctx context.Context, faq *entities.FAQ)
	// NotifyFAQDeleted - удаление FAQ
	NotifyFAQDeleted(ctx context.Context, faqID string)
	// NotifyFAQActivated - активация FAQ
	NotifyFAQActivated(ctx context.Context, faq *entities.FAQ)
	// NotifyFAQDeactivated - деактивация FAQ
	NotifyFAQDeactivated(ctx context.Context, faq *entities.FAQ)
	// NotifyFAQPriorityChanged - изменение приоритета FAQ
	NotifyFAQPriorityChanged(ctx context.Context, faq *entities.FAQ, oldPriority int)
	// NotifyFAQCategoryChanged - изменение категории FAQ
	NotifyFAQCategoryChanged(ctx context.Context, faq *entities.FAQ, oldCategory string)
	// NotifyFAQBatchCreated - создание пачки FAQ
	NotifyFAQBatchCreated(ctx context.Context, faqs []*entities.FAQ)
	// NotifyFAQBatchDeleted - удаление пачки FAQ
	NotifyFAQBatchDeleted(ctx context.Context, faqIDs []string)

	// Системные события

	// NotifySystemEvent - системное событие
	NotifySystemEvent(ctx context.Context, event string, data interface{})

	// Статистика и состояние

	// GetStats - получение статистики
	GetStats() map[string]interface{}
	// IsEnabled - проверка, включен ли сервис
	IsEnabled() bool
}
