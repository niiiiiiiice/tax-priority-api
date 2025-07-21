package models

import (
	"tax-priority-api/src/application/faq/commands"
	"tax-priority-api/src/application/faq/queries"
)

// ToCreateFAQCommand преобразует HTTP-модель в команду создания FAQ
func (r *CreateFAQRequest) ToCreateFAQCommand() commands.CreateFAQCommand {
	return commands.CreateFAQCommand{
		Question: r.Question,
		Answer:   r.Answer,
		Category: r.Category,
		Priority: r.Priority,
	}
}

// ToUpdateFAQCommand преобразует HTTP-модель в команду обновления FAQ
func (r *UpdateFAQRequest) ToUpdateFAQCommand(id string) commands.UpdateFAQCommand {
	return commands.UpdateFAQCommand{
		ID:       id,
		Question: r.Question,
		Answer:   r.Answer,
		Category: r.Category,
		Priority: r.Priority,
	}
}

// ToUpdateFAQPriorityCommand преобразует HTTP-модель в команду обновления приоритета FAQ
func (r *UpdateFAQPriorityRequest) ToUpdateFAQPriorityCommand(id string) commands.UpdateFAQPriorityCommand {
	return commands.UpdateFAQPriorityCommand{
		ID:       id,
		Priority: r.Priority,
	}
}

// ToBulkDeleteFAQCommand преобразует HTTP-модель в команду массового удаления FAQ
func (r *BulkDeleteFAQRequest) ToBulkDeleteFAQCommand() commands.BulkDeleteFAQCommand {
	return commands.BulkDeleteFAQCommand{
		IDs: r.IDs,
	}
}

// ToGetFAQsByIDsQuery преобразует HTTP-модель в запрос получения FAQ по ID
func (r *GetFAQsByIDsRequest) ToGetFAQsByIDsQuery() queries.GetFAQsByIDsQuery {
	return queries.GetFAQsByIDsQuery{
		IDs: r.IDs,
	}
}

// ToGetFAQsQuery преобразует HTTP-модель в запрос получения списка FAQ
func (r *GetFAQsQuery) ToGetFAQsQuery() queries.GetFAQsQuery {
	filters := make(map[string]interface{})

	if r.Category != "" {
		filters["category"] = r.Category
	}

	// Добавляем фильтр по активности, если он был передан
	if r.IsActive {
		filters["isActive"] = r.IsActive
	}

	return queries.GetFAQsQuery{
		Limit:     r.Limit,
		Offset:    r.Offset,
		SortBy:    r.SortBy,
		SortOrder: r.SortOrder,
		Filters:   filters,
	}
}

// ToSearchFAQsQuery преобразует HTTP-модель в запрос поиска FAQ
func (r *SearchFAQsQuery) ToSearchFAQsQuery() queries.SearchFAQsQuery {
	return queries.SearchFAQsQuery{
		Query:      r.Query,
		Category:   r.Category,
		Limit:      r.Limit,
		Offset:     r.Offset,
		SortBy:     r.SortBy,
		SortOrder:  r.SortOrder,
		ActiveOnly: r.ActiveOnly,
	}
}

// ToGetFAQsByCategoryQuery преобразует HTTP-модель в запрос получения FAQ по категории
func (r *GetFAQsByCategoryQuery) ToGetFAQsByCategoryQuery(category string) queries.GetFAQsByCategoryQuery {
	return queries.GetFAQsByCategoryQuery{
		Category:   category,
		Limit:      r.Limit,
		Offset:     r.Offset,
		SortBy:     r.SortBy,
		SortOrder:  r.SortOrder,
		ActiveOnly: r.ActiveOnly,
	}
}

// ToGetFAQCategoriesQuery преобразует HTTP-модель в запрос получения категорий FAQ
func (r *GetFAQCategoriesQuery) ToGetFAQCategoriesQuery() queries.GetFAQCategoriesQuery {
	return queries.GetFAQCategoriesQuery{
		WithCounts: r.WithCounts,
	}
}

// ToGetFAQCountQuery преобразует HTTP-модель в запрос получения количества FAQ
func (r *GetFAQCountQuery) ToGetFAQCountQuery() queries.GetFAQCountQuery {
	filters := make(map[string]interface{})

	if r.Category != "" {
		filters["category"] = r.Category
	}

	if r.IsActive {
		filters["isActive"] = r.IsActive
	}

	return queries.GetFAQCountQuery{
		Filters: filters,
	}
}
