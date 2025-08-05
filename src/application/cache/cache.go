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
	DefaultTTL       time.Duration
	ShortTTL         time.Duration
	LongTTL          time.Duration
	Enabled          bool
	EnableStatistics bool
	WarmupOnStart    bool
}

func NewCacheConfig() *CacheConfig {
	return &CacheConfig{
		DefaultTTL:       30 * time.Minute,
		ShortTTL:         15 * time.Minute,
		LongTTL:          1 * time.Hour,
		Enabled:          true,
		EnableStatistics: true,
		WarmupOnStart:    true,
	}
}

type Error struct {
	Operation string
	Key       string
	Err       error
}

func (e *Error) Error() string {
	return fmt.Sprintf("cache %s failed for key '%s': %v", e.Operation, e.Key, e.Err)
}

// NewCacheError создает новую ошибку кеширования
func NewCacheError(operation, key string, err error) *Error {
	return &Error{
		Operation: operation,
		Key:       key,
		Err:       err,
	}
}
