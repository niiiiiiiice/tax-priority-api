package models

type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Offset     int   `json:"offset"`
	Limit      int   `json:"limit"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
	TotalPages int   `json:"totalPages"`
}
