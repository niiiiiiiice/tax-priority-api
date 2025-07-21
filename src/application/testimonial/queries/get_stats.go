package queries

import (
	"context"
	"fmt"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/application/testimonial/dtos"
	"time"
)

type GetTestimonialStatsQueryHandler struct {
	testimonialRepo repositories.TestimonialRepository
}

func NewGetTestimonialStatsQueryHandler(repo repositories.TestimonialRepository) *GetTestimonialStatsQueryHandler {
	return &GetTestimonialStatsQueryHandler{
		testimonialRepo: repo,
	}
}

func (h *GetTestimonialStatsQueryHandler) Handle(ctx context.Context, query dtos.GetTestimonialStatsQuery) (*dtos.QueryResult, error) {
	stats := &dtos.TestimonialStats{}

	// Получаем общее количество активных отзывов
	totalFilters := map[string]interface{}{
		"is_active": true,
	}
	totalCount, err := h.testimonialRepo.Count(ctx, totalFilters)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get total count: %v", err),
			Timestamp: time.Now(),
		}, err
	}
	stats.TotalCount = totalCount

	// Получаем количество одобренных отзывов
	approvedCount, err := h.testimonialRepo.CountByApprovalStatus(ctx, true)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get approved count: %v", err),
			Timestamp: time.Now(),
		}, err
	}
	stats.ApprovedCount = approvedCount

	// Получаем количество ожидающих одобрения отзывов
	pendingCount, err := h.testimonialRepo.CountByApprovalStatus(ctx, false)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get pending count: %v", err),
			Timestamp: time.Now(),
		}, err
	}
	stats.PendingCount = pendingCount

	// Получаем средний рейтинг
	avgRating, err := h.testimonialRepo.GetAverageRating(ctx)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get average rating: %v", err),
			Timestamp: time.Now(),
		}, err
	}
	stats.AverageRating = avgRating

	// Получаем распределение рейтингов
	ratingDistribution, err := h.testimonialRepo.GetRatingDistribution(ctx)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get rating distribution: %v", err),
			Timestamp: time.Now(),
		}, err
	}
	stats.RatingDistribution = ratingDistribution

	// Получаем количество отзывов с файлами
	withFilesFilters := map[string]interface{}{
		"is_active":             true,
		"file_path IS NOT NULL": "",
		"file_path != ''":       "",
	}
	withFilesCount, err := h.testimonialRepo.Count(ctx, withFilesFilters)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get files count: %v", err),
			Timestamp: time.Now(),
		}, err
	}
	stats.WithFilesCount = withFilesCount

	// Получаем количество недавних отзывов (за последние 30 дней)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	recentFilters := map[string]interface{}{
		"is_active":     true,
		"created_at >=": thirtyDaysAgo,
	}
	recentCount, err := h.testimonialRepo.Count(ctx, recentFilters)
	if err != nil {
		return &dtos.QueryResult{
			Success:   false,
			Error:     fmt.Sprintf("failed to get recent count: %v", err),
			Timestamp: time.Now(),
		}, err
	}
	stats.RecentCount = recentCount

	return &dtos.QueryResult{
		Success:   true,
		Message:   "Testimonial statistics retrieved successfully",
		Stats:     stats,
		Timestamp: time.Now(),
	}, nil
}
