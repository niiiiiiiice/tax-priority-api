package repositories

import (
	"context"
	"tax-priority-api/src/application/shared/models"
	"tax-priority-api/src/domain/entities"
)

// FAQRepository определяет интерфейс для работы с FAQ
type FAQRepository interface {
	// Базовые CRUD операции
	Create(ctx context.Context, faq *entities.FAQ) error
	CreateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error)

	FindByID(ctx context.Context, id string) (*entities.FAQ, error)
	FindByIDs(ctx context.Context, ids []string) ([]*entities.FAQ, error)

	Update(ctx context.Context, faq *entities.FAQ) error
	UpdateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error)
	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error

	Delete(ctx context.Context, id string) error
	DeleteBatch(ctx context.Context, ids []string) (*models.BulkOperationResult, error)

	// Расширенные операции поиска
	FindAll(ctx context.Context, opts *models.QueryOptions) ([]*entities.FAQ, error)
	FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.FAQ], error)

	// Специфичные для FAQ операции
	FindByCategory(ctx context.Context, category string, opts *models.QueryOptions) ([]*entities.FAQ, error)
	FindActive(ctx context.Context, opts *models.QueryOptions) ([]*entities.FAQ, error)
	FindByPriority(ctx context.Context, minPriority int, opts *models.QueryOptions) ([]*entities.FAQ, error)

	// Поиск
	Search(ctx context.Context, query string, opts *models.QueryOptions) ([]*entities.FAQ, error)
	SearchByCategory(ctx context.Context, query string, category string, opts *models.QueryOptions) ([]*entities.FAQ, error)

	// Операции подсчета
	Count(ctx context.Context, filters map[string]interface{}) (int64, error)
	CountByCategory(ctx context.Context, category string) (int64, error)
	CountActive(ctx context.Context) (int64, error)

	// Утилиты
	Exists(ctx context.Context, id string) (bool, error)
	ExistsByQuestion(ctx context.Context, question string) (bool, error)

	// Категории
	GetCategories(ctx context.Context) ([]string, error)
	GetCategoriesWithCounts(ctx context.Context) (map[string]int64, error)

	// Транзакции
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
