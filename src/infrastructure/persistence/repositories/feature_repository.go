package repositories

import (
	"tax-priority-api/src/domain/entities"
)
import "tax-priority-api/src/application/repositories"

type FeatureRepositoryImpl struct {
	repositories.GenericRepository[*entities.Feature, string]
}

func NewFeatureRepository(generic repositories.GenericRepository[*entities.Feature, string]) repositories.FeatureRepository {
	return &FeatureRepositoryImpl{generic}
}
