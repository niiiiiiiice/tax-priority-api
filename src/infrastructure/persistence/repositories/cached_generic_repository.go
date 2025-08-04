package repositories

import (
	"context"
	"sync"
	"tax-priority-api/src/application/cache"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/persistence"
	"time"
)

// CacheConfig конфигурация кеширования
type CacheConfig struct {
	DefaultTTL       time.Duration
	ShortTTL         time.Duration
	LongTTL          time.Duration
	EnableMetrics    bool
	WarmupOnStart    bool
	InvalidationMode string // "aggressive" | "selective"
}

// CacheMetrics метрики кеша
type CacheMetrics struct {
	hits   int64
	misses int64
	errors int64
	mu     sync.RWMutex
}

// CacheInvalidator управляет инвалидацией кеша
type CacheInvalidator struct {
	cache cache.Cache
	keys  *persistence.RedisKeys
	mode  string
}

type CachedGenericRepositoryImpl[T entities.Entity[ID], ID comparable] struct {
	generic     repositories.GenericRepository[T, ID]
	cache       cache.Cache
	keys        *persistence.RedisKeys
	config      CacheConfig
	metrics     *CacheMetrics
	invalidator *CacheInvalidator
}

func NewCachedGenericRepository[T entities.Entity[ID], ID comparable](
	generic repositories.GenericRepository[T, ID],
	cache cache.Cache,
	keys *persistence.RedisKeys,
	config *CacheConfig) repositories.GenericRepository[T, ID] {
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

	return &CachedGenericRepositoryImpl[T, ID]{
		generic: generic,
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
}

func (repo CachedGenericRepositoryImpl[T, ID]) Create(ctx context.Context, entity T) error {
	err := repo.Create(ctx, entity)
	if err != nil {
		return err
	}

	repo.invalidator.InvalidateForFAQ(ctx, entity)
	return nil
}

func (repo CachedGenericRepositoryImpl[T, ID]) CreateBatch(ctx context.Context, entities []T) (*models.BulkOperationResult, error) {
	result, err := repo.generic.CreateBatch(ctx, entities)
	if err != nil {
		return result, err
	}

	repo.invalidator.InvalidateBatch(ctx, entities)
	return result, nil
}

func (repo CachedGenericRepositoryImpl[T, ID]) FindByID(ctx context.Context, id ID) (T, error) {
	cacheKey := repo.keys.GetFAQByIDKey(id)

	result, err := repo.cacheOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.repo.FindByID(ctx, id)
	}, r.config.DefaultTTL)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo CachedGenericRepositoryImpl[T, ID]) FindByIDs(ctx context.Context, ids []ID) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}

	type cacheResult struct {
		id  string
		faq *entities.FAQ
		err error
	}

	results := make(chan cacheResult, len(ids))
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)
		go func(faqID string) {
			defer wg.Done()

			var faq entities.FAQ
			cacheKey := cr.keys.GetFAQByIDKey(faqID)
			found, err := cr.getFromCache(ctx, cacheKey, &faq)

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
		missingFAQs, err := cr.repo.FindByIDs(ctx, missingIDs)
		if err != nil {
			return nil, err
		}

		// Кешируем полученные FAQ параллельно
		for _, faq := range missingFAQs {
			foundFAQs[faq.ID] = faq
			go func(f *entities.FAQ) {
				cacheKey := cr.keys.GetFAQByIDKey(f.ID)
				cr.setToCache(ctx, cacheKey, f, cr.config.DefaultTTL)
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

func (repo CachedGenericRepositoryImpl[T, ID]) Update(ctx context.Context, entity T) error {
	oldEntity, _ := repo.generic.FindByID(ctx, entity.GetID())

	err := repo.generic.Update(ctx, entity)
	if err != nil {
		return err
	}

	repo.invalidator.InvalidateForUpdate(ctx, oldEntity, entity)
	return nil
}

func (repo CachedGenericRepositoryImpl[T, ID]) UpdateBatch(ctx context.Context, entities []T) (*models.BulkOperationResult, error) {
	result, err := repo.generic.UpdateBatch(ctx, entities)
	if err != nil {
		return result, err
	}

	repo.invalidator.InvalidateBatch(ctx, entities)
	return result, nil
}

func (repo CachedGenericRepositoryImpl[T, ID]) UpdateFields(ctx context.Context, id ID, fields map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) Delete(ctx context.Context, id ID) error {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) DeleteBatch(ctx context.Context, ids []ID) (*models.BulkOperationResult, error) {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) SoftDelete(ctx context.Context, id ID) error {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) FindAll(ctx context.Context, opts *models.QueryOptions) ([]T, error) {
	cacheKey := repo.generateCacheKey("faq:all", opts)

	ttl := repo.config.DefaultTTL
	if opts != nil && (opts.Filters != nil || opts.Pagination != nil) {
		ttl = repo.config.ShortTTL // Меньший TTL для фильтрованных запросов
	}

	result, err := repo.cacheOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return repo.generic.FindAll(ctx, opts)
	}, ttl)

	if err != nil {
		return nil, err
	}

	return result.([]*entities.FAQ), nil
}

func (repo CachedGenericRepositoryImpl[T, ID]) FindOne(ctx context.Context, opts *models.QueryOptions) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[T], error) {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) Exists(ctx context.Context, id ID) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) Refresh(ctx context.Context, entity T) error {
	//TODO implement me
	panic("implement me")
}

func (repo CachedGenericRepositoryImpl[T, ID]) Clear(ctx context.Context) error {
	err := repo.generic.Clear(ctx)

	if err != nil {
		return err
	}

	repo.invalidator.invalidateAll(ctx)

	if repo.config.EnableMetrics {
		repo.metrics.mu.Lock()
		repo.metrics.hits = 0
		repo.metrics.misses = 0
		repo.metrics.errors = 0
		repo.metrics.mu.Unlock()
	}

	return nil
}
