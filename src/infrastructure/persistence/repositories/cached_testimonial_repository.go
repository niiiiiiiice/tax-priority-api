package repositories

import (
	"context"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
)

// CachedTestimonialRepositoryImpl кешированный репозиторий отзывов
type CachedTestimonialRepositoryImpl struct {
	repo repositories.TestimonialRepository
}

// NewCachedTestimonialRepositoryImpl создает новый кешированный репозиторий FAQ
func NewCachedTestimonialRepositoryImpl(
	repo repositories.TestimonialRepository,
) repositories.CachedTestimonialRepository {
	return &CachedTestimonialRepositoryImpl{
		repo: repo,
	}
}

func (r *CachedTestimonialRepositoryImpl) Create(ctx context.Context, testimonial *entities.Testimonial) error {
	return r.repo.Create(ctx, testimonial)
}

func (r *CachedTestimonialRepositoryImpl) CreateBatch(ctx context.Context, testimonials []*entities.Testimonial) (*models.BulkOperationResult, error) {
	return r.repo.CreateBatch(ctx, testimonials)
}

func (r *CachedTestimonialRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Testimonial, error) {
	return r.repo.FindByID(ctx, id)
}

func (r *CachedTestimonialRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.Testimonial, error) {
	return r.repo.FindByIDs(ctx, ids)
}

func (r *CachedTestimonialRepositoryImpl) Update(ctx context.Context, testimonial *entities.Testimonial) error {
	return r.repo.Update(ctx, testimonial)
}

func (r *CachedTestimonialRepositoryImpl) UpdateBatch(ctx context.Context, testimonials []*entities.Testimonial) (*models.BulkOperationResult, error) {
	return r.repo.UpdateBatch(ctx, testimonials)
}

func (r *CachedTestimonialRepositoryImpl) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.repo.UpdateFields(ctx, id, fields)
}

func (r *CachedTestimonialRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.repo.Delete(ctx, id)
}

func (r *CachedTestimonialRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*models.BulkOperationResult, error) {
	return r.repo.DeleteBatch(ctx, ids)
}

func (r *CachedTestimonialRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	return r.repo.SoftDelete(ctx, id)
}

func (r *CachedTestimonialRepositoryImpl) FindAll(ctx context.Context, opts *models.QueryOptions) ([]*entities.Testimonial, error) {
	return r.repo.FindAll(ctx, opts)
}

func (r *CachedTestimonialRepositoryImpl) FindOne(ctx context.Context, opts *models.QueryOptions) (*entities.Testimonial, error) {
	return r.repo.FindOne(ctx, opts)
}

func (r *CachedTestimonialRepositoryImpl) FindWithPagination(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	return r.repo.FindWithPagination(ctx, opts)
}

func (r *CachedTestimonialRepositoryImpl) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	return r.repo.Count(ctx, filters)
}

func (r *CachedTestimonialRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	return r.repo.Exists(ctx, id)
}

func (r *CachedTestimonialRepositoryImpl) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	return r.repo.ExistsByFields(ctx, filters)
}

func (r *CachedTestimonialRepositoryImpl) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	return r.repo.WithTransaction(ctx, fn)
}

func (r *CachedTestimonialRepositoryImpl) Refresh(ctx context.Context, entity *entities.Testimonial) error {
	return r.repo.Refresh(ctx, entity)
}

func (r *CachedTestimonialRepositoryImpl) Clear(ctx context.Context) error {
	return r.repo.Clear(ctx)
}
