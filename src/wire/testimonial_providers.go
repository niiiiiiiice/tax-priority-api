package wire

import (
	"gorm.io/gorm"
	appRepos "tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	infraModels "tax-priority-api/src/infrastructure/persistence/models"
	infraRepos "tax-priority-api/src/infrastructure/persistence/repositories"
)

func CreateTestimonialGenericRepository(db *gorm.DB) appRepos.GenericRepository[*entities.Testimonial, string] {
	domainToModel := func(entity *entities.Testimonial) *infraModels.TestimonialModel {
		return infraModels.NewTestimonialModelFromEntity(entity)
	}
	modelToDomain := func(model *infraModels.TestimonialModel) *entities.Testimonial {
		return model.ToEntity()
	}
	return infraRepos.NewGenericRepository(
		db,
		domainToModel,
		modelToDomain,
	)
}
