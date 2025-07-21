package repositories

import (
	"context"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/domain/entities"
)

// FAQRepository определяет интерфейс для работы с FAQ
type FAQRepository interface {
	GenericRepository[*entities.FAQ, string]

	UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error

	Delete(ctx context.Context, id string) error
	DeleteBatch(ctx context.Context, ids []string) (*models.BulkOperationResult, error)

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
}
