package repositories

import (
	"context"
	appCache "tax-priority-api/src/application/cache"
	"tax-priority-api/src/application/faq/dtos"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/cache"
	"time"
)

type CachedFAQRepositoryImpl struct {
	repositories.GenericRepository[*entities.FAQ, string]
	faqRepo      repositories.FAQRepository
	cacheManager cache.CacheManager[*entities.FAQ, string]
	keyGen       appCache.KeyGenerator[*entities.FAQ, string]
	config       *appCache.CacheConfig
}

// NewCachedFAQRepository создает кешированный FAQ репозиторий
func NewCachedFAQRepository(
	baseRepo repositories.GenericRepository[*entities.FAQ, string],
	faqRepo repositories.FAQRepository,
	cacheManager cache.CacheManager[*entities.FAQ, string],
	keyGen appCache.KeyGenerator[*entities.FAQ, string],
	config *appCache.CacheConfig,
) repositories.CachedFAQRepository {
	return &CachedFAQRepositoryImpl{
		GenericRepository: NewCachedGenericRepository(baseRepo, cacheManager, keyGen, config),
		faqRepo:           faqRepo,
		cacheManager:      cacheManager,
		keyGen:            keyGen,
		config:            config,
	}
}

// GetCategories возвращает список категорий FAQ с кешированием
func (r *CachedFAQRepositoryImpl) GetCategories(ctx context.Context, withCounts bool) ([]string, map[string]int64, error) {
	cacheKey := GenerateFAQCategoriesKey(withCounts)

	cached, err := r.cacheManager.GetQuery(ctx, cacheKey, func() (interface{}, error) {
		categories, categoryCounts, err := r.faqRepo.GetCategories(ctx, withCounts)
		if err != nil {
			return nil, err
		}

		result := &dtos.CategoriesResult{
			Categories:     categories,
			CategoryCounts: categoryCounts,
			WithCounts:     withCounts,
			CachedAt:       time.Now(),
		}

		return result, nil
	}, r.config.DefaultTTL)

	if err != nil {
		return r.faqRepo.GetCategories(ctx, withCounts)
	}

	if result, ok := cached.(*dtos.CategoriesResult); ok {
		return result.Categories, result.CategoryCounts, nil
	}

	return r.faqRepo.GetCategories(ctx, withCounts)
}

func (r *CachedFAQRepositoryImpl) invalidateCategoriesCache(ctx context.Context) error {
	return r.cacheManager.InvalidatePattern(ctx, FAQCategoriesPattern)
}

func (r *CachedFAQRepositoryImpl) Create(ctx context.Context, entity *entities.FAQ) error {
	err := r.GenericRepository.Create(ctx, entity)
	if err != nil {
		return err
	}

	_ = r.invalidateCategoriesCache(ctx)
	return nil
}

func (r *CachedFAQRepositoryImpl) Update(ctx context.Context, entity *entities.FAQ) error {
	err := r.GenericRepository.Update(ctx, entity)
	if err != nil {
		return err
	}

	_ = r.invalidateCategoriesCache(ctx)
	return nil
}

func (r *CachedFAQRepositoryImpl) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	err := r.GenericRepository.UpdateFields(ctx, id, fields)
	if err != nil {
		return err
	}

	if _, hasCategory := fields["category"]; hasCategory {
		_ = r.invalidateCategoriesCache(ctx)
	}
	return nil
}

func (r *CachedFAQRepositoryImpl) Delete(ctx context.Context, id string) error {
	err := r.GenericRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	_ = r.invalidateCategoriesCache(ctx)
	return nil
}

func (r *CachedFAQRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	err := r.GenericRepository.SoftDelete(ctx, id)
	if err != nil {
		return err
	}

	_ = r.invalidateCategoriesCache(ctx)
	return nil
}

func (r *CachedFAQRepositoryImpl) CreateBatch(ctx context.Context, entities []*entities.FAQ) (*models.BulkOperationResult, error) {
	result, err := r.GenericRepository.CreateBatch(ctx, entities)
	if err != nil {
		return result, err
	}

	_ = r.invalidateCategoriesCache(ctx)
	return result, nil
}

func (r *CachedFAQRepositoryImpl) UpdateBatch(ctx context.Context, entities []*entities.FAQ) (*models.BulkOperationResult, error) {
	result, err := r.GenericRepository.UpdateBatch(ctx, entities)
	if err != nil {
		return result, err
	}

	_ = r.invalidateCategoriesCache(ctx)
	return result, nil
}

func (r *CachedFAQRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*models.BulkOperationResult, error) {
	result, err := r.GenericRepository.DeleteBatch(ctx, ids)
	if err != nil {
		return result, err
	}

	_ = r.invalidateCategoriesCache(ctx)
	return result, nil
}
