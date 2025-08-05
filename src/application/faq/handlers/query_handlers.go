package handlers

import (
	"tax-priority-api/src/application/faq/queries"
	"tax-priority-api/src/application/repositories"
)

type FAQQueryHandlers struct {
	GetByID  *queries.GetFAQByIDQueryHandler
	GetByIDs *queries.GetFAQsByIDsQueryHandler
	GetCount *queries.GetFAQCountQueryHandler
	GetMany  *queries.GetFAQsQueryHandler
}

func NewFAQQueryHandlers(repo repositories.CachedFAQRepository) *FAQQueryHandlers {
	return &FAQQueryHandlers{
		GetByID:  queries.NewGetFAQByIDQueryHandler(repo),
		GetByIDs: queries.NewGetFAQsByIDsQueryHandler(repo),
		GetCount: queries.NewGetFAQCountQueryHandler(repo),
		GetMany:  queries.NewGetFAQsQueryHandler(repo),
	}
}
