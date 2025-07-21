package cache

import (
	"context"
	"fmt"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	Get(ctx context.Context, key string) (string, error)

	GetJSON(ctx context.Context, key string, dest interface{}) error

	SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error

	Delete(ctx context.Context, key string) error

	DeletePattern(ctx context.Context, pattern string) error

	Exists(ctx context.Context, key string) (bool, error)

	SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error)

	Expire(ctx context.Context, key string, ttl time.Duration) error

	TTL(ctx context.Context, key string) (time.Duration, error)

	Clear(ctx context.Context) error

	Close() error
}

type CacheConfig struct {
	DefaultTTL time.Duration
	Enabled    bool
}

func NewCacheConfig() *CacheConfig {
	return &CacheConfig{
		DefaultTTL: 15 * time.Minute,
		Enabled:    true,
	}
}

type CacheError struct {
	Operation string
	Key       string
	Err       error
}

func (e *CacheError) Error() string {
	return fmt.Sprintf("cache %s failed for key '%s': %v", e.Operation, e.Key, e.Err)
}

// NewCacheError создает новую ошибку кеширования
func NewCacheError(operation, key string, err error) *CacheError {
	return &CacheError{
		Operation: operation,
		Key:       key,
		Err:       err,
	}
}

// CacheStats статистика кеша
type CacheStats struct {
	Hits        int64
	Misses      int64
	Sets        int64
	Deletes     int64
	Errors      int64
	LastUpdated time.Time
}

// HitRate возвращает процент попаданий в кеш
func (s *CacheStats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total) * 100
}
