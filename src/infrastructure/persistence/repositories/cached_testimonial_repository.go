package repositories

import (
	appCache "tax-priority-api/src/application/cache"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/cache"
)

// NewCachedTestimonialRepository создает кешированный Testimonial репозиторий
func NewCachedTestimonialRepository(
	baseRepo repositories.GenericRepository[*entities.Testimonial, string],
	cacheManager cache.CacheManager[*entities.Testimonial, string],
	keyGen appCache.KeyGenerator[*entities.Testimonial, string],
	config *appCache.CacheConfig,
) repositories.CachedTestimonialRepository {
	return NewCachedGenericRepository(baseRepo, cacheManager, keyGen, config)
}
