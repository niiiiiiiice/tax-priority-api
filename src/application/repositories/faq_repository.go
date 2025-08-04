package repositories

import (
	"tax-priority-api/src/domain/entities"
)

// FAQRepository определяет интерфейс для работы с FAQ
type FAQRepository interface {
	GenericRepository[*entities.FAQ, string]
}
