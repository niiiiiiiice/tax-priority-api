package wire

import (
	"gorm.io/gorm"

	appCache "tax-priority-api/src/application/cache"
	appRepos "tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	infraCache "tax-priority-api/src/infrastructure/cache"
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

// CreateTestimonialKeyGenerator создает генератор ключей для Testimonial
func CreateTestimonialKeyGenerator() appCache.KeyGenerator[*entities.Testimonial, string] {
	return appCache.NewKeyGenerator(
		"testimonial",
		func(testimonial *entities.Testimonial) string { return testimonial.GetID() },
		func(id string) string { return id },
	)
}

// CreateTestimonialInvalidationConfig создает конфигурацию инвалидации для Testimonial
func CreateTestimonialInvalidationConfig() *appCache.InvalidationConfig {
	return &appCache.InvalidationConfig{
		Mode:              appCache.InvalidationModeSelective,
		BatchSize:         100,
		InvalidateRelated: true,
	}
}

// CreateTestimonialCacheManager создает менеджер кеша для Testimonial
func CreateTestimonialCacheManager(
	cache appCache.Cache,
	keyGen appCache.KeyGenerator[*entities.Testimonial, string],
	cacheConfig *appCache.CacheConfig,
	invalidationConfig *appCache.InvalidationConfig,
) infraCache.CacheManager[*entities.Testimonial, string] {
	return infraCache.NewCacheManager(cache, keyGen, cacheConfig, invalidationConfig)
}
