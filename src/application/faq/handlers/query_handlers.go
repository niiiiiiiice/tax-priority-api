package handlers

import (
	"tax-priority-api/src/application/faq/queries"
	"tax-priority-api/src/application/repositories"
)

type FAQQueryHandlers struct {
	GetActive     *queries.GetActiveFAQsQueryHandler
	GetByCategory *queries.GetFAQsByCategoryQueryHandler
	GetByID       *queries.GetFAQByIDQueryHandler
	GetByIDs      *queries.GetFAQsByIDsQueryHandler
	GetByPriority *queries.GetFAQsByPriorityQueryHandler
	GetCategories *queries.GetFAQCategoriesQueryHandler
	GetCount      *queries.GetFAQCountQueryHandler
	GetMany       *queries.GetFAQsQueryHandler
	Search        *queries.SearchFAQsQueryHandler
}

func NewFAQQueryHandlers(repo repositories.FAQRepository) *FAQQueryHandlers {
	return &FAQQueryHandlers{
		GetActive:     queries.NewGetActiveFAQsQueryHandler(repo),
		GetByCategory: queries.NewGetFAQsByCategoryQueryHandler(repo),
		GetByID:       queries.NewGetFAQByIDQueryHandler(repo),
		GetByIDs:      queries.NewGetFAQsByIDsQueryHandler(repo),
		GetByPriority: queries.NewGetFAQsByPriorityQueryHandler(repo),
		GetCategories: queries.NewGetFAQCategoriesQueryHandler(repo),
		GetCount:      queries.NewGetFAQCountQueryHandler(repo),
		GetMany:       queries.NewGetFAQsQueryHandler(repo),
		Search:        queries.NewSearchFAQsQueryHandler(repo),
	}
}
