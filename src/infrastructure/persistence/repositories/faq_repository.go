package repositories

import (
	"context"

	"tax-priority-api/src/application/repositories"
	sharedModels "tax-priority-api/src/application/shared/models"
	"tax-priority-api/src/domain/entities"
)

// FAQRepositoryImpl реализация FAQRepository для GORM
type FAQRepositoryImpl struct {
	generic repositories.GenericRepository[*entities.FAQ, string]
}

// NewFAQRepository создает новый репозиторий FAQ
func NewFAQRepository(generic repositories.GenericRepository[*entities.FAQ, string]) repositories.FAQRepository {
	return &FAQRepositoryImpl{generic}
}

// Delegate all GenericRepository methods
func (r *FAQRepositoryImpl) Create(ctx context.Context, entity *entities.FAQ) error {
	return r.generic.Create(ctx, entity)
}
func (r *FAQRepositoryImpl) CreateBatch(ctx context.Context, entities []*entities.FAQ) (*sharedModels.BulkOperationResult, error) {
	return r.generic.CreateBatch(ctx, entities)
}
func (r *FAQRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.FAQ, error) {
	return r.generic.FindByID(ctx, id)
}
func (r *FAQRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.FAQ, error) {
	return r.generic.FindByIDs(ctx, ids)
}
func (r *FAQRepositoryImpl) Update(ctx context.Context, entity *entities.FAQ) error {
	return r.generic.Update(ctx, entity)
}
func (r *FAQRepositoryImpl) UpdateBatch(ctx context.Context, entities []*entities.FAQ) (*sharedModels.BulkOperationResult, error) {
	return r.generic.UpdateBatch(ctx, entities)
}
func (r *FAQRepositoryImpl) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.generic.UpdateFields(ctx, id, fields)
}
func (r *FAQRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.generic.Delete(ctx, id)
}
func (r *FAQRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*sharedModels.BulkOperationResult, error) {
	return r.generic.DeleteBatch(ctx, ids)
}
func (r *FAQRepositoryImpl) FindAll(ctx context.Context, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	return r.generic.FindAll(ctx, opts)
}
func (r *FAQRepositoryImpl) FindWithPagination(ctx context.Context, opts *sharedModels.QueryOptions) (*sharedModels.PaginatedResult[*entities.FAQ], error) {
	return r.generic.FindWithPagination(ctx, opts)
}
func (r *FAQRepositoryImpl) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	return r.generic.Count(ctx, filters)
}
func (r *FAQRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	return r.generic.Exists(ctx, id)
}
func (r *FAQRepositoryImpl) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	return r.generic.ExistsByFields(ctx, filters)
}
func (r *FAQRepositoryImpl) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return r.generic.WithTransaction(ctx, repositories.TransactionFunc(fn))
}

// FindByCategory находит FAQ по категории
func (r *FAQRepositoryImpl) FindByCategory(ctx context.Context, category string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	localOpts := &sharedModels.QueryOptions{}
	if opts != nil {
		*localOpts = *opts
	}
	if localOpts.Filters == nil {
		localOpts.Filters = make(map[string]interface{})
	}
	localOpts.Filters["category"] = category
	return r.FindAll(ctx, localOpts)
}

// FindActive находит активные FAQ
func (r *FAQRepositoryImpl) FindActive(ctx context.Context, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	localOpts := &sharedModels.QueryOptions{}
	if opts != nil {
		*localOpts = *opts
	}
	if localOpts.Filters == nil {
		localOpts.Filters = make(map[string]interface{})
	}
	localOpts.Filters["is_active"] = true

	return r.FindAll(ctx, localOpts)
}

// FindByPriority находит FAQ по приоритету
func (r *FAQRepositoryImpl) FindByPriority(ctx context.Context, minPriority int, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	localOpts := &sharedModels.QueryOptions{}
	if opts != nil {
		*localOpts = *opts
	}
	if localOpts.Filters == nil {
		localOpts.Filters = make(map[string]interface{})
	}
	localOpts.Filters["priority >="] = minPriority
	return r.FindAll(ctx, localOpts)
}

// Search выполняет поиск FAQ
func (r *FAQRepositoryImpl) Search(ctx context.Context, searchQuery string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	// TODO: Implement search logic using GenericRepository or add search support to GenericRepository
	return nil, nil
}

// SearchByCategory выполняет поиск FAQ по категории
func (r *FAQRepositoryImpl) SearchByCategory(ctx context.Context, searchQuery string, category string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	// TODO: Implement search by category logic using GenericRepository or add search support to GenericRepository
	return nil, nil
}

// CountByCategory подсчитывает количество FAQ по категории
func (r *FAQRepositoryImpl) CountByCategory(ctx context.Context, category string) (int64, error) {
	return r.Count(ctx, map[string]interface{}{"category": category})
}

// CountActive подсчитывает количество активных FAQ
func (r *FAQRepositoryImpl) CountActive(ctx context.Context) (int64, error) {
	return r.Count(ctx, map[string]interface{}{"is_active": true})
}

// ExistsByQuestion проверяет существование FAQ по вопросу
func (r *FAQRepositoryImpl) ExistsByQuestion(ctx context.Context, question string) (bool, error) {
	return r.ExistsByFields(ctx, map[string]interface{}{"question": question})
}

// GetCategories получает все категории FAQ
func (r *FAQRepositoryImpl) GetCategories(ctx context.Context) ([]string, error) {
	// TODO: Implement get categories logic using GenericRepository or add distinct support to GenericRepository
	return nil, nil
}

// GetCategoriesWithCounts получает категории FAQ с количеством
func (r *FAQRepositoryImpl) GetCategoriesWithCounts(ctx context.Context) (map[string]int64, error) {
	// TODO: Implement get categories with counts logic using GenericRepository or add grouping support to GenericRepository
	return nil, nil
}
