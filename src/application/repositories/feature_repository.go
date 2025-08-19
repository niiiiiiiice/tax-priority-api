package repositories

import "tax-priority-api/src/domain/entities"

// FeatureRepository определяет интерфейс для работы с Feature
type FeatureRepository interface {
	GenericRepository[*entities.Feature, string]
}
