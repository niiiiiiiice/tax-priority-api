package repositories

import (
	"context"
	sharedModels "tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
)

type TestimonialRepositoryImpl struct {
	generic repositories.GenericRepository[*entities.Testimonial, string]
}

func NewTestimonialRepository(generic repositories.GenericRepository[*entities.Testimonial, string]) repositories.TestimonialRepository {
	return &TestimonialRepositoryImpl{generic}
}

// Delegate all GenericRepository methods
func (r *TestimonialRepositoryImpl) Create(ctx context.Context, entity *entities.Testimonial) error {
	return r.generic.Create(ctx, entity)
}
func (r *TestimonialRepositoryImpl) CreateBatch(ctx context.Context, entities []*entities.Testimonial) (*sharedModels.BulkOperationResult, error) {
	return r.generic.CreateBatch(ctx, entities)
}

func (r *TestimonialRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Testimonial, error) {
	return r.generic.FindByID(ctx, id)
}
func (r *TestimonialRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.Testimonial, error) {
	return r.generic.FindByIDs(ctx, ids)
}

func (r *TestimonialRepositoryImpl) Update(ctx context.Context, entity *entities.Testimonial) error {
	return r.generic.Update(ctx, entity)
}
func (r *TestimonialRepositoryImpl) UpdateBatch(ctx context.Context, entities []*entities.Testimonial) (*sharedModels.BulkOperationResult, error) {
	return r.generic.UpdateBatch(ctx, entities)
}
func (r *TestimonialRepositoryImpl) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.generic.UpdateFields(ctx, id, fields)
}

func (r *TestimonialRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.generic.Delete(ctx, id)
}
func (r *TestimonialRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*sharedModels.BulkOperationResult, error) {
	return r.generic.DeleteBatch(ctx, ids)
}
func (r *TestimonialRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	return r.generic.SoftDelete(ctx, id)
}

func (r *TestimonialRepositoryImpl) FindAll(ctx context.Context, opts *sharedModels.QueryOptions) ([]*entities.Testimonial, error) {
	return r.generic.FindAll(ctx, opts)
}
func (r *TestimonialRepositoryImpl) FindOne(ctx context.Context, opts *sharedModels.QueryOptions) (*entities.Testimonial, error) {
	return r.generic.FindOne(ctx, opts)
}
func (r *TestimonialRepositoryImpl) FindWithPagination(ctx context.Context, opts *sharedModels.QueryOptions) (*sharedModels.PaginatedResult[*entities.Testimonial], error) {
	return r.generic.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	return r.generic.Count(ctx, filters)
}
func (r *TestimonialRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	return r.generic.Exists(ctx, id)
}
func (r *TestimonialRepositoryImpl) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	return r.generic.ExistsByFields(ctx, filters)
}

func (r *TestimonialRepositoryImpl) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	return r.generic.WithTransaction(ctx, repositories.TransactionFunc(fn))
}

func (r *TestimonialRepositoryImpl) Refresh(ctx context.Context, entity *entities.Testimonial) error {
	return r.generic.Refresh(ctx, entity)
}
func (r *TestimonialRepositoryImpl) Clear(ctx context.Context) error {
	return r.generic.Clear(ctx)
}
