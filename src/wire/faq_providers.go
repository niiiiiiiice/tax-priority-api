package wire

import (
	"gorm.io/gorm"

	appRepos "tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	infraModels "tax-priority-api/src/infrastructure/persistence/models"
	infraRepos "tax-priority-api/src/infrastructure/persistence/repositories"
)

// CreateFAQGenericRepository создает GenericRepository для FAQ
// Эта функция изолирована от Wire чтобы избежать проблем с AST
func CreateFAQGenericRepository(db *gorm.DB) appRepos.GenericRepository[*entities.FAQ, string] {
	domainToModel := func(entity *entities.FAQ) *infraModels.FAQModel {
		return infraModels.NewFAQModelFromEntity(entity)
	}
	modelToDomain := func(model *infraModels.FAQModel) *entities.FAQ {
		return model.ToEntity()
	}
	return infraRepos.NewGenericRepository(
		db,
		domainToModel,
		modelToDomain,
	)
}
