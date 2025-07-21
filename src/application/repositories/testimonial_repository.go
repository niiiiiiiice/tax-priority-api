package repositories

import (
	"context"
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/domain/entities"
)

type TestimonialRepository interface {
	GenericRepository[*entities.Testimonial, string]

	// Специфические методы для отзывов
	FindByApprovalStatus(ctx context.Context, isApproved bool, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error)
	FindByRating(ctx context.Context, rating int, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error)
	FindByAuthor(ctx context.Context, author string, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error)
	FindByAuthorEmail(ctx context.Context, authorEmail string, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error)
	FindApprovedAndActive(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error)
	FindWithFiles(ctx context.Context, opts *models.QueryOptions) (*models.PaginatedResult[*entities.Testimonial], error)
	GetAverageRating(ctx context.Context) (float64, error)
	GetRatingDistribution(ctx context.Context) (map[int]int64, error)
	CountByApprovalStatus(ctx context.Context, isApproved bool) (int64, error)

	// Методы для массовых операций
	ApproveMany(ctx context.Context, ids []string, approvedBy string) (*models.BulkOperationResult, error)
	DeactivateMany(ctx context.Context, ids []string) (*models.BulkOperationResult, error)
	ActivateMany(ctx context.Context, ids []string) (*models.BulkOperationResult, error)
	DeleteMany(ctx context.Context, ids []string) (*models.BulkOperationResult, error)
}
