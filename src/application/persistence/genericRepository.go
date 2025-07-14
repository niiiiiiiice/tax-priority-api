package persistence

import (
	"context"
	"tax-priority-api/src/application/shared/models"
	"tax-priority-api/src/domain/entities"
)

type TransactionFunc func(ctx context.Context) error

type GenericRepository[T entities.Entity[ID], ID comparable] interface {
	// Базовые CRUD операции
	Create(ctx context.Context, entity T) error
	CreateBatch(ctx context.Context, entities []T) (*models.BulkOperationResult, error)

	FindByID(ctx context.Context, id ID) (T, error)
	FindByIDs(ctx context.Context, ids []ID) ([]T, error)

	Update(ctx context.Context, entity T) error
	UpdateBatch(ctx context.Context, entities []T) (*models.BulkOperationResult, error)
	UpdateFields(ctx context.Context, id ID, fields map[string]interface{}) error

	Delete(ctx context.Context, id ID) error
	DeleteBatch(ctx context.Context, ids []ID) (*models.BulkOperationResult, error)
	SoftDelete(ctx context.Context, id ID) error

	// Расширенные операции поиска
	FindAll(ctx context.Context, opts *models.QueryOptions) ([]T, error)
	FindOne(ctx context.Context, opts *models.QueryOptions) (T, error)
	FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[T], error)

	// Операции подсчета
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	Exists(ctx context.Context, id ID) (bool, error)
	ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error)

	// Транзакции
	WithTransaction(ctx context.Context, fn TransactionFunc) error

	// Утилиты
	Refresh(ctx context.Context, entity T) error
	Clear(ctx context.Context) error
}
