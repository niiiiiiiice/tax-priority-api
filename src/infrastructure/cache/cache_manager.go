package cache

import (
	"context"
	"sync"
	appCache "tax-priority-api/src/application/cache"
	"time"
)

type CacheManager[T any, ID comparable] interface {
	Get(ctx context.Context, id ID) (T, error)
	GetMultiple(ctx context.Context, ids []ID, loader func([]ID) (map[ID]T, error)) ([]T, error)
	GetQuery(ctx context.Context, queryKey string, loader func() (interface{}, error), ttl time.Duration) (interface{}, error)
	GetOrLoad(ctx context.Context, id ID, loader func() (T, error)) (T, error)
	Set(ctx context.Context, entity T, ttl time.Duration) error
	Invalidate(ctx context.Context, entity T) error
	InvalidateMultiple(ctx context.Context, entities []T) error
	InvalidateByID(ctx context.Context, id ID) error
	InvalidateQuery(ctx context.Context, queryKey string) error
	InvalidatePattern(ctx context.Context, pattern string) error
	InvalidateAll(ctx context.Context) error
}

type DefaultCacheManager[T any, ID comparable] struct {
	cache       appCache.Cache
	stats       appCache.StatsCollector
	invalidator appCache.Invalidator[T, ID]
	keyGen      appCache.KeyGenerator[T, ID]
	config      *appCache.CacheConfig
}

func NewCacheManager[T any, ID comparable](
	cache appCache.Cache,
	keyGen appCache.KeyGenerator[T, ID],
	cacheConfig *appCache.CacheConfig,
	invalidationConfig *appCache.InvalidationConfig,
) CacheManager[T, ID] {
	stats := appCache.NewStatsCollector(cacheConfig.EnableStatistics)
	invalidator := appCache.NewInvalidator(cache, keyGen, invalidationConfig)

	return &DefaultCacheManager[T, ID]{
		cache:       cache,
		stats:       stats,
		invalidator: invalidator,
		keyGen:      keyGen,
		config:      cacheConfig,
	}
}

func (m *DefaultCacheManager[T, ID]) Get(ctx context.Context, id ID) (T, error) {
	var result T
	key := m.keyGen.GenerateKeyByID(id)

	err := m.cache.GetJSON(ctx, key, &result)
	if err != nil {
		m.stats.RecordMiss()
		return result, err
	}

	m.stats.RecordHit()
	return result, nil
}

func (m *DefaultCacheManager[T, ID]) GetMultiple(
	ctx context.Context,
	ids []ID,
	loader func([]ID) (map[ID]T, error),
) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}

	type cacheResult struct {
		id     ID
		entity T
		found  bool
	}

	results := make(chan cacheResult, len(ids))
	var wg sync.WaitGroup

	for _, id := range ids {
		wg.Add(1)
		go func(entityID ID) {
			defer wg.Done()

			entity, err := m.Get(ctx, entityID)
			if err == nil {
				results <- cacheResult{id: entityID, entity: entity, found: true}
			} else {
				var zero T
				results <- cacheResult{id: entityID, entity: zero, found: false}
			}
		}(id)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	foundEntities := make(map[ID]T)
	var missingIDs []ID

	for result := range results {
		if result.found {
			foundEntities[result.id] = result.entity
		} else {
			missingIDs = append(missingIDs, result.id)
		}
	}

	if len(missingIDs) > 0 {
		loadedEntities, err := loader(missingIDs)
		if err != nil {
			return nil, err
		}

		for id, entity := range loadedEntities {
			foundEntities[id] = entity
			go func(e T) {
				_ = m.Set(ctx, e, m.config.DefaultTTL)
			}(entity)
		}
	}

	result := make([]T, 0, len(ids))
	for _, id := range ids {
		if entity, ok := foundEntities[id]; ok {
			result = append(result, entity)
		}
	}

	return result, nil
}

func (m *DefaultCacheManager[T, ID]) GetQuery(
	ctx context.Context,
	queryKey string,
	loader func() (interface{}, error),
	ttl time.Duration,
) (interface{}, error) {
	if !m.config.Enabled {
		return loader()
	}

	cached, err := m.cache.Get(ctx, queryKey)
	if err == nil {
		m.stats.RecordHit()
		return cached, nil
	}

	m.stats.RecordMiss()

	result, err := loader()
	if err != nil {
		return nil, err
	}

	if ttl == 0 {
		ttl = m.config.DefaultTTL
	}
	_ = m.cache.SetJSON(ctx, queryKey, result, ttl)

	return result, nil
}

func (m *DefaultCacheManager[T, ID]) GetOrLoad(
	ctx context.Context,
	id ID,
	loader func() (T, error),
) (T, error) {
	result, err := m.Get(ctx, id)
	if err == nil {
		return result, nil
	}

	result, err = loader()
	if err != nil {
		return result, err
	}

	_ = m.Set(ctx, result, m.config.DefaultTTL)

	return result, nil
}

func (m *DefaultCacheManager[T, ID]) Set(ctx context.Context, entity T, ttl time.Duration) error {
	if !m.config.Enabled {
		return nil
	}

	key := m.keyGen.GenerateKey(entity)
	if ttl == 0 {
		ttl = m.config.DefaultTTL
	}

	err := m.cache.SetJSON(ctx, key, entity, ttl)
	if err != nil {
		m.stats.RecordError()
		return err
	}

	m.stats.RecordSet()
	return nil
}

func (m *DefaultCacheManager[T, ID]) Invalidate(ctx context.Context, entity T) error {
	return m.invalidator.InvalidateEntity(ctx, entity)
}

func (m *DefaultCacheManager[T, ID]) InvalidateMultiple(ctx context.Context, entities []T) error {
	return m.invalidator.InvalidateBatch(ctx, entities)
}

func (m *DefaultCacheManager[T, ID]) InvalidateByID(ctx context.Context, id ID) error {
	return m.invalidator.InvalidateByID(ctx, id)
}

func (m *DefaultCacheManager[T, ID]) InvalidateQuery(ctx context.Context, queryKey string) error {
	m.stats.RecordDelete()
	return m.cache.Delete(ctx, queryKey)
}

func (m *DefaultCacheManager[T, ID]) InvalidatePattern(ctx context.Context, pattern string) error {
	return m.cache.DeletePattern(ctx, pattern)
}

func (m *DefaultCacheManager[T, ID]) InvalidateAll(ctx context.Context) error {
	return m.invalidator.InvalidateAll(ctx)
}
