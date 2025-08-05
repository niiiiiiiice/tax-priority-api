package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"tax-priority-api/src/application/cache"

	"github.com/redis/go-redis/v9"
)

// Cache реализация кеша для Redis
type Cache struct {
	client *redis.Client
	config *cache.CacheConfig
	stats  *cache.Stats
}

// NewRedisCache создает новый экземпляр Redis кеша
func NewRedisCache(client *redis.Client, config *cache.CacheConfig) cache.Cache {
	return &Cache{
		client: client,
		config: config,
		stats: &cache.Stats{
			LastUpdated: time.Now(),
		},
	}
}

// Set сохраняет значение в кеш
func (r *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !r.config.Enabled {
		return nil
	}

	if ttl == 0 {
		ttl = r.config.DefaultTTL
	}

	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError(cache.Set, key, err)
	}

	atomic.AddInt64(&r.stats.Sets, 1)
	return nil
}

// Get получает значение из кеша
func (r *Cache) Get(ctx context.Context, key string) (string, error) {
	if !r.config.Enabled {
		return "", cache.NewCacheError(cache.Get, key, fmt.Errorf("cache disabled"))
	}

	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			atomic.AddInt64(&r.stats.Misses, 1)
			return "", cache.NewCacheError(cache.Get, key, fmt.Errorf("key not found"))
		}
		atomic.AddInt64(&r.stats.Errors, 1)
		return "", cache.NewCacheError(cache.Get, key, err)
	}

	atomic.AddInt64(&r.stats.Hits, 1)
	return value, nil
}

// GetJSON получает JSON значение из кеша и десериализует его
func (r *Cache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	if !r.config.Enabled {
		return cache.NewCacheError(cache.GetJSON, key, fmt.Errorf("cache disabled"))
	}

	value, err := r.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError(cache.GetJSON, key, fmt.Errorf("failed to unmarshal JSON: %w", err))
	}

	return nil
}

// SetJSON сериализует объект в JSON и сохраняет в кеш
func (r *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !r.config.Enabled {
		return nil
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError(cache.SetJSON, key, fmt.Errorf("failed to marshal JSON: %w", err))
	}

	return r.Set(ctx, key, jsonData, ttl)
}

// Delete удаляет значение из кеша
func (r *Cache) Delete(ctx context.Context, key string) error {
	if !r.config.Enabled {
		return nil
	}

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError(cache.Delete, key, err)
	}

	atomic.AddInt64(&r.stats.Deletes, 1)
	return nil
}

// DeletePattern удаляет все ключи по паттерну
func (r *Cache) DeletePattern(ctx context.Context, pattern string) error {
	if !r.config.Enabled {
		return nil
	}

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError(cache.DeletePattern, pattern, err)
	}

	if len(keys) == 0 {
		return nil
	}

	err = r.client.Del(ctx, keys...).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError(cache.DeletePattern, pattern, err)
	}

	atomic.AddInt64(&r.stats.Deletes, int64(len(keys)))
	return nil
}

// Exists проверяет существование ключа
func (r *Cache) Exists(ctx context.Context, key string) (bool, error) {
	if !r.config.Enabled {
		return false, nil
	}

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return false, cache.NewCacheError(cache.Exists, key, err)
	}

	return exists > 0, nil
}

func (r *Cache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	if !r.config.Enabled {
		return false, nil
	}

	if ttl == 0 {
		ttl = r.config.DefaultTTL
	}

	success, err := r.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return false, cache.NewCacheError(cache.SetNX, key, err)
	}

	if success {
		atomic.AddInt64(&r.stats.Sets, 1)
	}

	return success, nil
}

func (r *Cache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if !r.config.Enabled {
		return nil
	}

	err := r.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError(cache.Expire, key, err)
	}

	return nil
}

func (r *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	if !r.config.Enabled {
		return 0, cache.NewCacheError(cache.TTL, key, fmt.Errorf("cache disabled"))
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return 0, cache.NewCacheError(cache.TTL, key, err)
	}

	return ttl, nil
}

func (r *Cache) Clear(ctx context.Context) error {
	if !r.config.Enabled {
		return nil
	}

	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError(cache.Clear, "all", err)
	}

	return nil
}

func (r *Cache) Close() error {
	return r.client.Close()
}

func (r *Cache) GetStats() *cache.Stats {
	r.stats.LastUpdated = time.Now()
	return r.stats
}
