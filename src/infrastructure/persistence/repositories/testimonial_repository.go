package repositories

import (
	"context"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
	"time"
)

type TestimonialRepositoryImpl struct {
	repositories.GenericRepository[*entities.Testimonial, string]
}

func NewTestimonialRepository(generic repositories.GenericRepository[*entities.Testimonial, string]) repositories.TestimonialRepository {
	return &TestimonialRepositoryImpl{generic}
}

func (r *TestimonialRepositoryImpl) FindByApprovalStatus(ctx context.Context, isApproved bool, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["isApproved"] = isApproved

	return r.GenericRepository.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) FindByRating(ctx context.Context, rating int, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["rating"] = rating

	return r.GenericRepository.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) FindByAuthor(ctx context.Context, author string, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["author"] = author

	return r.GenericRepository.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) FindByAuthorEmail(ctx context.Context, authorEmail string, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error) {
	if opts == nil {
		opts = &models.QueryOptions{}
	}
	if opts.Filters == nil {
		opts.Filters = make(map[string]interface{})
	}
	opts.Filters["authorEmail"] = authorEmail

	return r.GenericRepository.FindWithPagination(ctx, opts)
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

	return r.GenericRepository.FindWithPagination(ctx, opts)
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

	return r.GenericRepository.FindWithPagination(ctx, opts)
}

func (r *TestimonialRepositoryImpl) GetAverageRating(ctx context.Context) (float64, error) {
	// Для получения среднего рейтинга нужно использовать SQL агрегацию
	// Пока возвращаем простую реализацию через получение всех записей
	testimonials, err := r.GenericRepository.FindAll(ctx, &models.QueryOptions{
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
	testimonials, err := r.GenericRepository.FindAll(ctx, &models.QueryOptions{
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
	return r.GenericRepository.Count(ctx, map[string]interface{}{
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
		err := r.GenericRepository.UpdateFields(ctx, id, fields)
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
		err := r.GenericRepository.UpdateFields(ctx, id, fields)
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
		err := r.GenericRepository.UpdateFields(ctx, id, fields)
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
	return r.GenericRepository.DeleteBatch(ctx, ids)
}
