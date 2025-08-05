package repositories

import (
	"context"

	"tax-priority-api/src/application/cache"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/infrastructure/persistence"
)

// CachedFAQRepositoryImpl кешированный репозиторий FAQ
type CachedFAQRepositoryImpl struct {
	repo repositories.FAQRepository
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
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*models.BulkOperationResult, error) {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) FindAll(ctx context.Context, opts *models.QueryOptions) ([]*entities.FAQ, error) {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) FindOne(ctx context.Context, opts *models.QueryOptions) (*entities.FAQ, error) {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.FAQ], error) {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) Refresh(ctx context.Context, faq *entities.FAQ) error {
	//TODO implement me
	panic("implement me")
}

func (cr *CachedFAQRepositoryImpl) Clear(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewCachedFAQRepository(
	repo repositories.FAQRepository,
	cache cache.Cache,
	keys *persistence.RedisKeys,
	config *CacheConfig,
) repositories.CachedFAQRepository {

	r := &CachedFAQRepositoryImpl{
		repo: repo,
	}

	//if config.WarmupOnStart {
	//	go r.warmupCache(context.Background())
	//}

	return r
}
