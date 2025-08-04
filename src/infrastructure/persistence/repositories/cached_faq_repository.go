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

// CachedFAQRepositoryImpl кешированный репозиторий FAQ
type CachedFAQRepositoryImpl struct {
	repo repositories.FAQRepository
}

// NewCachedFAQRepository создает новый кешированный репозиторий FAQ
func NewCachedFAQRepository(
	repo repositories.FAQRepository,
	cache cache.Cache,
	keys *persistence.RedisKeys,
	config *CacheConfig,
) repositories.CachedFAQRepository {

	r := &CachedFAQRepositoryImpl{
		repo: repo,
	}

	if config.WarmupOnStart {
		go r.warmupCache(context.Background())
	}

	return r
}

func (cr *CachedFAQRepositoryImpl) getFromCache(ctx context.Context, key string, target interface{}) (bool, error) {
	err := cr.cache.GetJSON(ctx, key, target)
	if err == nil {
		cr.recordHit()
		return true, nil
	}

	if err == cache.ErrCacheMiss {
		cr.recordMiss()
		return false, nil
	}

	cr.recordError()
	return false, err
}

// setToCache устанавливает данные в кеш с обработкой ошибок
func (cr *CachedFAQRepositoryImpl) setToCache(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	if err := cr.cache.SetJSON(ctx, key, value, ttl); err != nil {
		log.Printf("Failed to cache data for key %s: %v", key, err)
		cr.recordError()
	}
}

// cacheOrLoad загружает из кеша или базы данных
func (cr *CachedFAQRepositoryImpl) cacheOrLoad(
	ctx context.Context,
	cacheKey string,
	loader func() (interface{}, error),
	ttl time.Duration,
) (interface{}, error) {
	// Проверяем кеш
	var result interface{}
	found, err := cr.getFromCache(ctx, cacheKey, &result)
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
	cr.setToCache(ctx, cacheKey, data, ttl)
	return data, nil
}

func (cr *CachedFAQRepositoryImpl) Create(ctx context.Context, faq *entities.FAQ) error {
	return cr.repo.Create(ctx, faq)
}

func (cr *CachedFAQRepositoryImpl) CreateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error) {
	return cr.repo.CreateBatch(ctx, faqs)
}

func (cr *CachedFAQRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.FAQ, error) {
	return cr.repo.FindByID(ctx, id)
}

func (cr *CachedFAQRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.FAQ, error) {
	return cr.repo.FindByIDs(ctx, ids)
}

func (cr *CachedFAQRepositoryImpl) Update(ctx context.Context, faq *entities.FAQ) error {
	return cr.repo.Update(ctx, faq)
}

func (cr *CachedFAQRepositoryImpl) UpdateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error) {
	return cr.repo.UpdateBatch(ctx, faqs)
}

// Методы метрик

func (cr *CachedFAQRepositoryImpl) recordHit() {
	if cr.config.EnableMetrics {
		cr.metrics.mu.Lock()
		cr.metrics.hits++
		cr.metrics.mu.Unlock()
	}
}

func (cr *CachedFAQRepositoryImpl) recordMiss() {
	if cr.config.EnableMetrics {
		cr.metrics.mu.Lock()
		cr.metrics.misses++
		cr.metrics.mu.Unlock()
	}
}

func (cr *CachedFAQRepositoryImpl) recordError() {
	if cr.config.EnableMetrics {
		cr.metrics.mu.Lock()
		cr.metrics.errors++
		cr.metrics.mu.Unlock()
	}
}

// GetMetrics возвращает метрики кеша
func (cr *CachedFAQRepositoryImpl) GetMetrics() map[string]int64 {
	cr.metrics.mu.RLock()
	defer cr.metrics.mu.RUnlock()

	total := cr.metrics.hits + cr.metrics.misses
	hitRate := int64(0)
	if total > 0 {
		hitRate = (cr.metrics.hits * 100) / total
	}

	return map[string]int64{
		"hits":     cr.metrics.hits,
		"misses":   cr.metrics.misses,
		"errors":   cr.metrics.errors,
		"total":    total,
		"hit_rate": hitRate,
	}
}

// warmupCache прогревает кеш при старте
func (cr *CachedFAQRepositoryImpl) warmupCache(ctx context.Context) {
	log.Println("Starting cache warmup...")

	// Прогреваем активные FAQ
	if faqs, err := cr.repo.FindActive(ctx, nil); err == nil {
		cacheKey := cr.keys.FAQActive
		cr.setToCache(ctx, cacheKey, faqs, cr.config.DefaultTTL)

		// Кешируем каждый FAQ по ID
		for _, faq := range faqs {
			key := cr.keys.GetFAQByIDKey(faq.ID)
			cr.setToCache(ctx, key, faq, cr.config.DefaultTTL)
		}
	}

	// Прогреваем категории
	if categories, err := cr.repo.GetCategories(ctx); err == nil {
		cr.setToCache(ctx, cr.keys.FAQCategories, categories, cr.config.LongTTL)
	}

	log.Println("Cache warmup completed")
}

// generateCacheKey генерирует ключ кеша на основе параметров запроса
func (cr *CachedFAQRepositoryImpl) generateCacheKey(prefix string, opts *models.QueryOptions) string {
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

// InvalidateForFAQ инвалидирует кеш для конкретного FAQ
func (ci *CacheInvalidator) InvalidateForFAQ(ctx context.Context, faq *entities.FAQ) {
	if faq == nil {
		return
	}

	// Всегда удаляем кеш по ID
	if err := ci.cache.Delete(ctx, ci.keys.GetFAQByIDKey(faq.ID)); err != nil {
		log.Printf("Failed to invalidate cache for FAQ ID %s: %v", faq.ID, err)
	}

	// В зависимости от режима инвалидации
	if ci.mode == "aggressive" {
		ci.invalidateAll(ctx)
	} else {
		ci.invalidateSelective(ctx, faq)
	}
}

// InvalidateForUpdate инвалидирует кеш при обновлении
func (ci *CacheInvalidator) InvalidateForUpdate(ctx context.Context, oldFAQ, newFAQ *entities.FAQ) {
	// Удаляем кеш по ID
	if newFAQ != nil {
		if err := ci.cache.Delete(ctx, ci.keys.GetFAQByIDKey(newFAQ.ID)); err != nil {
			log.Printf("Failed to invalidate cache for FAQ ID %s: %v", newFAQ.ID, err)
		}
	}

	// Если изменилась категория, инвалидируем обе
	if oldFAQ != nil && newFAQ != nil && oldFAQ.Category != newFAQ.Category {
		ci.cache.Delete(ctx, ci.keys.GetFAQByCategoryKey(oldFAQ.Category))
		ci.cache.Delete(ctx, ci.keys.GetFAQByCategoryKey(newFAQ.Category))
	} else if newFAQ != nil {
		ci.cache.Delete(ctx, ci.keys.GetFAQByCategoryKey(newFAQ.Category))
	}

	// Инвалидируем агрегированные данные
	ci.invalidateAggregates(ctx)
}

// InvalidateBatch инвалидирует кеш для множества FAQ
func (ci *CacheInvalidator) InvalidateBatch(ctx context.Context, faqs []*entities.FAQ) {
	// Используем pipeline для батчевого удаления
	keys := make([]string, 0, len(faqs)*2)
	categories := make(map[string]bool)

	for _, faq := range faqs {
		if faq != nil {
			keys = append(keys, ci.keys.GetFAQByIDKey(faq.ID))
			categories[faq.Category] = true
		}
	}

	// Добавляем ключи категорий
	for category := range categories {
		keys = append(keys, ci.keys.GetFAQByCategoryKey(category))
	}

	// Батчевое удаление
	if len(keys) > 0 {
		if err := ci.cache.DeleteBatch(ctx, keys); err != nil {
			log.Printf("Failed to batch invalidate cache: %v", err)
		}
	}

	ci.invalidateAggregates(ctx)
}

func (ci *CacheInvalidator) invalidateSelective(ctx context.Context, faq *entities.FAQ) {
	// Удаляем только связанные кеши
	keys := []string{
		ci.keys.GetFAQByCategoryKey(faq.Category),
		ci.keys.FAQActive,
		"faq:all",
	}

	for _, key := range keys {
		ci.cache.Delete(ctx, key)
	}

	// Инвалидируем кешированные результаты с этой категорией
	pattern := fmt.Sprintf("faq:all:filter:*category:%s*", faq.Category)
	ci.cache.DeletePattern(ctx, pattern)
}

func (ci *CacheInvalidator) invalidateAll(ctx context.Context) {
	patterns := []string{
		"faq:*",
		ci.keys.FAQActive,
		ci.keys.FAQCount,
		ci.keys.FAQCategories,
	}

	for _, pattern := range patterns {
		if err := ci.cache.DeletePattern(ctx, pattern); err != nil {
			log.Printf("Failed to clear cache pattern %s: %v", pattern, err)
		}
	}
}

func (ci *CacheInvalidator) invalidateAggregates(ctx context.Context) {
	keys := []string{
		ci.keys.FAQCount,
		ci.keys.FAQCategories,
		ci.keys.FAQActive,
		"faq:all",
	}

	for _, key := range keys {
		ci.cache.Delete(ctx, key)
	}
}
