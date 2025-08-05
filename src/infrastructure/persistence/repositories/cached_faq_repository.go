package repositories

import (
	"context"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
)

// CachedFAQRepositoryImpl кешированный репозиторий FAQ
type CachedFAQRepositoryImpl struct {
	repo repositories.FAQRepository
}

func NewCachedFAQRepository(
	repo repositories.FAQRepository,
) repositories.CachedFAQRepository {
	return &CachedFAQRepositoryImpl{
		repo: repo,
	}
}

func (cr *CachedFAQRepositoryImpl) Create(ctx context.Context, faq *entities.FAQ) error {
	return cr.repo.Create(ctx, faq)
}

func (cr *CachedFAQRepositoryImpl) CreateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error) {
	return cr.repo.CreateBatch(ctx, faqs)
}

func (cr *CachedFAQRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.FAQ, error) {
	return cr.repo.FindByID(ctx, id)
}

func (cr *CachedFAQRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.FAQ, error) {
	return cr.repo.FindByIDs(ctx, ids)
}

func (cr *CachedFAQRepositoryImpl) Update(ctx context.Context, faq *entities.FAQ) error {
	return cr.repo.Update(ctx, faq)
}

func (cr *CachedFAQRepositoryImpl) UpdateBatch(ctx context.Context, faqs []*entities.FAQ) (*models.BulkOperationResult, error) {
	return cr.repo.UpdateBatch(ctx, faqs)
}

func (cr *CachedFAQRepositoryImpl) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return cr.repo.UpdateFields(ctx, id, fields)
}

func (cr *CachedFAQRepositoryImpl) Delete(ctx context.Context, id string) error {
	return cr.repo.Delete(ctx, id)
}

func (cr *CachedFAQRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*models.BulkOperationResult, error) {
	return cr.repo.DeleteBatch(ctx, ids)
}

func (cr *CachedFAQRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	return cr.repo.SoftDelete(ctx, id)
}

func (cr *CachedFAQRepositoryImpl) FindAll(ctx context.Context, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	return cr.repo.FindAll(ctx, opts)
}

func (cr *CachedFAQRepositoryImpl) FindOne(ctx context.Context, opts *models.QueryOptions) (*entities.FAQ, error) {
	return cr.repo.FindOne(ctx, opts)
}

func (cr *CachedFAQRepositoryImpl) FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.FAQ], error) {
	return cr.repo.FindWithPagination(ctx, opts)
}

func (cr *CachedFAQRepositoryImpl) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	return cr.repo.Count(ctx, filters)
}

func (cr *CachedFAQRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	return cr.repo.Exists(ctx, id)
}

func (cr *CachedFAQRepositoryImpl) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	return cr.repo.ExistsByFields(ctx, filters)
}

func (cr *CachedFAQRepositoryImpl) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	return cr.repo.WithTransaction(ctx, fn)
}

func (cr *CachedFAQRepositoryImpl) Refresh(ctx context.Context, faq *entities.FAQ) error {
	return cr.repo.Refresh(ctx, faq)
}

func (cr *CachedFAQRepositoryImpl) Clear(ctx context.Context) error {
	return cr.repo.Clear(ctx)
}
