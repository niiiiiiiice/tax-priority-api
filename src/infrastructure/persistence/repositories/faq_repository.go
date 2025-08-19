package repositories

import (
	"context"
	sharedModels "tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
)

type FAQRepositoryImpl struct {
	repositories.GenericRepository[*entities.FAQ, string]
}

func NewFAQRepository(generic repositories.GenericRepository[*entities.FAQ, string]) repositories.FAQRepository {
	return &FAQRepositoryImpl{generic}
}

func (r *FAQRepositoryImpl) GetCategories(ctx context.Context, withCounts bool) ([]string, map[string]int64, error) {
	opts := &sharedModels.QueryOptions{
		Filters: map[string]interface{}{
			"isActive": true,
		},
	}

	faqs, err := r.GenericRepository.FindAll(ctx, opts)
	if err != nil {
		return nil, nil, err
	}

	categoryMap := make(map[string]int64)

	for _, faq := range faqs {
		if faq.Category != "" {
			categoryMap[faq.Category]++
		}
	}

	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}

	if !withCounts {
		return categories, nil, nil
	}

	return categories, categoryMap, nil
}
