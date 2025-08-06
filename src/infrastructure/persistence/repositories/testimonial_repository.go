package repositories

import (
	"context"
	"tax-priority-api/src/application/models"
	sharedModels "tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"time"
)

type TestimonialRepositoryImpl struct {
	generic repositories.GenericRepository[*entities.Testimonial, string]
}

func NewTestimonialRepository(generic repositories.GenericRepository[*entities.Testimonial, string]) repositories.TestimonialRepository {
	return &TestimonialRepositoryImpl{generic}
}

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

func (r *TestimonialRepositoryImpl) FindByApprovalStatus(ctx context.Context, isApproved bool, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["isApproved"] = isApproved

	return r.generic.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) FindByRating(ctx context.Context, rating int, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["rating"] = rating

	return r.generic.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) FindByAuthor(ctx context.Context, author string, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["author"] = author

	return r.generic.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) FindByAuthorEmail(ctx context.Context, authorEmail string, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["authorEmail"] = authorEmail

	return r.generic.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) FindApprovedAndActive(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["isApproved"] = true
	opts.Filters["isActive"] = true

	return r.generic.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) FindWithFiles(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	// Ищем отзывы, у которых есть файлы (filePath не пустой)
	opts.Filters["filePath"] = map[string]interface{}{"$ne": ""}

	return r.generic.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) GetAverageRating(ctx context.Context) (float64, error) {
	// Для получения среднего рейтинга нужно использовать SQL агрегацию
	// Пока возвращаем простую реализацию через получение всех записей
	testimonials, err := r.generic.FindAll(ctx, &models.QueryOptions{
		Filters: map[string]interface{}{
			"isApproved": true,
			"isActive":   true,
		},
	})
	if err != nil {
		return 0, err
	}

	if len(testimonials) == 0 {
		return 0, nil
	}

	var total float64
	for _, testimonial := range testimonials {
		total += float64(testimonial.Rating)
	}

	return total / float64(len(testimonials)), nil
}

func (r *TestimonialRepositoryImpl) GetRatingDistribution(ctx context.Context) (map[int]int64, error) {
	// Получаем все одобренные и активные отзывы
	testimonials, err := r.generic.FindAll(ctx, &models.QueryOptions{
		Filters: map[string]interface{}{
			"isApproved": true,
			"isActive":   true,
		},
	})
	if err != nil {
		return nil, err
	}

	distribution := make(map[int]int64)
	for i := 1; i <= 5; i++ {
		distribution[i] = 0
	}

	for _, testimonial := range testimonials {
		distribution[testimonial.Rating]++
	}

	return distribution, nil
}

func (r *TestimonialRepositoryImpl) CountByApprovalStatus(ctx context.Context, isApproved bool) (int64, error) {
	return r.generic.Count(ctx, map[string]interface{}{
		"isApproved": isApproved,
	})
}

func (r *TestimonialRepositoryImpl) ApproveMany(ctx context.Context, ids []string, approvedBy string) (*models.BulkOperationResult, error) {
	if len(ids) == 0 {
		return &models.BulkOperationResult{SuccessCount: 0, FailureCount: 0}, nil
	}

	now := time.Now()
	fields := map[string]interface{}{
		"isApproved": true,
		"approvedAt": now,
		"approvedBy": approvedBy,
		"updatedAt":  now,
	}

	var successCount, failureCount int
	var errors []error

	for _, id := range ids {
		err := r.generic.UpdateFields(ctx, id, fields)
		if err != nil {
			failureCount++
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	return &models.BulkOperationResult{
		SuccessCount: successCount,
		FailureCount: failureCount,
		Errors:       errors,
	}, nil
}

func (r *TestimonialRepositoryImpl) DeactivateMany(ctx context.Context, ids []string) (*models.BulkOperationResult, error) {
	if len(ids) == 0 {
		return &models.BulkOperationResult{SuccessCount: 0, FailureCount: 0}, nil
	}

	fields := map[string]interface{}{
		"isActive":  false,
		"updatedAt": time.Now(),
	}

	var successCount, failureCount int
	var errors []error

	for _, id := range ids {
		err := r.generic.UpdateFields(ctx, id, fields)
		if err != nil {
			failureCount++
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	return &models.BulkOperationResult{
		SuccessCount: successCount,
		FailureCount: failureCount,
		Errors:       errors,
	}, nil
}

func (r *TestimonialRepositoryImpl) ActivateMany(ctx context.Context, ids []string) (*models.BulkOperationResult, error) {
	if len(ids) == 0 {
		return &models.BulkOperationResult{SuccessCount: 0, FailureCount: 0}, nil
	}

	fields := map[string]interface{}{
		"isActive":  true,
		"updatedAt": time.Now(),
	}

	var successCount, failureCount int
	var errors []error

	for _, id := range ids {
		err := r.generic.UpdateFields(ctx, id, fields)
		if err != nil {
			failureCount++
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	return &models.BulkOperationResult{
		SuccessCount: successCount,
		FailureCount: failureCount,
		Errors:       errors,
	}, nil
}

func (r *TestimonialRepositoryImpl) DeleteMany(ctx context.Context, ids []string) (*models.BulkOperationResult, error) {
	return r.generic.DeleteBatch(ctx, ids)
}
