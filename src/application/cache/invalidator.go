package cache

import (
	"context"
	"fmt"
)

type InvalidationConfig struct {
	Mode              InvalidationMode
	BatchSize         int
	InvalidateRelated bool
}

type InvalidationMode string

const (
	InvalidationModeAggressive InvalidationMode = "aggressive"
	InvalidationModeSelective  InvalidationMode = "selective"
)

type Invalidator[T any, ID comparable] interface {
	InvalidateEntity(ctx context.Context, entity T) error
	InvalidateBatch(ctx context.Context, entities []T) error
	InvalidateByID(ctx context.Context, id ID) error
	InvalidateAll(ctx context.Context) error
}

type InvalidationStrategy[T any, ID comparable] interface {
	Execute(ctx context.Context, cache Cache, keyGen KeyGenerator[T, ID], entities []T) error
}

type AggressiveInvalidation[T any, ID comparable] struct{}

func (s *AggressiveInvalidation[T, ID]) Execute(
	ctx context.Context,
	cache Cache,
	_ KeyGenerator[T, ID],
	_ []T,
) error {
	return cache.Clear(ctx)
}

type SelectiveInvalidation[T any, ID comparable] struct {
	relatedPatterns []string
}

func (s *SelectiveInvalidation[T, ID]) Execute(
	ctx context.Context,
	cache Cache,
	keyGen KeyGenerator[T, ID],
	entities []T,
) error {

	for _, entity := range entities {
		key := keyGen.GenerateKey(entity)
		if err := cache.Delete(ctx, key); err != nil {
			return err
		}
	}

	for _, pattern := range s.relatedPatterns {
		if err := cache.DeletePattern(ctx, pattern); err != nil {
			return err
		}
	}

	return nil
}

type StrategyInvalidator[T any, ID comparable] struct {
	cache    Cache
	keyGen   KeyGenerator[T, ID]
	strategy InvalidationStrategy[T, ID]
}

func NewInvalidator[T any, ID comparable](
	cache Cache,
	keyGen KeyGenerator[T, ID],
	config *InvalidationConfig,
) Invalidator[T, ID] {
	var strategy InvalidationStrategy[T, ID]

	switch config.Mode {
	case InvalidationModeAggressive:
		strategy = &AggressiveInvalidation[T, ID]{}
	case InvalidationModeSelective:
		strategy = &SelectiveInvalidation[T, ID]{
			relatedPatterns: []string{fmt.Sprintf("%s:*", keyGen.GetPrefix())},
		}
	default:
		strategy = &SelectiveInvalidation[T, ID]{}
	}

	return &StrategyInvalidator[T, ID]{
		cache:    cache,
		keyGen:   keyGen,
		strategy: strategy,
	}
}

func (i *StrategyInvalidator[T, ID]) InvalidateEntity(ctx context.Context, entity T) error {
	return i.strategy.Execute(ctx, i.cache, i.keyGen, []T{entity})
}

func (i *StrategyInvalidator[T, ID]) InvalidateBatch(ctx context.Context, entities []T) error {
	return i.strategy.Execute(ctx, i.cache, i.keyGen, entities)
}

func (i *StrategyInvalidator[T, ID]) InvalidateByID(ctx context.Context, id ID) error {
	key := i.keyGen.GenerateKeyByID(id)
	return i.cache.Delete(ctx, key)
}

func (i *StrategyInvalidator[T, ID]) InvalidateAll(ctx context.Context) error {
	return i.cache.Clear(ctx)
}
