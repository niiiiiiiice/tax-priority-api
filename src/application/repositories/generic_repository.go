package repositories

import (
	"context"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/domain/entities"
)

type TransactionFunc func(ctx context.Context) error

type GenericRepository[T entities.Entity[ID], ID comparable] interface {
	// Базовые CRUD операции

	// Create - создание сущности
	Create(ctx context.Context, entity T) error
	// CreateBatch - создание пачки сущностей
	CreateBatch(ctx context.Context, entities []T) (*models.BulkOperationResult, error)

	// FindByID - поиск сущности по ID
	FindByID(ctx context.Context, id ID) (T, error)
	// FindByIDs - поиск сущностей по ID
	FindByIDs(ctx context.Context, ids []ID) ([]T, error)

	// Update - обновление сущности
	Update(ctx context.Context, entity T) error
	// UpdateBatch - обновление пачки сущностей
	UpdateBatch(ctx context.Context, entities []T) (*models.BulkOperationResult, error)
	// UpdateFields - обновление полей сущности
	UpdateFields(ctx context.Context, id ID, fields map[string]interface{}) error

	// Delete - удаление сущности
	Delete(ctx context.Context, id ID) error
	// DeleteBatch - удаление пачки сущностей
	DeleteBatch(ctx context.Context, ids []ID) (*models.BulkOperationResult, error)
	// SoftDelete - мягкое удаление сущности
	SoftDelete(ctx context.Context, id ID) error

	// Расширенные операции поиска

	// FindAll - поиск всех сущностей
	FindAll(ctx context.Context, opts *models.QueryOptions) ([]T, error)
	// FindOne - поиск одной сущности
	FindOne(ctx context.Context, opts *models.QueryOptions) (T, error)
	// FindWithPagination - поиск с пагинацией
	FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[T], error)

	// Операции подсчета

	// Count - подсчет количества сущностей
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	// Exists - проверка на существование сущности
	Exists(ctx context.Context, id ID) (bool, error)
	// ExistsByFields - проверка на существование сущности по полям
	ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error)

	// Транзакции

	// WithTransaction - выполнение транзакции
	WithTransaction(ctx context.Context, fn TransactionFunc) error

	// Утилиты

	// Refresh - обновление сущности
	Refresh(ctx context.Context, entity T) error
	// Clear - очистка кеша
	Clear(ctx context.Context) error
}
