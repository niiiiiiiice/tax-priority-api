package events

import (
	"context"
	"log"
	"time"

	"tax-priority-api/src/application/events"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/websocket"
)

// NotificationServiceImpl реализация сервиса для отправки уведомлений
type NotificationServiceImpl struct {
	hub *websocket.Hub
}

// NewNotificationService создает новый сервис уведомлений
func NewNotificationService(hub *websocket.Hub) events.NotificationService {
	return &NotificationServiceImpl{
		hub: hub,
	}
}

// Event типы событий
type Event struct {
	Entity   string      `json:"entity"`
	Action   string      `json:"action"`
	EntityID string      `json:"entityId,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

// FAQ события
const (
	FAQEntity = "faq"

	// Действия
	ActionCreated         = "created"
	ActionUpdated         = "updated"
	ActionDeleted         = "deleted"
	ActionActivated       = "activated"
	ActionDeactivated     = "deactivated"
	ActionPriorityChanged = "priority_changed"
	ActionCategoryChanged = "category_changed"
)

// NotifyFAQCreated отправляет уведомление о создании FAQ
func (s *NotificationServiceImpl) NotifyFAQCreated(ctx context.Context, faq *entities.FAQ) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity:   FAQEntity,
		Action:   ActionCreated,
		EntityID: faq.ID,
		Data:     faq,
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, event.EntityID, event.Data)
	log.Printf("Sent FAQ created notification for ID: %s", faq.ID)
}

// NotifyFAQUpdated отправляет уведомление об обновлении FAQ
func (s *NotificationServiceImpl) NotifyFAQUpdated(ctx context.Context, faq *entities.FAQ) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity:   FAQEntity,
		Action:   ActionUpdated,
		EntityID: faq.ID,
		Data:     faq,
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, event.EntityID, event.Data)
	log.Printf("Sent FAQ updated notification for ID: %s", faq.ID)
}

// NotifyFAQDeleted отправляет уведомление об удалении FAQ
func (s *NotificationServiceImpl) NotifyFAQDeleted(ctx context.Context, faqID string) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity:   FAQEntity,
		Action:   ActionDeleted,
		EntityID: faqID,
		Data:     map[string]string{"id": faqID},
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, event.EntityID, event.Data)
	log.Printf("Sent FAQ deleted notification for ID: %s", faqID)
}

// NotifyFAQActivated отправляет уведомление об активации FAQ
func (s *NotificationServiceImpl) NotifyFAQActivated(ctx context.Context, faq *entities.FAQ) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity:   FAQEntity,
		Action:   ActionActivated,
		EntityID: faq.ID,
		Data:     faq,
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, event.EntityID, event.Data)
	log.Printf("Sent FAQ activated notification for ID: %s", faq.ID)
}

// NotifyFAQDeactivated отправляет уведомление о деактивации FAQ
func (s *NotificationServiceImpl) NotifyFAQDeactivated(ctx context.Context, faq *entities.FAQ) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity:   FAQEntity,
		Action:   ActionDeactivated,
		EntityID: faq.ID,
		Data:     faq,
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, event.EntityID, event.Data)
	log.Printf("Sent FAQ deactivated notification for ID: %s", faq.ID)
}

// NotifyFAQPriorityChanged отправляет уведомление об изменении приоритета FAQ
func (s *NotificationServiceImpl) NotifyFAQPriorityChanged(ctx context.Context, faq *entities.FAQ, oldPriority int) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity:   FAQEntity,
		Action:   ActionPriorityChanged,
		EntityID: faq.ID,
		Data: map[string]interface{}{
			"faq":         faq,
			"oldPriority": oldPriority,
			"newPriority": faq.Priority,
		},
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, event.EntityID, event.Data)
	log.Printf("Sent FAQ priority changed notification for ID: %s (from %d to %d)", faq.ID, oldPriority, faq.Priority)
}

// NotifyFAQCategoryChanged отправляет уведомление об изменении категории FAQ
func (s *NotificationServiceImpl) NotifyFAQCategoryChanged(ctx context.Context, faq *entities.FAQ, oldCategory string) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity:   FAQEntity,
		Action:   ActionCategoryChanged,
		EntityID: faq.ID,
		Data: map[string]interface{}{
			"faq":         faq,
			"oldCategory": oldCategory,
			"newCategory": faq.Category,
		},
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, event.EntityID, event.Data)
	log.Printf("Sent FAQ category changed notification for ID: %s (from %s to %s)", faq.ID, oldCategory, faq.Category)
}

// NotifyFAQBatchCreated отправляет уведомление о массовом создании FAQ
func (s *NotificationServiceImpl) NotifyFAQBatchCreated(ctx context.Context, faqs []*entities.FAQ) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity: FAQEntity,
		Action: "batch_created",
		Data: map[string]interface{}{
			"count": len(faqs),
			"faqs":  faqs,
		},
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, "", event.Data)
	log.Printf("Sent FAQ batch created notification for %d items", len(faqs))
}

// NotifyFAQBatchDeleted отправляет уведомление о массовом удалении FAQ
func (s *NotificationServiceImpl) NotifyFAQBatchDeleted(ctx context.Context, faqIDs []string) {
	if s.hub == nil {
		return
	}

	event := Event{
		Entity: FAQEntity,
		Action: "batch_deleted",
		Data: map[string]interface{}{
			"count": len(faqIDs),
			"ids":   faqIDs,
		},
	}

	s.hub.BroadcastEvent(event.Entity, event.Action, "", event.Data)
	log.Printf("Sent FAQ batch deleted notification for %d items", len(faqIDs))
}

// NotifySystemEvent отправляет системное уведомление
func (s *NotificationServiceImpl) NotifySystemEvent(ctx context.Context, event string, data interface{}) {
	if s.hub == nil {
		return
	}

	message := websocket.Message{
		Type:      "system",
		Event:     event,
		Data:      data,
		Timestamp: time.Now(),
	}

	s.hub.BroadcastToAll(message)
	log.Printf("Sent system notification: %s", event)
}

// GetStats возвращает статистику WebSocket подключений
func (s *NotificationServiceImpl) GetStats() map[string]interface{} {
	if s.hub == nil {
		return map[string]interface{}{
			"enabled":     false,
			"connections": 0,
		}
	}

	return map[string]interface{}{
		"enabled":         true,
		"connections":     s.hub.GetClientCount(),
		"faq_subscribers": s.hub.GetSubscriptionCount(FAQEntity),
	}
}

// IsEnabled проверяет, включен ли сервис уведомлений
func (s *NotificationServiceImpl) IsEnabled() bool {
	return s.hub != nil
}
