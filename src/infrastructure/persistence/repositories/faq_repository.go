package repositories

import (
	"context"
	sharedModels "tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
)

type FAQRepositoryImpl struct {
	generic repositories.GenericRepository[*entities.FAQ, string]
}

func NewFAQRepository(generic repositories.GenericRepository[*entities.FAQ, string]) repositories.FAQRepository {
	return &FAQRepositoryImpl{generic}
}

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

func (r *FAQRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	return r.generic.SoftDelete(ctx, id)
}

func (r *FAQRepositoryImpl) FindAll(ctx context.Context, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	return r.generic.FindAll(ctx, opts)
}

func (r *FAQRepositoryImpl) FindOne(ctx context.Context, opts *sharedModels.QueryOptions) (*entities.FAQ, error) {
	return r.generic.FindOne(ctx, opts)
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

func (r *FAQRepositoryImpl) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	return r.generic.WithTransaction(ctx, fn)
}

func (r *FAQRepositoryImpl) Refresh(ctx context.Context, entity *entities.FAQ) error {
	return r.generic.Refresh(ctx, entity)
}

func (r *FAQRepositoryImpl) Clear(ctx context.Context) error {
	return r.generic.Clear(ctx)
}

func (r *FAQRepositoryImpl) GetCategories(ctx context.Context, withCounts bool) ([]string, map[string]int64, error) {
	opts := &sharedModels.QueryOptions{
		Filters: map[string]interface{}{
			"isActive": true,
		},
	}

	faqs, err := r.generic.FindAll(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	categoryMap := make(map[string]int64)

	for _, faq := range faqs {
		if faq.Category != "" {
			categoryMap[faq.Category]++
		}
	}

	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}

	if !withCounts {
		return categories, nil, nil
	}

	return categories, categoryMap, nil
}
