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

// CreateFAQKeyGenerator создает генератор ключей для FAQ
func CreateFAQKeyGenerator() appCache.KeyGenerator[*entities.FAQ, string] {
	return appCache.NewKeyGenerator(
		"faq",
		func(faq *entities.FAQ) string { return faq.GetID() },
		func(id string) string { return id },
	)
}

// CreateFAQInvalidationConfig создает конфигурацию инвалидации для FAQ
func CreateFAQInvalidationConfig() *appCache.InvalidationConfig {
	return &appCache.InvalidationConfig{
		Mode:              appCache.InvalidationModeSelective,
		BatchSize:         100,
		InvalidateRelated: true,
	}
}

// CreateFAQCacheManager создает менеджер кеша для FAQ
func CreateFAQCacheManager(
	cache appCache.Cache,
	keyGen appCache.KeyGenerator[*entities.FAQ, string],
	cacheConfig *appCache.CacheConfig,
	invalidationConfig *appCache.InvalidationConfig,
) infraCache.CacheManager[*entities.FAQ, string] {
	return infraCache.NewCacheManager(cache, keyGen, cacheConfig, invalidationConfig)
}
