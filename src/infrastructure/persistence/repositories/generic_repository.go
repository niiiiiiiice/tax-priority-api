package repositories

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	sharedModels "tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	persistence "tax-priority-api/src/infrastructure/persistence"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GenericRepositoryImpl[T entities.Entity[ID], M any, ID comparable] struct {
	db            *gorm.DB
	domainToModel func(T) *M
	modelToDomain func(*M) T
}

func NewGenericRepository[T entities.Entity[ID], M any, ID comparable](
	db *gorm.DB,
	domainToModel func(T) *M,
	modelToDomain func(*M) T,
) repositories.GenericRepository[T, ID] {
	return &GenericRepositoryImpl[T, M, ID]{
		db:            db,
		domainToModel: domainToModel,
		modelToDomain: modelToDomain,
	}
}

func (r *GenericRepositoryImpl[T, M, ID]) Create(ctx context.Context, entity T) error {
	model := r.domainToModel(entity)

	now := time.Now()
	entity.SetCreatedAt(now)
	entity.SetUpdatedAt(now)

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return persistence.NewInternalError("failed to create entity", result.Error)
	}

	return nil
}

func (r *GenericRepositoryImpl[T, M, ID]) CreateBatch(ctx context.Context, entities []T) (*sharedModels.BulkOperationResult, error) {
	if len(entities) == 0 {
		return &sharedModels.BulkOperationResult{SuccessCount: 0, FailureCount: 0}, nil
	}

	models := make([]M, len(entities))
	now := time.Now()

	for i, entity := range entities {
		entity.SetCreatedAt(now)
		entity.SetUpdatedAt(now)
		models[i] = *r.domainToModel(entity)
	}

	result := r.db.WithContext(ctx).CreateInBatches(models, 100)

	if result.Error != nil {
		return &sharedModels.BulkOperationResult{
			SuccessCount: 0,
			FailureCount: len(entities),
			Errors:       []error{result.Error},
		}, persistence.NewInternalError("failed to create entities batch", result.Error)
	}

	return &sharedModels.BulkOperationResult{
		SuccessCount: len(entities),
		FailureCount: 0,
	}, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) FindByID(ctx context.Context, id ID) (T, error) {
	var model M
	var zero T

	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return zero, persistence.NewNotFoundError(fmt.Sprintf("entity with id %v not found", id), result.Error)
		}
		return zero, persistence.NewInternalError("failed to find entity by id", result.Error)
	}

	entity := r.modelToDomain(&model) // возвращает T (то есть *entities.User)
	return entity, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) FindByIDs(ctx context.Context, ids []ID) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}

	var models []M
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find entities by ids", result.Error)
	}

	_entities := make([]T, len(models))
	for i, model := range models {
		_entities[i] = r.modelToDomain(&model)
	}

	return _entities, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) Update(ctx context.Context, entity T) error {
	model := r.domainToModel(entity)
	entity.SetUpdatedAt(time.Now())

	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return persistence.NewInternalError("failed to update entity", result.Error)
	}

	return nil
}

func (r *GenericRepositoryImpl[T, M, ID]) UpdateBatch(ctx context.Context, entities []T) (*sharedModels.BulkOperationResult, error) {
	if len(entities) == 0 {
		return &sharedModels.BulkOperationResult{SuccessCount: 0, FailureCount: 0}, nil
	}

	models := make([]*M, len(entities))
	now := time.Now()

	for i, entity := range entities {
		entity.SetUpdatedAt(now)
		models[i] = r.domainToModel(entity)
	}

	result := r.db.WithContext(ctx).Save(models)
	if result.Error != nil {
		return &sharedModels.BulkOperationResult{
			SuccessCount: 0,
			FailureCount: len(entities),
			Errors:       []error{result.Error},
		}, persistence.NewInternalError("failed to update entities batch", result.Error)
	}

	return &sharedModels.BulkOperationResult{
		SuccessCount: len(entities),
		FailureCount: 0,
	}, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) UpdateFields(ctx context.Context, id ID, fields map[string]interface{}) error {
	fields["updated_at"] = time.Now()

	result := r.db.WithContext(ctx).Model(new(M)).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return persistence.NewInternalError("failed to update entity fields", result.Error)
	}

	if result.RowsAffected == 0 {
		return persistence.NewNotFoundError(fmt.Sprintf("entity with id %v not found", id), nil)
	}

	return nil
}

func (r *GenericRepositoryImpl[T, M, ID]) Delete(ctx context.Context, id ID) error {
	result := r.db.WithContext(ctx).Delete(new(M), "id = ?", id)
	if result.Error != nil {
		return persistence.NewInternalError("failed to delete entity", result.Error)
	}

	if result.RowsAffected == 0 {
		return persistence.NewNotFoundError(fmt.Sprintf("entity with id %v not found", id), nil)
	}

	return nil
}

func (r *GenericRepositoryImpl[T, M, ID]) DeleteBatch(ctx context.Context, ids []ID) (*sharedModels.BulkOperationResult, error) {
	if len(ids) == 0 {
		return &sharedModels.BulkOperationResult{SuccessCount: 0, FailureCount: 0}, nil
	}

	result := r.db.WithContext(ctx).Delete(new(M), "id IN ?", ids)
	if result.Error != nil {
		return &sharedModels.BulkOperationResult{
			SuccessCount: 0,
			FailureCount: len(ids),
			Errors:       []error{result.Error},
		}, persistence.NewInternalError("failed to delete entities batch", result.Error)
	}

	return &sharedModels.BulkOperationResult{
		SuccessCount: int(result.RowsAffected),
		FailureCount: len(ids) - int(result.RowsAffected),
	}, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) SoftDelete(ctx context.Context, id ID) error {
	result := r.db.WithContext(ctx).Model(new(M)).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return persistence.NewInternalError("failed to soft delete entity", result.Error)
	}

	if result.RowsAffected == 0 {
		return persistence.NewNotFoundError(fmt.Sprintf("entity with id %v not found", id), nil)
	}

	return nil
}

func (r *GenericRepositoryImpl[T, M, ID]) FindAll(ctx context.Context, opts *sharedModels.QueryOptions) ([]T, error) {
	var models []M
	query := r.db.WithContext(ctx)

	// Применяем фильтры
	if opts != nil {
		query = r.applyFilters(query, opts.Filters)
		query = r.applySorting(query, opts.SortBy)
		query = r.applyIncludes(query, opts.Includes)

		if opts.Pagination != nil {
			query = query.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)
		}
	}

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find entities", result.Error)
	}

	_entities := make([]T, len(models))
	for i, model := range models {
		_entities[i] = r.modelToDomain(&model)
	}

	return _entities, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) FindOne(ctx context.Context, opts *sharedModels.QueryOptions) (T, error) {
	var model M
	var zero T

	query := r.db.WithContext(ctx)

	if opts != nil {
		query = r.applyFilters(query, opts.Filters)
		query = r.applySorting(query, opts.SortBy)
		query = r.applyIncludes(query, opts.Includes)
	}

	result := query.First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return zero, persistence.NewNotFoundError("entity not found", result.Error)
		}
		return zero, persistence.NewInternalError("failed to find entity", result.Error)
	}

	entity := r.modelToDomain(&model)
	return entity, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) FindWithPagination(ctx context.Context, opts *sharedModels.QueryOptions) (*sharedModels.PaginatedResult[T], error) {
	if opts == nil || opts.Pagination == nil {
		return nil, persistence.NewInvalidInputError("pagination options are required", nil)
	}

	var models []M
	var total int64

	countQuery := r.db.WithContext(ctx).Model(new(M))
	countQuery = r.applyFilters(countQuery, opts.Filters)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, persistence.NewInternalError("failed to count _entities", err)
	}

	query := r.db.WithContext(ctx)
	query = r.applyFilters(query, opts.Filters)
	query = r.applySorting(query, opts.SortBy)
	query = r.applyIncludes(query, opts.Includes)
	query = query.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find entities with pagination", result.Error)
	}

	_entities := make([]T, len(models))
	for i, model := range models {
		_entities[i] = r.modelToDomain(&model)
	}

	totalPages := int((total + int64(opts.Pagination.Limit) - 1) / int64(opts.Pagination.Limit))
	hasNext := opts.Pagination.Offset+opts.Pagination.Limit < int(total)
	hasPrev := opts.Pagination.Offset > 0

	return &sharedModels.PaginatedResult[T]{
		Items:      _entities,
		Total:      total,
		Offset:     opts.Pagination.Offset,
		Limit:      opts.Pagination.Limit,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		TotalPages: totalPages,
	}, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(M))

	query = r.applyFilters(query, filters)

	result := query.Count(&count)
	if result.Error != nil {
		return 0, persistence.NewInternalError("failed to count entities", result.Error)
	}

	return count, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) Exists(ctx context.Context, id ID) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(new(M)).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, persistence.NewInternalError("failed to check entity existence", result.Error)
	}

	return count > 0, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(M))

	query = r.applyFilters(query, filters)

	result := query.Count(&count)
	if result.Error != nil {
		return false, persistence.NewInternalError("failed to check entity existence by fields", result.Error)
	}

	return count > 0, nil
}

func (r *GenericRepositoryImpl[T, M, ID]) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &GenericRepositoryImpl[T, M, ID]{
			db:            tx,
			domainToModel: r.domainToModel,
			modelToDomain: r.modelToDomain,
		}

		txCtx := context.WithValue(ctx, "tx_repo", txRepo)

		return fn(txCtx)
	})
}

func (r *GenericRepositoryImpl[T, M, ID]) Refresh(ctx context.Context, entity T) error {
	id := entity.GetID()

	freshEntity, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}

	reflect.ValueOf(entity).Elem().Set(reflect.ValueOf(freshEntity).Elem())

	return nil
}

func (r *GenericRepositoryImpl[T, M, ID]) Clear(ctx context.Context) error {
	result := r.db.WithContext(ctx).Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(new(M))
	if result.Error != nil {
		return persistence.NewInternalError("failed to clear entities", result.Error)
	}

	return nil
}

func (r *GenericRepositoryImpl[T, M, ID]) mapFieldToColumn(field string) string {
	fieldMappings := map[string]string{
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"isActive":  "is_active",
	}

	if dbColumn, exists := fieldMappings[field]; exists {
		return dbColumn
	}

	return r.camelToSnake(field)
}

func (r *GenericRepositoryImpl[T, M, ID]) camelToSnake(str string) string {
	var result strings.Builder

	for i, char := range str {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(char)
	}

	return strings.ToLower(result.String())
}

func (r *GenericRepositoryImpl[T, M, ID]) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		if value != nil {
			dbColumn := r.mapFieldToColumn(key)
			query = query.Where(dbColumn+" = ?", value)
		}
	}

	return query
}

func (r *GenericRepositoryImpl[T, M, ID]) applySorting(query *gorm.DB, sortBy []sharedModels.SortBy) *gorm.DB {
	if sortBy == nil {
		return query
	}

	for _, sort := range sortBy {
		dbColumn := r.mapFieldToColumn(sort.Field)
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: dbColumn}, Desc: sort.Order == sharedModels.DESC})
	}

	return query
}

func (r *GenericRepositoryImpl[T, M, ID]) applyIncludes(query *gorm.DB, includes []string) *gorm.DB {
	if includes == nil {
		return query
	}

	for _, include := range includes {
		query = query.Preload(include)
	}

	return query
}
