package models

import "time"

// CreateFAQRequest модель для создания FAQ
type CreateFAQRequest struct {
	Question string `json:"question" validate:"required,min=10,max=500" example:"Как подать налоговую декларацию?"`
	Answer   string `json:"answer" validate:"required,min=10,max=2000" example:"Для подачи налоговой декларации необходимо..."`
	Category string `json:"category" validate:"required,max=100" example:"налоги"`
	Priority int    `json:"priority" validate:"min=0,max=100" example:"50"`
}

// UpdateFAQRequest модель для обновления FAQ
type UpdateFAQRequest struct {
	Question string `json:"question" validate:"required,min=10,max=500" example:"Как подать налоговую декларацию?"`
	Answer   string `json:"answer" validate:"required,min=10,max=2000" example:"Для подачи налоговой декларации необходимо..."`
	Category string `json:"category" validate:"required,max=100" example:"налоги"`
	Priority int    `json:"priority" validate:"min=0,max=100" example:"50"`
}

// UpdateFAQPriorityRequest модель для обновления приоритета FAQ
type UpdateFAQPriorityRequest struct {
	Priority int `json:"priority" validate:"min=0,max=100" example:"75"`
}

// BulkDeleteFAQRequest модель для массового удаления FAQ
type BulkDeleteFAQRequest struct {
	IDs []string `json:"ids" validate:"required,min=1" example:"[\"uuid1\", \"uuid2\"]"`
}

// GetFAQsByIDsRequest модель для получения FAQ по списку ID
type GetFAQsByIDsRequest struct {
	IDs []string `json:"ids" validate:"required,min=1" example:"[\"uuid1\", \"uuid2\"]"`
}

// FAQResponse модель ответа FAQ
type FAQResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Question  string    `json:"question" example:"Как подать налоговую декларацию?"`
	Answer    string    `json:"answer" example:"Для подачи налоговой декларации необходимо..."`
	Category  string    `json:"category" example:"налоги"`
	IsActive  bool      `json:"isActive" example:"true"`
	Priority  int       `json:"priority" example:"50"`
	CreatedAt time.Time `json:"createdAt" example:"2023-12-01T10:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt" example:"2023-12-01T10:00:00Z"`
}

// PaginatedFAQResponse модель пагинированного ответа FAQ
type PaginatedFAQResponse struct {
	Items      []FAQResponse `json:"items"`
	Total      int64         `json:"total" example:"100"`
	Offset     int           `json:"offset" example:"0"`
	Limit      int           `json:"limit" example:"10"`
	HasNext    bool          `json:"hasNext" example:"true"`
	HasPrev    bool          `json:"hasPrev" example:"false"`
	TotalPages int           `json:"totalPages" example:"10"`
}

// CategoryResponse модель ответа категории
type CategoryResponse struct {
	Name  string `json:"name" example:"налоги"`
	Count int64  `json:"count,omitempty" example:"25"`
}

// CommandResult модель результата выполнения команды
type CommandResult struct {
	ID        string    `json:"id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	Success   bool      `json:"success" example:"true"`
	Message   string    `json:"message,omitempty" example:"FAQ created successfully"`
	Error     string    `json:"error,omitempty" example:"Validation failed"`
	CreatedAt time.Time `json:"createdAt,omitempty" example:"2023-12-01T10:00:00Z"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" example:"2023-12-01T10:00:00Z"`
}

// BatchCommandResult модель результата выполнения batch команды
type BatchCommandResult struct {
	SuccessCount int             `json:"successCount" example:"8"`
	FailureCount int             `json:"failureCount" example:"2"`
	Results      []CommandResult `json:"results"`
	Errors       []string        `json:"errors,omitempty" example:"[\"Validation failed for item 1\"]"`
}

// ErrorResponse модель ошибки
type ErrorResponse struct {
	Error string `json:"error" example:"Validation failed"`
}

// HealthResponse модель ответа health check
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Tax Priority API is running"`
}

// CountResponse модель ответа с количеством
type CountResponse struct {
	Count int64 `json:"count" example:"42"`
}

// SearchFAQsQuery модель для поиска FAQ
type SearchFAQsQuery struct {
	Query      string `form:"q" binding:"required,min=3" example:"налоги"`
	Category   string `form:"category" example:"налоги"`
	Limit      int    `form:"_limit" example:"10"`
	Offset     int    `form:"_offset" example:"0"`
	SortBy     string `form:"_sort" example:"priority"`
	SortOrder  string `form:"_order" example:"desc"`
	ActiveOnly bool   `form:"activeOnly" example:"true"`
}

// GetFAQsQuery модель для получения списка FAQ
type GetFAQsQuery struct {
	Limit     int    `form:"_limit" example:"10"`
	Offset    int    `form:"_offset" example:"0"`
	SortBy    string `form:"_sort" example:"createdAt"`
	SortOrder string `form:"_order" example:"desc"`
	Category  string `form:"category" example:"налоги"`
	IsActive  bool   `form:"isActive" example:"true"`
}

// GetFAQsByCategoryQuery модель для получения FAQ по категории
type GetFAQsByCategoryQuery struct {
	Limit      int    `form:"_limit" example:"10"`
	Offset     int    `form:"_offset" example:"0"`
	SortBy     string `form:"_sort" example:"priority"`
	SortOrder  string `form:"_order" example:"desc"`
	ActiveOnly bool   `form:"activeOnly" example:"true"`
}

// GetFAQCategoriesQuery модель для получения категорий FAQ
type GetFAQCategoriesQuery struct {
	WithCounts bool `form:"withCounts" example:"false"`
}

// GetFAQCountQuery модель для получения количества FAQ
type GetFAQCountQuery struct {
	Category string `form:"category" example:"налоги"`
	IsActive bool   `form:"isActive" example:"true"`
}
