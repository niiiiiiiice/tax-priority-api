package repositories

import (
	"context"
	"fmt"
	"strings"

	sharedModels "tax-priority-api/src/application/shared/models"
	"tax-priority-api/src/domain/entities"
	"tax-priority-api/src/domain/repositories"
	"tax-priority-api/src/infrastructure/persistence"
	infraModels "tax-priority-api/src/infrastructure/persistence/models"

	"gorm.io/gorm"
)

// FAQRepositoryImpl реализация FAQRepository для GORM
type FAQRepositoryImpl struct {
	db *gorm.DB
}

// NewFAQRepository создает новый репозиторий FAQ
func NewFAQRepository(db *gorm.DB) repositories.FAQRepository {
	return &FAQRepositoryImpl{
		db: db,
	}
}

// Create создает новую FAQ
func (r *FAQRepositoryImpl) Create(ctx context.Context, faq *entities.FAQ) error {
	model := infraModels.NewFAQModelFromEntity(faq)

	result := r.db.WithContext(ctx).Create(model)
	if result.Error != nil {
		return persistence.NewInternalError("failed to create FAQ", result.Error)
	}

	return nil
}

// CreateBatch создает несколько FAQ
func (r *FAQRepositoryImpl) CreateBatch(ctx context.Context, faqs []*entities.FAQ) (*sharedModels.BulkOperationResult, error) {
	models := infraModels.FromFAQEntities(faqs)

	result := r.db.WithContext(ctx).CreateInBatches(models, 100)
	if result.Error != nil {
		return &sharedModels.BulkOperationResult{
			SuccessCount: 0,
			FailureCount: len(faqs),
			Errors:       []error{result.Error},
		}, persistence.NewInternalError("failed to create FAQ batch", result.Error)
	}

	return &sharedModels.BulkOperationResult{
		SuccessCount: len(faqs),
		FailureCount: 0,
		Errors:       []error{},
	}, nil
}

// FindByID находит FAQ по ID
func (r *FAQRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.FAQ, error) {
	var model infraModels.FAQModel

	result := r.db.WithContext(ctx).First(&model, "id = ?", id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, persistence.NewNotFoundError(fmt.Sprintf("FAQ with ID %s not found", id), result.Error)
		}
		return nil, persistence.NewInternalError("failed to find FAQ", result.Error)
	}

	return model.ToEntity(), nil
}

// FindByIDs находит FAQ по списку ID
func (r *FAQRepositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]*entities.FAQ, error) {
	var models []*infraModels.FAQModel

	result := r.db.WithContext(ctx).Find(&models, "id IN ?", ids)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find FAQs by IDs", result.Error)
	}

	return infraModels.ToFAQEntities(models), nil
}

// Update обновляет FAQ
func (r *FAQRepositoryImpl) Update(ctx context.Context, faq *entities.FAQ) error {
	model := infraModels.NewFAQModelFromEntity(faq)

	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return persistence.NewInternalError("failed to update FAQ", result.Error)
	}

	if result.RowsAffected == 0 {
		return persistence.NewNotFoundError(fmt.Sprintf("FAQ with ID %s not found", faq.ID), nil)
	}

	return nil
}

// UpdateBatch обновляет несколько FAQ
func (r *FAQRepositoryImpl) UpdateBatch(ctx context.Context, faqs []*entities.FAQ) (*sharedModels.BulkOperationResult, error) {
	successCount := 0
	failureCount := 0
	errors := make([]error, 0)

	for _, faq := range faqs {
		if err := r.Update(ctx, faq); err != nil {
			failureCount++
			errors = append(errors, err)
		} else {
			successCount++
		}
	}

	return &sharedModels.BulkOperationResult{
		SuccessCount: successCount,
		FailureCount: failureCount,
		Errors:       errors,
	}, nil
}

// UpdateFields обновляет определенные поля FAQ
func (r *FAQRepositoryImpl) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&infraModels.FAQModel{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return persistence.NewInternalError("failed to update FAQ fields", result.Error)
	}

	if result.RowsAffected == 0 {
		return persistence.NewNotFoundError(fmt.Sprintf("FAQ with ID %s not found", id), nil)
	}

	return nil
}

// Delete удаляет FAQ
func (r *FAQRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&infraModels.FAQModel{}, "id = ?", id)
	if result.Error != nil {
		return persistence.NewInternalError("failed to delete FAQ", result.Error)
	}

	if result.RowsAffected == 0 {
		return persistence.NewNotFoundError(fmt.Sprintf("FAQ with ID %s not found", id), nil)
	}

	return nil
}

// DeleteBatch удаляет несколько FAQ
func (r *FAQRepositoryImpl) DeleteBatch(ctx context.Context, ids []string) (*sharedModels.BulkOperationResult, error) {
	result := r.db.WithContext(ctx).Delete(&infraModels.FAQModel{}, "id IN ?", ids)
	if result.Error != nil {
		return &sharedModels.BulkOperationResult{
			SuccessCount: 0,
			FailureCount: len(ids),
			Errors:       []error{result.Error},
		}, persistence.NewInternalError("failed to delete FAQ batch", result.Error)
	}

	return &sharedModels.BulkOperationResult{
		SuccessCount: int(result.RowsAffected),
		FailureCount: len(ids) - int(result.RowsAffected),
		Errors:       []error{},
	}, nil
}

// FindAll находит все FAQ с опциями
func (r *FAQRepositoryImpl) FindAll(ctx context.Context, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	var models []*infraModels.FAQModel
	query := r.db.WithContext(ctx)

	// Применяем фильтры, сортировку и пагинацию
	query = r.applyQueryOptions(query, opts)

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find FAQs", result.Error)
	}

	return infraModels.ToFAQEntities(models), nil
}

// FindWithPagination находит FAQ с пагинацией
func (r *FAQRepositoryImpl) FindWithPagination(ctx context.Context, opts *sharedModels.QueryOptions) (*sharedModels.PaginatedResult[*entities.FAQ], error) {
	if opts == nil || opts.Pagination == nil {
		return nil, persistence.NewInvalidInputError("pagination options are required", nil)
	}

	var models []*infraModels.FAQModel
	var total int64

	// Подсчитываем общее количество
	countQuery := r.db.WithContext(ctx).Model(&infraModels.FAQModel{})
	countQuery = r.applyFilters(countQuery, opts.Filters)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, persistence.NewInternalError("failed to count FAQs", err)
	}

	// Получаем данные
	query := r.db.WithContext(ctx)
	query = r.applyQueryOptions(query, opts)

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find FAQs with pagination", result.Error)
	}

	entity := infraModels.ToFAQEntities(models)

	// Вычисляем пагинацию
	totalPages := int((total + int64(opts.Pagination.Limit) - 1) / int64(opts.Pagination.Limit))
	hasNext := opts.Pagination.Offset+opts.Pagination.Limit < int(total)
	hasPrev := opts.Pagination.Offset > 0

	return &sharedModels.PaginatedResult[*entities.FAQ]{
		Items:      entity,
		Total:      total,
		Offset:     opts.Pagination.Offset,
		Limit:      opts.Pagination.Limit,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		TotalPages: totalPages,
	}, nil
}

// FindByCategory находит FAQ по категории
func (r *FAQRepositoryImpl) FindByCategory(ctx context.Context, category string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	var models []*infraModels.FAQModel
	query := r.db.WithContext(ctx).Where("category = ?", category)

	query = r.applyQueryOptions(query, opts)

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find FAQs by category", result.Error)
	}

	return infraModels.ToFAQEntities(models), nil
}

// FindActive находит активные FAQ
func (r *FAQRepositoryImpl) FindActive(ctx context.Context, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	var models []*infraModels.FAQModel
	query := r.db.WithContext(ctx).Where("is_active = ?", true)

	query = r.applyQueryOptions(query, opts)

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find active FAQs", result.Error)
	}

	return infraModels.ToFAQEntities(models), nil
}

// FindByPriority находит FAQ по приоритету
func (r *FAQRepositoryImpl) FindByPriority(ctx context.Context, minPriority int, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	var models []*infraModels.FAQModel
	query := r.db.WithContext(ctx).Where("priority >= ?", minPriority)

	query = r.applyQueryOptions(query, opts)

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to find FAQs by priority", result.Error)
	}

	return infraModels.ToFAQEntities(models), nil
}

// Search выполняет поиск FAQ
func (r *FAQRepositoryImpl) Search(ctx context.Context, searchQuery string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	var models []*infraModels.FAQModel
	query := r.db.WithContext(ctx).Where(
		"LOWER(question) LIKE ? OR LOWER(answer) LIKE ? OR LOWER(category) LIKE ?",
		"%"+strings.ToLower(searchQuery)+"%",
		"%"+strings.ToLower(searchQuery)+"%",
		"%"+strings.ToLower(searchQuery)+"%",
	)

	query = r.applyQueryOptions(query, opts)

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to search FAQs", result.Error)
	}

	return infraModels.ToFAQEntities(models), nil
}

// SearchByCategory выполняет поиск FAQ по категории
func (r *FAQRepositoryImpl) SearchByCategory(ctx context.Context, searchQuery string, category string, opts *sharedModels.QueryOptions) ([]*entities.FAQ, error) {
	var models []*infraModels.FAQModel
	query := r.db.WithContext(ctx).Where(
		"category = ? AND (LOWER(question) LIKE ? OR LOWER(answer) LIKE ?)",
		category,
		"%"+strings.ToLower(searchQuery)+"%",
		"%"+strings.ToLower(searchQuery)+"%",
	)

	query = r.applyQueryOptions(query, opts)

	result := query.Find(&models)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to search FAQs by category", result.Error)
	}

	return infraModels.ToFAQEntities(models), nil
}

// Count подсчитывает количество FAQ
func (r *FAQRepositoryImpl) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&infraModels.FAQModel{})

	query = r.applyFilters(query, filters)

	result := query.Count(&count)
	if result.Error != nil {
		return 0, persistence.NewInternalError("failed to count FAQs", result.Error)
	}

	return count, nil
}

// CountByCategory подсчитывает количество FAQ по категории
func (r *FAQRepositoryImpl) CountByCategory(ctx context.Context, category string) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&infraModels.FAQModel{}).Where("category = ?", category).Count(&count)
	if result.Error != nil {
		return 0, persistence.NewInternalError("failed to count FAQs by category", result.Error)
	}

	return count, nil
}

// CountActive подсчитывает количество активных FAQ
func (r *FAQRepositoryImpl) CountActive(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&infraModels.FAQModel{}).Where("is_active = ?", true).Count(&count)
	if result.Error != nil {
		return 0, persistence.NewInternalError("failed to count active FAQs", result.Error)
	}

	return count, nil
}

// Exists проверяет существование FAQ
func (r *FAQRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&infraModels.FAQModel{}).Where("id = ?", id).Count(&count)
	if result.Error != nil {
		return false, persistence.NewInternalError("failed to check FAQ existence", result.Error)
	}

	return count > 0, nil
}

// ExistsByQuestion проверяет существование FAQ по вопросу
func (r *FAQRepositoryImpl) ExistsByQuestion(ctx context.Context, question string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&infraModels.FAQModel{}).Where("question = ?", question).Count(&count)
	if result.Error != nil {
		return false, persistence.NewInternalError("failed to check FAQ existence by question", result.Error)
	}

	return count > 0, nil
}

// GetCategories получает все категории FAQ
func (r *FAQRepositoryImpl) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	result := r.db.WithContext(ctx).Model(&infraModels.FAQModel{}).Distinct("category").Pluck("category", &categories)
	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to get FAQ categories", result.Error)
	}

	return categories, nil
}

// GetCategoriesWithCounts получает категории FAQ с количеством
func (r *FAQRepositoryImpl) GetCategoriesWithCounts(ctx context.Context) (map[string]int64, error) {
	type categoryCount struct {
		Category string
		Count    int64
	}

	var results []categoryCount
	result := r.db.WithContext(ctx).Model(&infraModels.FAQModel{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Scan(&results)

	if result.Error != nil {
		return nil, persistence.NewInternalError("failed to get FAQ categories with counts", result.Error)
	}

	categoryCounts := make(map[string]int64)
	for _, r := range results {
		categoryCounts[r.Category] = r.Count
	}

	return categoryCounts, nil
}

// WithTransaction выполняет операцию в транзакции
func (r *FAQRepositoryImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Создаем новый репозиторий с транзакционной БД
		txRepo := &FAQRepositoryImpl{db: tx}

		// Создаем новый контекст с транзакционным репозиторием
		txCtx := context.WithValue(ctx, "tx_repo", txRepo)

		return fn(txCtx)
	})
}

// Вспомогательные методы

func (r *FAQRepositoryImpl) applyQueryOptions(query *gorm.DB, opts *sharedModels.QueryOptions) *gorm.DB {
	if opts == nil {
		return query
	}

	// Применяем фильтры
	query = r.applyFilters(query, opts.Filters)

	// Применяем сортировку
	query = r.applySorting(query, opts.SortBy)

	// Применяем пагинацию
	if opts.Pagination != nil {
		query = query.Offset(opts.Pagination.Offset).Limit(opts.Pagination.Limit)
	}

	return query
}

func (r *FAQRepositoryImpl) applyFilters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	if filters == nil {
		return query
	}

	for key, value := range filters {
		if value != nil {
			// Преобразуем имя поля в имя колонки БД
			dbColumnName := r.mapFieldToDBColumn(key)
			query = query.Where(dbColumnName+" = ?", value)
		}
	}

	return query
}

func (r *FAQRepositoryImpl) applySorting(query *gorm.DB, sortBy []sharedModels.SortBy) *gorm.DB {
	if sortBy == nil {
		return query
	}

	for _, sort := range sortBy {
		// Преобразуем имя поля в имя колонки БД
		dbColumnName := r.mapFieldToDBColumn(sort.Field)
		query = query.Order(fmt.Sprintf("%s %s", dbColumnName, sort.Order))
	}

	return query
}

// mapFieldToDBColumn преобразует имя поля entity в имя колонки БД
func (r *FAQRepositoryImpl) mapFieldToDBColumn(fieldName string) string {
	// Маппинг полей entity на колонки БД
	fieldMapping := map[string]string{
		"id":        "id",
		"question":  "question",
		"answer":    "answer",
		"category":  "category",
		"isActive":  "is_active",
		"priority":  "priority",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
		"deletedAt": "deleted_at",
	}

	// Если есть прямое соответствие, используем его
	if dbColumn, exists := fieldMapping[fieldName]; exists {
		return dbColumn
	}

	// Иначе конвертируем camelCase в snake_case
	return r.camelToSnake(fieldName)
}

// camelToSnake преобразует camelCase в snake_case
func (r *FAQRepositoryImpl) camelToSnake(str string) string {
	var result strings.Builder
	for i, char := range str {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}
