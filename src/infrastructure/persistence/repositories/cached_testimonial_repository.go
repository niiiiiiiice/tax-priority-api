package repositories

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"tax-priority-api/src/application/cache"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/persistence"
)

// CachedTestimonialRepositoryImpl кешированный репозиторий отзывов
type CachedTestimonialRepositoryImpl struct {
	repo        repositories.TestimonialRepository
	cache       cache.Cache
	keys        *persistence.RedisKeys
	config      CacheConfig
	metrics     *CacheMetrics
	invalidator *CacheInvalidator
}

// NewCachedTestimonialRepositoryImpl создает новый кешированный репозиторий FAQ
func NewCachedTestimonialRepositoryImpl(
	repo repositories.TestimonialRepository,
	cache cache.Cache,
	keys *persistence.RedisKeys,
	config *CacheConfig,
) repositories.CachedTestimonialRepository {
	if config == nil {
		config = &CacheConfig{
			DefaultTTL:       15 * time.Minute,
			ShortTTL:         5 * time.Minute,
			LongTTL:          1 * time.Hour,
			EnableMetrics:    true,
			WarmupOnStart:    false,
			InvalidationMode: "selective",
		}
	}

	r := &CachedTestimonialRepositoryImpl{
		repo:    repo,
		cache:   cache,
		keys:    keys,
		config:  *config,
		metrics: &CacheMetrics{},
		invalidator: &CacheInvalidator{
			cache: cache,
			keys:  keys,
			mode:  config.InvalidationMode,
		},
	}

	if config.WarmupOnStart {
		go r.warmupCache(context.Background())
	}

	return r
}

// getFromCache получает данные из кеша с метриками
func (r *CachedTestimonialRepositoryImpl) getFromCache(ctx context.Context, key string, target interface{}) (bool, error) {
	err := r.cache.GetJSON(ctx, key, target)
	if err == nil {
		r.recordHit()
		return true, nil
	}

	if err == cache.ErrCacheMiss {
		r.recordMiss()
		return false, nil
	}

	r.recordError()
	return false, err
}

// setToCache устанавливает данные в кеш с обработкой ошибок
func (r *CachedTestimonialRepositoryImpl) setToCache(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	if err := r.cache.SetJSON(ctx, key, value, ttl); err != nil {
		log.Printf("Failed to cache data for key %s: %v", key, err)
		r.recordError()
	}
}

// cacheOrLoad загружает из кеша или базы данных
func (r *CachedTestimonialRepositoryImpl) cacheOrLoad(
	ctx context.Context,
	cacheKey string,
	loader func() (interface{}, error),
	ttl time.Duration,
) (interface{}, error) {
	// Проверяем кеш
	var result interface{}
	found, err := r.getFromCache(ctx, cacheKey, &result)
	if err != nil {
		log.Printf("Cache error for key %s: %v", cacheKey, err)
	}
	if found {
		return result, nil
	}

	// Загружаем из базы
	data, err := loader()
	if err != nil {
		return nil, err
	}

	// Кешируем результат
	r.setToCache(ctx, cacheKey, data, ttl)
	return data, nil
}

// Create создает новую FAQ и инвалидирует кеш
func (r *CachedTestimonialRepositoryImpl) Create(ctx context.Context, faq *entities.FAQ) error {
	err := r.repo.Create(ctx, faq)
	if err != nil {
		return err
	}

	r.invalidator.InvalidateForFAQ(ctx, faq)
	return nil
}

// CreateBatch создает несколько FAQ и инвалидирует кеш
func (r *CachedTestimonialRepositoryImpl) CreateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error) {
	result, err := r.repo.CreateBatch(ctx, faqs)
	if err != nil {
		return result, err
	}

	// Батчевая инвалидация
	r.invalidator.InvalidateBatch(ctx, faqs)
	return result, nil
}

// FindByID ищет FAQ по ID с кешированием
func (r *CachedTestimonialRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.FAQ, error) {
	cacheKey := r.keys.GetFAQByIDKey(id)

	result, err := r.cacheOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.repo.FindByID(ctx, id)
	}, r.config.DefaultTTL)

	if err != nil {
		return nil, err
	}

	return result.(*entities.FAQ), nil
}

// FindByIDs ищет FAQ по списку ID с оптимизированным кешированием
func (r *CachedTestimonialRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.FAQ, error) {
	if len(ids) == 0 {
		return []*entities.FAQ{}, nil
	}

	// Используем pipeline для батчевого получения из кеша
	type cacheResult struct {
		id  string
		faq *entities.FAQ
		err error
	}

	results := make(chan cacheResult, len(ids))
	var wg sync.WaitGroup

	// Параллельно проверяем кеш
	for _, id := range ids {
		wg.Add(1)
		go func(faqID string) {
			defer wg.Done()

			var faq entities.FAQ
			cacheKey := r.keys.GetFAQByIDKey(faqID)
			found, err := r.getFromCache(ctx, cacheKey, &faq)

			if found && err == nil {
				results <- cacheResult{id: faqID, faq: &faq}
			} else {
				results <- cacheResult{id: faqID, err: cache.ErrCacheMiss}
			}
		}(id)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Собираем результаты
	foundFAQs := make(map[string]*entities.FAQ)
	var missingIDs []string

	for result := range results {
		if result.err == nil && result.faq != nil {
			foundFAQs[result.id] = result.faq
		} else {
			missingIDs = append(missingIDs, result.id)
		}
	}

	// Загружаем отсутствующие из базы
	if len(missingIDs) > 0 {
		missingFAQs, err := r.repo.FindByIDs(ctx, missingIDs)
		if err != nil {
			return nil, err
		}

		// Кешируем полученные FAQ параллельно
		for _, faq := range missingFAQs {
			foundFAQs[faq.ID] = faq
			go func(f *entities.FAQ) {
				cacheKey := r.keys.GetFAQByIDKey(f.ID)
				r.setToCache(ctx, cacheKey, f, r.config.DefaultTTL)
			}(faq)
		}
	}

	// Восстанавливаем порядок
	result := make([]*entities.FAQ, 0, len(ids))
	for _, id := range ids {
		if faq, ok := foundFAQs[id]; ok {
			result = append(result, faq)
		}
	}

	return result, nil
}

// Update обновляет FAQ и инвалидирует кеш
func (r *CachedTestimonialRepositoryImpl) Update(ctx context.Context, faq *entities.FAQ) error {
	// Получаем старую версию для правильной инвалидации
	oldFAQ, _ := r.repo.FindByID(ctx, faq.ID)

	err := r.repo.Update(ctx, faq)
	if err != nil {
		return err
	}

	r.invalidator.InvalidateForUpdate(ctx, oldFAQ, faq)
	return nil
}

// UpdateBatch обновляет несколько FAQ и инвалидирует кеш
func (r *CachedTestimonialRepositoryImpl) UpdateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error) {
	result, err := r.repo.UpdateBatch(ctx, faqs)
	if err != nil {
		return result, err
	}

	r.invalidator.InvalidateBatch(ctx, faqs)
	return result, nil
}

// FindAll возвращает все FAQ с умным кешированием
func (r *CachedTestimonialRepositoryImpl) FindAll(ctx context.Context, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	// Генерируем ключ кеша на основе параметров запроса
	cacheKey := r.generateCacheKey("faq:all", opts)

	// Определяем TTL на основе сложности запроса
	ttl := r.config.DefaultTTL
	if opts != nil && (opts.Filters != nil || opts.Pagination != nil) {
		ttl = r.config.ShortTTL // Меньший TTL для фильтрованных запросов
	}

	result, err := r.cacheOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.repo.FindAll(ctx, opts)
	}, ttl)

	if err != nil {
		return nil, err
	}

	return result.([]*entities.FAQ), nil
}

// FindActive ищет активные FAQ с кешированием
func (r *CachedTestimonialRepositoryImpl) FindActive(ctx context.Context, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	cacheKey := r.generateCacheKey(r.keys.FAQActive, opts)

	result, err := r.cacheOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.repo.FindActive(ctx, opts)
	}, r.config.DefaultTTL)

	if err != nil {
		return nil, err
	}

	return result.([]*entities.FAQ), nil
}

// GetCategories возвращает список категорий с долгим кешированием
func (r *CachedTestimonialRepositoryImpl) GetCategories(ctx context.Context) ([]string, error) {
	cacheKey := r.keys.FAQCategories

	result, err := r.cacheOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.repo.GetCategories(ctx)
	}, r.config.LongTTL) // Категории меняются редко

	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

// Методы метрик

func (r *CachedTestimonialRepositoryImpl) recordHit() {
	if r.config.EnableMetrics {
		r.metrics.mu.Lock()
		r.metrics.hits++
		r.metrics.mu.Unlock()
	}
}

func (r *CachedTestimonialRepositoryImpl) recordMiss() {
	if r.config.EnableMetrics {
		r.metrics.mu.Lock()
		r.metrics.misses++
		r.metrics.mu.Unlock()
	}
}

func (r *CachedTestimonialRepositoryImpl) recordError() {
	if r.config.EnableMetrics {
		r.metrics.mu.Lock()
		r.metrics.errors++
		r.metrics.mu.Unlock()
	}
}

// GetMetrics возвращает метрики кеша
func (r *CachedTestimonialRepositoryImpl) GetMetrics() map[string]int64 {
	r.metrics.mu.RLock()
	defer r.metrics.mu.RUnlock()

	total := r.metrics.hits + r.metrics.misses
	hitRate := int64(0)
	if total > 0 {
		hitRate = (r.metrics.hits * 100) / total
	}

	return map[string]int64{
		"hits":     r.metrics.hits,
		"misses":   r.metrics.misses,
		"errors":   r.metrics.errors,
		"total":    total,
		"hit_rate": hitRate,
	}
}

// warmupCache прогревает кеш при старте
func (r *CachedTestimonialRepositoryImpl) warmupCache(ctx context.Context) {
	log.Println("Starting cache warmup...")

	// Прогреваем активные FAQ
	if faqs, err := r.repo.FindActive(ctx, nil); err == nil {
		cacheKey := r.keys.FAQActive
		r.setToCache(ctx, cacheKey, faqs, r.config.DefaultTTL)

		// Кешируем каждый FAQ по ID
		for _, faq := range faqs {
			key := r.keys.GetFAQByIDKey(faq.ID)
			r.setToCache(ctx, key, faq, r.config.DefaultTTL)
		}
	}

	// Прогреваем категории
	if categories, err := r.repo.GetCategories(ctx); err == nil {
		r.setToCache(ctx, r.keys.FAQCategories, categories, r.config.LongTTL)
	}

	log.Println("Cache warmup completed")
}

// generateCacheKey генерирует ключ кеша на основе параметров запроса
func (r *CachedTestimonialRepositoryImpl) generateCacheKey(prefix string, opts *models.QueryOptions) string {
	if opts == nil {
		return prefix
	}

	// Простая генерация ключа, можно улучшить хешированием
	key := prefix

	if opts.Filters != nil && len(opts.Filters) > 0 {
		key = fmt.Sprintf("%s:filter:%v", key, opts.Filters)
	}

	if opts.Pagination != nil {
		key = fmt.Sprintf("%s:page:%d:%d", key, opts.Pagination.Offset, opts.Pagination.Limit)
	}

	if opts.SortBy != nil && len(opts.SortBy) > 0 {
		key = fmt.Sprintf("%s:sort:%v", key, opts.SortBy)
	}

	return key
}

// Delete удаляет FAQ и инвалидирует кеш
func (r *CachedTestimonialRepositoryImpl) Delete(ctx context.Context, id string) error {
	// Получаем FAQ перед удалением для правильной инвалидации
	faq, _ := r.repo.FindByID(ctx, id)

	err := r.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if faq != nil {
		r.invalidator.InvalidateForFAQ(ctx, faq)
	} else {
		// Если не смогли получить FAQ, инвалидируем по ID
		r.cache.Delete(ctx, r.keys.GetFAQByIDKey(id))
		r.invalidator.invalidateAggregates(ctx)
	}

	return nil
}

// Clear очищает весь кеш
func (r *CachedTestimonialRepositoryImpl) Clear(ctx context.Context) error {
	r.invalidator.invalidateAll(ctx)

	// Сбрасываем метрики
	if r.config.EnableMetrics {
		r.metrics.mu.Lock()
		r.metrics.hits = 0
		r.metrics.misses = 0
		r.metrics.errors = 0
		r.metrics.mu.Unlock()
	}

	return nil
}
