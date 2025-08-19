package repositories

import (
	appCache "tax-priority-api/src/application/cache"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/cache"
)

type CachedFeatureRepositoryImpl struct {
	repositories.GenericRepository[*entities.Feature, string]
	faqRepo      repositories.FAQRepository
	cacheManager cache.CacheManager[*entities.Feature, string]
	keyGen       appCache.KeyGenerator[*entities.Feature, string]
	config       *appCache.CacheConfig
}

func NewCachedFeatureRepository(
	baseRepo repositories.GenericRepository[*entities.Feature, string],
	faqRepo repositories.FAQRepository,
	cacheManager cache.CacheManager[*entities.Feature, string],
	keyGen appCache.KeyGenerator[*entities.Feature, string],
	config *appCache.CacheConfig,
) repositories.CachedFeatureRepository {
	return &CachedFeatureRepositoryImpl{
		GenericRepository: NewCachedGenericRepository(baseRepo, cacheManager, keyGen, config),
		faqRepo:           faqRepo,
		cacheManager:      cacheManager,
		keyGen:            keyGen,
		config:            config,
	}
}
