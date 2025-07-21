package dtos

import (
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/domain/entities"
	"time"
)

// GetTestimonialsQuery для получения списка отзывов
type GetTestimonialsQuery struct {
	Limit     int                    `json:"limit" validate:"min=1,max=100"`
	Offset    int                    `json:"offset" validate:"min=0"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder" validate:"oneof=asc desc"`
	Filters   map[string]interface{} `json:"filters"`
}

// GetTestimonialByIDQuery для получения отзыва по ID
type GetTestimonialByIDQuery struct {
	ID string `json:"id" validate:"required"`
}

// GetTestimonialsByRatingQuery для получения отзывов по рейтингу
type GetTestimonialsByRatingQuery struct {
	Rating    int                    `json:"rating" validate:"required,min=1,max=5"`
	Limit     int                    `json:"limit" validate:"min=1,max=100"`
	Offset    int                    `json:"offset" validate:"min=0"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder" validate:"oneof=asc desc"`
	Filters   map[string]interface{} `json:"filters"`
}

// GetTestimonialsByAuthorQuery для получения отзывов по автору
type GetTestimonialsByAuthorQuery struct {
	Author    string                 `json:"author" validate:"required"`
	Limit     int                    `json:"limit" validate:"min=1,max=100"`
	Offset    int                    `json:"offset" validate:"min=0"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder" validate:"oneof=asc desc"`
	Filters   map[string]interface{} `json:"filters"`
}

// GetApprovedTestimonialsQuery для получения одобренных отзывов
type GetApprovedTestimonialsQuery struct {
	Limit     int                    `json:"limit" validate:"min=1,max=100"`
	Offset    int                    `json:"offset" validate:"min=0"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder" validate:"oneof=asc desc"`
	Filters   map[string]interface{} `json:"filters"`
}

// GetPendingTestimonialsQuery для получения ожидающих одобрения отзывов
type GetPendingTestimonialsQuery struct {
	Limit     int                    `json:"limit" validate:"min=1,max=100"`
	Offset    int                    `json:"offset" validate:"min=0"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder" validate:"oneof=asc desc"`
	Filters   map[string]interface{} `json:"filters"`
}

// GetTestimonialsWithFilesQuery для получения отзывов с файлами
type GetTestimonialsWithFilesQuery struct {
	Limit     int                    `json:"limit" validate:"min=1,max=100"`
	Offset    int                    `json:"offset" validate:"min=0"`
	SortBy    string                 `json:"sortBy"`
	SortOrder string                 `json:"sortOrder" validate:"oneof=asc desc"`
	Filters   map[string]interface{} `json:"filters"`
}

// GetTestimonialStatsQuery для получения статистики отзывов
type GetTestimonialStatsQuery struct {
	// Можно добавить фильтры для статистики
}

// QueryResult общий результат выполнения запроса
type QueryResult struct {
	Success   bool                                           `json:"success"`
	Message   string                                         `json:"message,omitempty"`
	Error     string                                         `json:"error,omitempty"`
	Data      *entities.Testimonial                          `json:"data,omitempty"`
	Paginated *models.PaginatedResult[*entities.Testimonial] `json:"paginated,omitempty"`
	Stats     *TestimonialStats                              `json:"stats,omitempty"`
	Timestamp time.Time                                      `json:"timestamp"`
}

// TestimonialStats статистика отзывов
type TestimonialStats struct {
	TotalCount         int64         `json:"totalCount"`
	ApprovedCount      int64         `json:"approvedCount"`
	PendingCount       int64         `json:"pendingCount"`
	AverageRating      float64       `json:"averageRating"`
	RatingDistribution map[int]int64 `json:"ratingDistribution"`
	WithFilesCount     int64         `json:"withFilesCount"`
	RecentCount        int64         `json:"recentCount"` // За последние 30 дней
}
