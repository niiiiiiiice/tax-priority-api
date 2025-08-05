package repositories

import (
	appCache "tax-priority-api/src/application/cache"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/cache"
)

// NewCachedFAQRepository создает кешированный FAQ репозиторий
func NewCachedFAQRepository(
	baseRepo repositories.GenericRepository[*entities.FAQ, string],
	cacheManager cache.CacheManager[*entities.FAQ, string],
	keyGen appCache.KeyGenerator[*entities.FAQ, string],
	config *appCache.CacheConfig,
) repositories.CachedFAQRepository {
	return NewCachedGenericRepository(baseRepo, cacheManager, keyGen, config)
}
