package repositories

import (
	"context"
	"strings"

	sharedModels "tax-priority-api/src/application/models"
	"tax-priority-api/src/application/repositories"
	"tax-priority-api/src/domain/entities"
)

// FAQRepositoryImpl реализация FAQRepository для GORM
type FAQRepositoryImpl struct {
	generic repositories.GenericRepository[*entities.FAQ, string]
}

// NewFAQRepository создает новый репозиторий FAQ
func NewFAQRepository(generic repositories.GenericRepository[*entities.FAQ, string]) repositories.FAQRepository {
	return &FAQRepositoryImpl{generic}
}

// Delegate all GenericRepository methods
func (r *FAQRepositoryImpl) Create(ctx context.Context, entity *entities.FAQ) error {
	return r.generic.Create(ctx, entity)
}
func (r *FAQRepositoryImpl) CreateBatch(ctx context.Context, entities []*entities.FAQ) (*sharedModels.BulkOperationResult, error) {
	return r.generic.CreateBatch(ctx, entities)
}

func (r *FAQRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.FAQ, error) {
	return r.generic.FindByID(ctx, id)
}
func (r *FAQRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.FAQ, error) {
	return r.generic.FindByIDs(ctx, ids)
}

func (r *FAQRepositoryImpl) Update(ctx context.Context, entity *entities.FAQ) error {
	return r.generic.Update(ctx, entity)
}
func (r *FAQRepositoryImpl) UpdateBatch(ctx context.Context, entities []*entities.FAQ) (*sharedModels.BulkOperationResult, error) {
	return r.generic.UpdateBatch(ctx, entities)
}
func (r *FAQRepositoryImpl) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	return r.generic.UpdateFields(ctx, id, fields)
}

func (r *FAQRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.generic.Delete(ctx, id)
}
func (r *FAQRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*sharedModels.BulkOperationResult, error) {
	return r.generic.DeleteBatch(ctx, ids)
}
func (r *FAQRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	return r.generic.SoftDelete(ctx, id)
}

func (r *FAQRepositoryImpl) FindAll(ctx context.Context, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	return r.generic.FindAll(ctx, opts)
}
func (r *FAQRepositoryImpl) FindOne(ctx context.Context, opts *sharedModels.QueryOptions) (*entities.FAQ, error) {
	return r.generic.FindOne(ctx, opts)
}
func (r *FAQRepositoryImpl) FindWithPagination(ctx context.Context, opts *sharedModels.QueryOptions) (*sharedModels.PaginatedResult[*entities.FAQ], error) {
	return r.generic.FindWithPagination(ctx, opts)
}

func (r *FAQRepositoryImpl) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	return r.generic.Count(ctx, filters)
}
func (r *FAQRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	return r.generic.Exists(ctx, id)
}
func (r *FAQRepositoryImpl) ExistsByFields(ctx context.Context, filters map[string]interface{}) (bool, error) {
	return r.generic.ExistsByFields(ctx, filters)
}

func (r *FAQRepositoryImpl) WithTransaction(ctx context.Context, fn repositories.TransactionFunc) error {
	return r.generic.WithTransaction(ctx, repositories.TransactionFunc(fn))
}

func (r *FAQRepositoryImpl) Refresh(ctx context.Context, entity *entities.FAQ) error {
	return r.generic.Refresh(ctx, entity)
}
func (r *FAQRepositoryImpl) Clear(ctx context.Context) error {
	return r.generic.Clear(ctx)
}

// FindByCategory находит FAQ по категории
func (r *FAQRepositoryImpl) FindByCategory(ctx context.Context, category string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	localOpts := &sharedModels.QueryOptions{}
	if opts != nil {
		*localOpts = *opts
	}
	if localOpts.Filters == nil {
		localOpts.Filters = make(map[string]interface{})
	}
	localOpts.Filters["category"] = category
	return r.FindAll(ctx, localOpts)
}

// FindActive находит активные FAQ
func (r *FAQRepositoryImpl) FindActive(ctx context.Context, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	localOpts := &sharedModels.QueryOptions{}
	if opts != nil {
		*localOpts = *opts
	}
	if localOpts.Filters == nil {
		localOpts.Filters = make(map[string]interface{})
	}
	localOpts.Filters["is_active"] = true

	return r.FindAll(ctx, localOpts)
}

// FindByPriority находит FAQ по приоритету
func (r *FAQRepositoryImpl) FindByPriority(ctx context.Context, minPriority int, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	localOpts := &sharedModels.QueryOptions{}
	if opts != nil {
		*localOpts = *opts
	}
	if localOpts.Filters == nil {
		localOpts.Filters = make(map[string]interface{})
	}
	localOpts.Filters["priority >="] = minPriority
	return r.FindAll(ctx, localOpts)
}

// Search выполняет поиск FAQ
func (r *FAQRepositoryImpl) Search(ctx context.Context, searchQuery string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	localOpts := &sharedModels.QueryOptions{}
	if opts != nil {
		*localOpts = *opts
	}
	if localOpts.Filters == nil {
		localOpts.Filters = make(map[string]interface{})
	}

	// Получаем все FAQ для последующей фильтрации
	allFAQs, err := r.FindAll(ctx, &sharedModels.QueryOptions{
		Filters:    localOpts.Filters,
		SortBy:     localOpts.SortBy,
		Pagination: localOpts.Pagination,
		Includes:   localOpts.Includes,
	})
	if err != nil {
		return nil, err
	}

	// Фильтруем результаты по поисковому запросу
	var filteredFAQs []*entities.FAQ
	searchQueryLower := strings.ToLower(searchQuery)

	for _, faq := range allFAQs {
		searchableText := faq.GetSearchableText()
		if strings.Contains(searchableText, searchQueryLower) {
			filteredFAQs = append(filteredFAQs, faq)
		}
	}

	return filteredFAQs, nil
}

// SearchByCategory выполняет поиск FAQ по категории
func (r *FAQRepositoryImpl) SearchByCategory(ctx context.Context, searchQuery string, category string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	localOpts := &sharedModels.QueryOptions{}
	if opts != nil {
		*localOpts = *opts
	}
	if localOpts.Filters == nil {
		localOpts.Filters = make(map[string]interface{})
	}

	// Добавляем фильтр по категории
	localOpts.Filters["category"] = category

	// Используем метод Search с дополнительным фильтром по категории
	return r.Search(ctx, searchQuery, localOpts)
}

// CountByCategory подсчитывает количество FAQ по категории
func (r *FAQRepositoryImpl) CountByCategory(ctx context.Context, category string) (int64, error) {
	return r.Count(ctx, map[string]interface{}{"category": category})
}

// CountActive подсчитывает количество активных FAQ
func (r *FAQRepositoryImpl) CountActive(ctx context.Context) (int64, error) {
	return r.Count(ctx, map[string]interface{}{"is_active": true})
}

// ExistsByQuestion проверяет существование FAQ по вопросу
func (r *FAQRepositoryImpl) ExistsByQuestion(ctx context.Context, question string) (bool, error) {
	return r.ExistsByFields(ctx, map[string]interface{}{"question": question})
}

// GetCategories получает все категории FAQ
func (r *FAQRepositoryImpl) GetCategories(ctx context.Context) ([]string, error) {
	// Получаем все FAQ
	allFAQs, err := r.FindAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Собираем уникальные категории
	categorySet := make(map[string]bool)
	for _, faq := range allFAQs {
		if faq.Category != "" {
			categorySet[faq.Category] = true
		}
	}

	// Преобразуем в слайс
	categories := make([]string, 0, len(categorySet))
	for category := range categorySet {
		categories = append(categories, category)
	}

	return categories, nil
}

// GetCategoriesWithCounts получает категории FAQ с количеством
func (r *FAQRepositoryImpl) GetCategoriesWithCounts(ctx context.Context) (map[string]int64, error) {
	// Получаем все FAQ
	allFAQs, err := r.FindAll(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Подсчитываем количество FAQ в каждой категории
	categoryCounts := make(map[string]int64)
	for _, faq := range allFAQs {
		if faq.Category != "" {
			categoryCounts[faq.Category]++
		}
	}

	return categoryCounts, nil
}
