package models

type QueryOptions struct {
	Pagination *PaginationParams
	SortBy     []SortBy
	Filters    map[string]any
	Includes   []string
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Filters: make(map[string]any),
	}
}

func (qo *QueryOptions) WithPagination(offset, limit int) *QueryOptions {
	qo.Pagination = &PaginationParams{Offset: offset, Limit: limit}
	return qo
}

func (qo *QueryOptions) WithSort(field string, order SortOrder) *QueryOptions {
	qo.SortBy = append(qo.SortBy, SortBy{Field: field, Order: order})
	return qo
}

func (qo *QueryOptions) WithFilter(key string, value any) *QueryOptions {
	qo.Filters[key] = value
	return qo
}

func (qo *QueryOptions) WithIncludes(includes ...string) *QueryOptions {
	qo.Includes = append(qo.Includes, includes...)
	return qo
}

func NewQueryOptionsWithFilters(filters map[string]any) *QueryOptions {
	return &QueryOptions{
		Filters: filters,
	}
}

func NewQueryOptionsWithFilter(key string, value any) *QueryOptions {
	return &QueryOptions{
		Filters: map[string]any{key: value},
	}
}
