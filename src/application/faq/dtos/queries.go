package dtos

import (
	"tax-priority-api/src/application/models"
	"tax-priority-api/src/domain/entities"
	"time"
)

type QueryResult struct {
	FAQ            *entities.FAQ                          `json:"faq,omitempty"`
	FAQs           []*entities.FAQ                        `json:"faqs,omitempty"`
	Paginated      *models.PaginatedResult[*entities.FAQ] `json:"paginated,omitempty"`
	Count          int64                                  `json:"count,omitempty"`
	Categories     []string                               `json:"categories,omitempty"`
	CategoryCounts map[string]int64                       `json:"categoryCounts,omitempty"`
	Success        bool                                   `json:"success"`
	Message        string                                 `json:"message,omitempty"`
	Error          string                                 `json:"error,omitempty"`
	Timestamp      time.Time                              `json:"timestamp"`
}

type FAQResponse struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	Category  string    `json:"category"`
	IsActive  bool      `json:"isActive"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PaginatedFAQResponse struct {
	Items      []FAQResponse `json:"items"`
	Total      int64         `json:"total"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
	HasNext    bool          `json:"hasNext"`
	HasPrev    bool          `json:"hasPrev"`
	TotalPages int           `json:"totalPages"`
}

type CategoryResponse struct {
	Name  string `json:"name"`
	Count int64  `json:"count,omitempty"`
}

func ToFAQResponse(faq *entities.FAQ) FAQResponse {
	return FAQResponse{
		ID:        faq.ID,
		Question:  faq.Question,
		Answer:    faq.Answer,
		Category:  faq.Category,
		IsActive:  faq.IsActive,
		Priority:  faq.Priority,
		CreatedAt: faq.CreatedAt,
		UpdatedAt: faq.UpdatedAt,
	}
}

func ToFAQResponses(faqs []*entities.FAQ) []FAQResponse {
	responses := make([]FAQResponse, len(faqs))
	for i, faq := range faqs {
		responses[i] = ToFAQResponse(faq)
	}
	return responses
}

func ToPaginatedFAQResponse(paginated *models.PaginatedResult[*entities.FAQ]) PaginatedFAQResponse {
	return PaginatedFAQResponse{
		Items:      ToFAQResponses(paginated.Items),
		Total:      paginated.Total,
		Offset:     paginated.Offset,
		Limit:      paginated.Limit,
		HasNext:    paginated.HasNext,
		HasPrev:    paginated.HasPrev,
		TotalPages: paginated.TotalPages,
	}
}
