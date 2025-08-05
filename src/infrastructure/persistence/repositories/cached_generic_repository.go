package repositories

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/cache"
)

type CachedGenericRepositoryImpl[T entities.Entity[ID], ID comparable] struct {
	genericRepo  repositories.GenericRepository[T, ID]
	cacheManager cache.CacheManager[T, ID]
}

func NewCachedGenericRepository[T entities.Entity[ID], ID comparable](
	genericRepo repositories.GenericRepository[T, ID],
	cacheManager cache.CacheManager[T, ID]) repositories.GenericRepository[T, ID] {

	return &CachedGenericRepositoryImpl[T, ID]{
		genericRepo:  genericRepo,
		cacheManager: cacheManager,
	}
}

func (r *CachedGenericRepositoryImpl[T, ID]) Create(ctx context.Context, entity T) error {
	err := r.genericRepo.Create(ctx, entity)
	if err != nil {
		return err
	}

	return r.cacheManager.Invalidate(ctx, entity)
}

func (r *CachedGenericRepositoryImpl[T, ID]) CreateBatch(ctx context.Context, entities []T) (*models.BulkOperationResult, error) {
	result, err := r.genericRepo.CreateBatch(ctx, entities)
	if err != nil {
		return result, err
	}

	_ = r.cacheManager.InvalidateMultiple(ctx, entities)
	_ = r.invalidateAggregatedQueries(ctx)

	return result, nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) FindByID(ctx context.Context, id ID) (T, error) {
	return r.cacheManager.GetOrLoad(ctx, id, func() (T, error) {
		return r.genericRepo.FindByID(ctx, id)
	})
}

func (r *CachedGenericRepositoryImpl[T, ID]) FindByIDs(ctx context.Context, ids []ID) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}

	return r.cacheManager.GetMultiple(ctx, ids, func(missingIDs []ID) (map[ID]T, error) {
		foundEntities, err := r.genericRepo.FindByIDs(ctx, missingIDs)
		if err != nil {
			return nil, err
		}

		result := make(map[ID]T)
		for _, entity := range foundEntities {
			result[entity.GetID()] = entity
		}
		return result, nil
	})
}

func (r *CachedGenericRepositoryImpl[T, ID]) Update(ctx context.Context, entity T) error {
	err := r.genericRepo.Update(ctx, entity)
	if err != nil {
		return err
	}

	return r.cacheManager.Invalidate(ctx, entity)
}

func (r *CachedGenericRepositoryImpl[T, ID]) UpdateBatch(ctx context.Context, entities []T) (*models.BulkOperationResult, error) {
	result, err := r.genericRepo.UpdateBatch(ctx, entities)
	if err != nil {
		return result, err
	}

	_ = r.cacheManager.InvalidateMultiple(ctx, entities)
	_ = r.invalidateAggregatedQueries(ctx)

	return result, nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) UpdateFields(ctx context.Context, id ID, fields map[string]interface{}) error {
	err := r.genericRepo.UpdateFields(ctx, id, fields)
	if err != nil {
		return err
	}

	_ = r.cacheManager.InvalidateByID(ctx, id)
	_ = r.invalidateAggregatedQueries(ctx)

	return nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) Delete(ctx context.Context, id ID) error {
	entity, findErr := r.FindByID(ctx, id)
	if findErr != nil {
		return findErr
	}

	delErr := r.genericRepo.Delete(ctx, id)
	if delErr != nil {
		return delErr
	}

	_ = r.cacheManager.InvalidateByID(ctx, id)
	if !isZero(entity) {
		_ = r.cacheManager.Invalidate(ctx, entity)
	}

	return nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) DeleteBatch(ctx context.Context, ids []ID) (*models.BulkOperationResult, error) {
	foundEntities, _ := r.FindByIDs(ctx, ids)

	result, err := r.genericRepo.DeleteBatch(ctx, ids)
	if err != nil {
		return result, err
	}

	for _, id := range ids {
		_ = r.cacheManager.InvalidateByID(ctx, id)
	}
	if len(foundEntities) > 0 {
		_ = r.cacheManager.InvalidateMultiple(ctx, foundEntities)
	}
	_ = r.invalidateAggregatedQueries(ctx)

	return result, nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) SoftDelete(ctx context.Context, id ID) error {
	entity, _ := r.FindByID(ctx, id)

	err := r.genericRepo.SoftDelete(ctx, id)
	if err != nil {
		return err
	}

	// Инвалидируем через менеджер
	_ = r.cacheManager.InvalidateByID(ctx, id)
	if !isZero(entity) {
		_ = r.cacheManager.Invalidate(ctx, entity)
	}

	return nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) FindAll(ctx context.Context, opts *models.QueryOptions) ([]T, error) {
	cacheKey := r.keyGen.GenerateQueryKey("all", opts)
	ttl := r.determineTTL(opts)

	// Используем QueryCacheManager для кеширования запросов
	result, err := r.cacheManager.GetOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.genericRepo.FindAll(ctx, opts)
	}, ttl)

	if err != nil {
		return nil, err
	}

	// Преобразование типа
	entities, ok := result.([]T)
	if !ok {
		return r.genericRepo.FindAll(ctx, opts)
	}

	// Кешируем отдельные сущности асинхронно
	go func() {
		for _, entity := range entities {
			_ = r.cacheManager.Set(ctx, entity, r.config.DefaultTTL)
		}
	}()

	return entities, nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) FindOne(ctx context.Context, opts *models.QueryOptions) (T, error) {
	cacheKey := r.keyGen.GenerateQueryKey("one", opts)

	result, err := r.queryCache.GetOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.genericRepo.FindOne(ctx, opts)
	}, r.config.ShortTTL)

	if err != nil {
		var zero T
		return zero, err
	}

	// Преобразование типа
	entity, ok := result.(T)
	if !ok {
		return r.genericRepo.FindOne(ctx, opts)
	}

	// Кешируем по ID
	_ = r.cacheManager.Set(ctx, entity, r.config.DefaultTTL)

	return entity, nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[T], error) {
	cacheKey := r.keyGen.GenerateQueryKey("paginated", opts)

	result, err := r.queryCache.GetOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.genericRepo.FindWithPagination(ctx, opts)
	}, r.config.ShortTTL)

	if err != nil {
		return nil, err
	}

	// Преобразование типа
	paginatedResult, ok := result.(*models.PaginatedResult[T])
	if !ok {
		return r.genericRepo.FindWithPagination(ctx, opts)
	}

	// Кешируем отдельные сущности асинхронно
	go func() {
		for _, entity := range paginatedResult.Items {
			_ = r.cacheManager.Set(ctx, entity, r.config.DefaultTTL)
		}
	}()

	return paginatedResult, nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	cacheKey := r.keyGen.GenerateQueryKey("count", filters)

	result, err := r.queryCache.GetOrLoad(ctx, cacheKey, func() (interface{}, error) {
		return r.genericRepo.Count(ctx, filters)
	}, r.config.ShortTTL)

	if err != nil {
		return 0, err
	}

	count, ok := result.(int64)
	if !ok {
		return r.genericRepo.Count(ctx, filters)
	}

	return count, nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) Exists(ctx context.Context, id ID) (bool, error) {
	_, err := r.cacheManager.Get(ctx, id)
	if err == nil {
		return true, nil
	}

	// Проверяем в базе
	return r.genericRepo.Exists(ctx, id)
}

func (r *CachedGenericRepositoryImpl[T, ID]) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	return r.genericRepo.ExistsByFields(ctx, filters)
}

func (r *CachedGenericRepositoryImpl[T, ID]) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	return r.genericRepo.WithTransaction(ctx, func(txCtx context.Context) error {
		err := fn(txCtx)
		if err != nil {
			return err
		}

		// После транзакции инвалидируем агрегированные запросы
		_ = r.invalidateAggregatedQueries(ctx)

		return nil
	})
}

func (r *CachedGenericRepositoryImpl[T, ID]) Refresh(ctx context.Context, entity T) error {
	id := entity.GetID()

	// Инвалидируем кеш через менеджер
	_ = r.cacheManager.InvalidateByID(ctx, id)

	// Обновляем из базы
	err := r.genericRepo.Refresh(ctx, entity)
	if err != nil {
		return err
	}

	// Кешируем обновленную сущность
	_ = r.cacheManager.Set(ctx, entity, r.config.DefaultTTL)

	return nil
}

func (r *CachedGenericRepositoryImpl[T, ID]) Clear(ctx context.Context) error {
	err := r.genericRepo.Clear(ctx)
	if err != nil {
		return err
	}

	// Полная очистка через менеджер
	return r.cacheManager.InvalidateAll(ctx)
}

func (r *CachedGenericRepositoryImpl[T, ID]) invalidateAggregatedQueries(ctx context.Context) error {
	// Получаем префикс из генератора ключей
	sampleKey := r.keyGen.GenerateKeyByID(*new(ID))
	prefix := sampleKey[:3] // Берем префикс (например, "use" от "user")

	patterns := []string{
		fmt.Sprintf("%s:all:*", prefix),
		fmt.Sprintf("%s:count:*", prefix),
		fmt.Sprintf("%s:paginated:*", prefix),
		fmt.Sprintf("%s:one:*", prefix),
	}

	// Делегируем инвалидацию менеджеру
	for _, pattern := range patterns {
		_ = r.queryCache.InvalidatePattern(ctx, pattern)
	}

	return nil
}

func isZero[T any](v T) bool {
	var zero T
	return fmt.Sprintf("%v", v) == fmt.Sprintf("%v", zero)
}
