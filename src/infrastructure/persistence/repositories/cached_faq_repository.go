package repositories

import (
	"context"
	"log"
	"time"

	"tax-priority-api/src/application/cache"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/persistence"
)

// CachedFAQRepositoryImpl кешированный репозиторий FAQ
type CachedFAQRepositoryImpl struct {
	repo       repositories.FAQRepository
	cache      cache.Cache
	keys       *persistence.RedisKeys
	defaultTTL time.Duration
}

// NewCachedFAQRepository создает новый кешированный репозиторий FAQ
func NewCachedFAQRepository(repo repositories.FAQRepository, cache cache.Cache, keys *persistence.RedisKeys) repositories.CachedFAQRepository {
	return &CachedFAQRepositoryImpl{
		repo:       repo,
		cache:      cache,
		keys:       keys,
		defaultTTL: 15 * time.Minute,
	}
}

// Create создает новую FAQ и инвалидирует кеш
func (r *CachedFAQRepositoryImpl) Create(ctx context.Context, faq *entities.FAQ) error {
	err := r.repo.Create(ctx, faq)
	if err != nil {
		return err
	}

	// Инвалидируем кеш
	r.invalidateCache(ctx, faq)
	return nil
}

// CreateBatch создает несколько FAQ и инвалидирует кеш
func (r *CachedFAQRepositoryImpl) CreateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error) {
	result, err := r.repo.CreateBatch(ctx, faqs)
	if err != nil {
		return result, err
	}

	// Инвалидируем кеш для всех FAQ
	for _, faq := range faqs {
		r.invalidateCache(ctx, faq)
	}

	return result, nil
}

// FindByID ищет FAQ по ID с кешированием
func (r *CachedFAQRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.FAQ, error) {
	// Пытаемся получить из кеша
	cacheKey := r.keys.GetFAQByIDKey(id)
	var faq entities.FAQ

	err := r.cache.GetJSON(ctx, cacheKey, &faq)
	if err == nil {
		return &faq, nil
	}

	// Если не найдено в кеше, получаем из базы
	result, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	if err := r.cache.SetJSON(ctx, cacheKey, result, r.defaultTTL); err != nil {
		log.Printf("Failed to cache FAQ by ID %s: %v", id, err)
	}

	return result, nil
}

// FindByIDs ищет FAQ по списку ID
func (r *CachedFAQRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.FAQ, error) {
	var faqs []*entities.FAQ
	var missingIDs []string

	// Проверяем кеш для каждого ID
	for _, id := range ids {
		cacheKey := r.keys.GetFAQByIDKey(id)
		var faq entities.FAQ

		err := r.cache.GetJSON(ctx, cacheKey, &faq)
		if err == nil {
			faqs = append(faqs, &faq)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	// Если есть отсутствующие ID, получаем их из базы
	if len(missingIDs) > 0 {
		missingFAQs, err := r.repo.FindByIDs(ctx, missingIDs)
		if err != nil {
			return nil, err
		}

		// Кешируем полученные FAQ
		for _, faq := range missingFAQs {
			cacheKey := r.keys.GetFAQByIDKey(faq.ID)
			if err := r.cache.SetJSON(ctx, cacheKey, faq, r.defaultTTL); err != nil {
				log.Printf("Failed to cache FAQ by ID %s: %v", faq.ID, err)
			}
		}

		faqs = append(faqs, missingFAQs...)
	}

	return faqs, nil
}

// Update обновляет FAQ и инвалидирует кеш
func (r *CachedFAQRepositoryImpl) Update(ctx context.Context, faq *entities.FAQ) error {
	err := r.repo.Update(ctx, faq)
	if err != nil {
		return err
	}

	// Инвалидируем кеш
	r.invalidateCache(ctx, faq)
	return nil
}

// UpdateBatch обновляет несколько FAQ и инвалидирует кеш
func (r *CachedFAQRepositoryImpl) UpdateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error) {
	result, err := r.repo.UpdateBatch(ctx, faqs)
	if err != nil {
		return result, err
	}

	// Инвалидируем кеш для всех FAQ
	for _, faq := range faqs {
		r.invalidateCache(ctx, faq)
	}

	return result, nil
}

// UpdateFields обновляет поля FAQ и инвалидирует кеш
func (r *CachedFAQRepositoryImpl) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	err := r.repo.UpdateFields(ctx, id, fields)
	if err != nil {
		return err
	}

	// Инвалидируем кеш по ID
	r.invalidateCacheByID(ctx, id)
	return nil
}

// Delete удаляет FAQ и инвалидирует кеш
func (r *CachedFAQRepositoryImpl) Delete(ctx context.Context, id string) error {
	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Инвалидируем кеш по ID
	r.invalidateCacheByID(ctx, id)
	return nil
}

// DeleteBatch удаляет несколько FAQ и инвалидирует кеш
func (r *CachedFAQRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*models.BulkOperationResult, error) {
	result, err := r.repo.DeleteBatch(ctx, ids)
	if err != nil {
		return result, err
	}

	// Инвалидируем кеш для всех ID
	for _, id := range ids {
		r.invalidateCacheByID(ctx, id)
	}

	return result, nil
}

// SoftDelete удаляет FAQ и инвалидирует кеш
func (r *CachedFAQRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	err := r.repo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}

	r.invalidateCacheByID(ctx, id)
	return nil
}

// FindAll возвращает все FAQ с кешированием
func (r *CachedFAQRepositoryImpl) FindAll(ctx context.Context, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	// Для простых запросов можем кешировать
	if opts == nil || (opts.Filters == nil && opts.Pagination == nil) {
		cacheKey := "faq:all"
		var faqs []*entities.FAQ

		err := r.cache.GetJSON(ctx, cacheKey, &faqs)
		if err == nil {
			return faqs, nil
		}

		// Получаем из базы
		result, err := r.repo.FindAll(ctx, opts)
		if err != nil {
			return nil, err
		}

		// Кешируем результат
		if err := r.cache.SetJSON(ctx, cacheKey, result, r.defaultTTL); err != nil {
			log.Printf("Failed to cache all FAQs: %v", err)
		}

		return result, nil
	}

	// Для сложных запросов не кешируем
	return r.repo.FindAll(ctx, opts)
}

// FindOne возвращает один FAQ с кешированием
func (r *CachedFAQRepositoryImpl) FindOne(ctx context.Context, opts *models.QueryOptions) (*entities.FAQ, error) {
	if opts == nil || opts.Filters == nil || len(opts.Filters) == 0 {
		activeResults, err := r.FindActive(ctx, &models.QueryOptions{
			Pagination: &models.PaginationParams{Offset: 0, Limit: 1},
		})
		if err != nil {
			return nil, err
		}
		if len(activeResults) > 0 {
			return activeResults[0], nil
		}
		return nil, nil
	}

	return r.repo.FindOne(ctx, opts)
}

// FindWithPagination возвращает FAQ с пагинацией
func (r *CachedFAQRepositoryImpl) FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.FAQ], error) {
	// Для пагинации обычно не кешируем, так как результаты могут быстро устареть
	return r.repo.FindWithPagination(ctx, opts)
}

// Count возвращает количество FAQ с кешированием
func (r *CachedFAQRepositoryImpl) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	// Кешируем только простые подсчеты
	if len(filters) == 0 {
		cacheKey := r.keys.FAQCount
		var count int64

		err := r.cache.GetJSON(ctx, cacheKey, &count)
		if err == nil {
			return count, nil
		}

		// Получаем из базы
		result, err := r.repo.Count(ctx, filters)
		if err != nil {
			return 0, err
		}

		// Кешируем результат
		if err := r.cache.SetJSON(ctx, cacheKey, result, r.defaultTTL); err != nil {
			log.Printf("Failed to cache FAQ count: %v", err)
		}

		return result, nil
	}

	return r.repo.Count(ctx, filters)
}

// Exists проверяет существование FAQ
func (r *CachedFAQRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	return r.repo.Exists(ctx, id)
}

// ExistsByFields проверяет существование FAQ по полям
func (r *CachedFAQRepositoryImpl) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	if len(filters) == 1 {
		if id, ok := filters["id"].(string); ok {
			cacheKey := r.keys.GetFAQByIDKey(id)
			exists, err := r.cache.Exists(ctx, cacheKey)
			if err == nil {
				return exists, nil
			}
		}

		if category, ok := filters["category"].(string); ok {
			cacheKey := r.keys.GetFAQByCategoryKey(category)
			var faqs []*entities.FAQ
			err := r.cache.GetJSON(ctx, cacheKey, &faqs)
			if err == nil {
				return len(faqs) > 0, nil
			}
		}
	}

	return r.repo.ExistsByFields(ctx, filters)
}

// WithTransaction выполняет операцию в транзакции
func (r *CachedFAQRepositoryImpl) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	return r.repo.WithTransaction(ctx, fn)
}

// Refresh обновляет кеш
func (r *CachedFAQRepositoryImpl) Refresh(ctx context.Context, entity *entities.FAQ) error {
	if entity == nil {
		return nil
	}

	fresh, err := r.repo.FindByID(ctx, entity.ID)
	if err != nil {
		return err
	}

	cacheKey := r.keys.GetFAQByIDKey(entity.ID)
	if err := r.cache.SetJSON(ctx, cacheKey, fresh, r.defaultTTL); err != nil {
		log.Printf("Failed to refresh cache for FAQ ID %s: %v", entity.ID, err)
		return err
	}

	r.invalidateCache(ctx, fresh)

	return nil
}

// Clear очищает кеш
func (r *CachedFAQRepositoryImpl) Clear(ctx context.Context) error {
	patterns := []string{
		"faq:*",
		r.keys.FAQActive,
		r.keys.FAQCount,
		r.keys.FAQCategories,
	}

	for _, pattern := range patterns {
		if err := r.cache.DeletePattern(ctx, pattern); err != nil {
			log.Printf("Failed to clear cache pattern %s: %v", pattern, err)
			return err
		}
	}

	return nil
}

// FindByCategory ищет FAQ по категории с кешированием
func (r *CachedFAQRepositoryImpl) FindByCategory(ctx context.Context, category string, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	// Кешируем только простые запросы по категории
	if opts == nil || (opts.Filters == nil && opts.Pagination == nil) {
		cacheKey := r.keys.GetFAQByCategoryKey(category)
		var faqs []*entities.FAQ

		err := r.cache.GetJSON(ctx, cacheKey, &faqs)
		if err == nil {
			return faqs, nil
		}

		// Получаем из базы
		result, err := r.repo.FindByCategory(ctx, category, opts)
		if err != nil {
			return nil, err
		}

		// Кешируем результат
		if err := r.cache.SetJSON(ctx, cacheKey, result, r.defaultTTL); err != nil {
			log.Printf("Failed to cache FAQs by category %s: %v", category, err)
		}

		return result, nil
	}

	return r.repo.FindByCategory(ctx, category, opts)
}

// FindActive ищет активные FAQ с кешированием
func (r *CachedFAQRepositoryImpl) FindActive(ctx context.Context, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	// Кешируем только простые запросы активных FAQ
	if opts == nil || (opts.Filters == nil && opts.Pagination == nil) {
		cacheKey := r.keys.FAQActive
		var faqs []*entities.FAQ

		err := r.cache.GetJSON(ctx, cacheKey, &faqs)
		if err == nil {
			return faqs, nil
		}

		// Получаем из базы
		result, err := r.repo.FindActive(ctx, opts)
		if err != nil {
			return nil, err
		}

		// Кешируем результат
		if err := r.cache.SetJSON(ctx, cacheKey, result, r.defaultTTL); err != nil {
			log.Printf("Failed to cache active FAQs: %v", err)
		}

		return result, nil
	}

	return r.repo.FindActive(ctx, opts)
}

// FindByPriority ищет FAQ по приоритету
func (r *CachedFAQRepositoryImpl) FindByPriority(ctx context.Context, minPriority int, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	return r.repo.FindByPriority(ctx, minPriority, opts)
}

// Search выполняет поиск FAQ
func (r *CachedFAQRepositoryImpl) Search(ctx context.Context, query string, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	return r.repo.Search(ctx, query, opts)
}

// SearchByCategory выполняет поиск FAQ по категории
func (r *CachedFAQRepositoryImpl) SearchByCategory(ctx context.Context, query string, category string, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	return r.repo.SearchByCategory(ctx, query, category, opts)
}

// CountByCategory возвращает количество FAQ по категории
func (r *CachedFAQRepositoryImpl) CountByCategory(ctx context.Context, category string) (int64, error) {
	return r.repo.CountByCategory(ctx, category)
}

// CountActive возвращает количество активных FAQ
func (r *CachedFAQRepositoryImpl) CountActive(ctx context.Context) (int64, error) {
	return r.repo.CountActive(ctx)
}

// ExistsByQuestion проверяет существование FAQ по вопросу
func (r *CachedFAQRepositoryImpl) ExistsByQuestion(ctx context.Context, question string) (bool, error) {
	return r.repo.ExistsByQuestion(ctx, question)
}

// GetCategories возвращает список категорий с кешированием
func (r *CachedFAQRepositoryImpl) GetCategories(ctx context.Context) ([]string, error) {
	cacheKey := r.keys.FAQCategories
	var categories []string

	err := r.cache.GetJSON(ctx, cacheKey, &categories)
	if err == nil {
		return categories, nil
	}

	// Получаем из базы
	result, err := r.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	// Кешируем результат
	if err := r.cache.SetJSON(ctx, cacheKey, result, r.defaultTTL); err != nil {
		log.Printf("Failed to cache FAQ categories: %v", err)
	}

	return result, nil
}

// GetCategoriesWithCounts возвращает категории с количеством
func (r *CachedFAQRepositoryImpl) GetCategoriesWithCounts(ctx context.Context) (map[string]int64, error) {
	return r.repo.GetCategoriesWithCounts(ctx)
}

// invalidateCache инвалидирует кеш для FAQ
func (r *CachedFAQRepositoryImpl) invalidateCache(ctx context.Context, faq *entities.FAQ) {
	// Удаляем кеш по ID
	cacheKey := r.keys.GetFAQByIDKey(faq.ID)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate cache for FAQ ID %s: %v", faq.ID, err)
	}

	// Удаляем кеш по категории
	categoryKey := r.keys.GetFAQByCategoryKey(faq.Category)
	if err := r.cache.Delete(ctx, categoryKey); err != nil {
		log.Printf("Failed to invalidate cache for FAQ category %s: %v", faq.Category, err)
	}

	// Удаляем общие кеши
	generalKeys := []string{
		r.keys.FAQActive,
		r.keys.FAQCount,
		r.keys.FAQCategories,
		"faq:all",
	}

	for _, key := range generalKeys {
		if err := r.cache.Delete(ctx, key); err != nil {
			log.Printf("Failed to invalidate cache for key %s: %v", key, err)
		}
	}
}

// invalidateCacheByID инвалидирует кеш по ID
func (r *CachedFAQRepositoryImpl) invalidateCacheByID(ctx context.Context, id string) {
	// Удаляем кеш по ID
	cacheKey := r.keys.GetFAQByIDKey(id)
	if err := r.cache.Delete(ctx, cacheKey); err != nil {
		log.Printf("Failed to invalidate cache for FAQ ID %s: %v", id, err)
	}

	// Удаляем общие кеши
	generalKeys := []string{
		r.keys.FAQActive,
		r.keys.FAQCount,
		r.keys.FAQCategories,
		"faq:all",
	}

	for _, key := range generalKeys {
		if err := r.cache.Delete(ctx, key); err != nil {
			log.Printf("Failed to invalidate cache for key %s: %v", key, err)
		}
	}

	categoryPattern := "faq:category:*"
	if err := r.cache.DeletePattern(ctx, categoryPattern); err != nil {
		log.Printf("Failed to invalidate cache pattern %s: %v", categoryPattern, err)
	}
}
