package repositories

import (
	"context"
	"tax-priority-api/src/domain/entities"
)

// FAQRepository определяет интерфейс для работы с FAQ
type FAQRepository interface {
	GenericRepository[*entities.FAQ, string]
	// GetCategories возвращает список категорий FAQ
	// Если withCounts = true, также возвращает количество FAQ в каждой категории
	GetCategories(ctx context.Context, withCounts bool) ([]string, map[string]int64, error)
}
