package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"tax-priority-api/src/application/cache"

	"github.com/redis/go-redis/v9"
)

// RedisCache реализация кеша для Redis
type RedisCache struct {
	client *redis.Client
	config *cache.CacheConfig
	stats  *cache.CacheStats
}

// NewRedisCache создает новый экземпляр Redis кеша
func NewRedisCache(client *redis.Client, config *cache.CacheConfig) cache.Cache {
	return &RedisCache{
		client: client,
		config: config,
		stats: &cache.CacheStats{
			LastUpdated: time.Now(),
		},
	}
}

// Set сохраняет значение в кеш
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !r.config.Enabled {
		return nil
	}

	if ttl == 0 {
		ttl = r.config.DefaultTTL
	}

	err := r.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError("set", key, err)
	}

	atomic.AddInt64(&r.stats.Sets, 1)
	return nil
}

// Get получает значение из кеша
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if !r.config.Enabled {
		return "", cache.NewCacheError("get", key, fmt.Errorf("cache disabled"))
	}

	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			atomic.AddInt64(&r.stats.Misses, 1)
			return "", cache.NewCacheError("get", key, fmt.Errorf("key not found"))
		}
		atomic.AddInt64(&r.stats.Errors, 1)
		return "", cache.NewCacheError("get", key, err)
	}

	atomic.AddInt64(&r.stats.Hits, 1)
	return value, nil
}

// GetJSON получает JSON значение из кеша и десериализует его
func (r *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	if !r.config.Enabled {
		return cache.NewCacheError("getjson", key, fmt.Errorf("cache disabled"))
	}

	value, err := r.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError("getjson", key, fmt.Errorf("failed to unmarshal JSON: %w", err))
	}

	return nil
}

// SetJSON сериализует объект в JSON и сохраняет в кеш
func (r *RedisCache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if !r.config.Enabled {
		return nil
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError("setjson", key, fmt.Errorf("failed to marshal JSON: %w", err))
	}

	return r.Set(ctx, key, jsonData, ttl)
}

// Delete удаляет значение из кеша
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if !r.config.Enabled {
		return nil
	}

	err := r.client.Del(ctx, key).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError("delete", key, err)
	}

	atomic.AddInt64(&r.stats.Deletes, 1)
	return nil
}

// DeletePattern удаляет все ключи по паттерну
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	if !r.config.Enabled {
		return nil
	}

	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError("deletepattern", pattern, err)
	}

	if len(keys) == 0 {
		return nil
	}

	err = r.client.Del(ctx, keys...).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError("deletepattern", pattern, err)
	}

	atomic.AddInt64(&r.stats.Deletes, int64(len(keys)))
	return nil
}

// Exists проверяет существование ключа
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	if !r.config.Enabled {
		return false, nil
	}

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return false, cache.NewCacheError("exists", key, err)
	}

	return exists > 0, nil
}

// SetNX устанавливает значение только если ключ не существует
func (r *RedisCache) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	if !r.config.Enabled {
		return false, nil
	}

	if ttl == 0 {
		ttl = r.config.DefaultTTL
	}

	success, err := r.client.SetNX(ctx, key, value, ttl).Result()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return false, cache.NewCacheError("setnx", key, err)
	}

	if success {
		atomic.AddInt64(&r.stats.Sets, 1)
	}

	return success, nil
}

// Expire устанавливает TTL для существующего ключа
func (r *RedisCache) Expire(ctx context.Context, key string, ttl time.Duration) error {
	if !r.config.Enabled {
		return nil
	}

	err := r.client.Expire(ctx, key, ttl).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError("expire", key, err)
	}

	return nil
}

// TTL получает оставшееся время жизни ключа
func (r *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	if !r.config.Enabled {
		return 0, cache.NewCacheError("ttl", key, fmt.Errorf("cache disabled"))
	}

	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return 0, cache.NewCacheError("ttl", key, err)
	}

	return ttl, nil
}

// Clear очищает весь кеш
func (r *RedisCache) Clear(ctx context.Context) error {
	if !r.config.Enabled {
		return nil
	}

	err := r.client.FlushDB(ctx).Err()
	if err != nil {
		atomic.AddInt64(&r.stats.Errors, 1)
		return cache.NewCacheError("clear", "all", err)
	}

	return nil
}

// Close закрывает соединение с кешем
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// GetStats возвращает статистику кеша
func (r *RedisCache) GetStats() *cache.CacheStats {
	r.stats.LastUpdated = time.Now()
	return r.stats
}
