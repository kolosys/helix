package helix

import "net/http"

// Pagination contains common pagination parameters.
// Use with struct embedding for automatic binding.
//
// Example:
//
//	type ListUsersRequest struct {
//	    helix.Pagination
//	    Status string `query:"status"`
//	}
type Pagination struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Sort   string `query:"sort"`
	Order  string `query:"order"`
	Cursor string `query:"cursor"`
}

// GetPage returns the page number with a default of 1.
func (p Pagination) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetLimit returns the limit with a default and maximum.
func (p Pagination) GetLimit(defaultLimit, maxLimit int) int {
	if p.Limit <= 0 {
		return defaultLimit
	}
	if p.Limit > maxLimit {
		return maxLimit
	}
	return p.Limit
}

// GetOffset calculates the offset for SQL queries.
func (p Pagination) GetOffset(limit int) int {
	return (p.GetPage() - 1) * limit
}

// GetSort returns the sort field with a default.
func (p Pagination) GetSort(defaultSort string, allowed []string) string {
	if p.Sort == "" {
		return defaultSort
	}
	for _, s := range allowed {
		if p.Sort == s {
			return p.Sort
		}
	}
	return defaultSort
}

// GetOrder returns the order (asc/desc) with a default of desc.
func (p Pagination) GetOrder() string {
	if p.Order == "asc" {
		return "asc"
	}
	return "desc"
}

// IsAscending returns true if the order is ascending.
func (p Pagination) IsAscending() bool {
	return p.Order == "asc"
}

// PaginatedResponse wraps a list response with pagination metadata.
type PaginatedResponse[T any] struct {
	Items      []T    `json:"items"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int    `json:"total_pages"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor,omitempty"`
}

// NewPaginatedResponse creates a new paginated response.
func NewPaginatedResponse[T any](items []T, total, page, limit int) PaginatedResponse[T] {
	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}

	return PaginatedResponse[T]{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}
}

// NewCursorResponse creates a new cursor-based paginated response.
func NewCursorResponse[T any](items []T, total int, nextCursor string) PaginatedResponse[T] {
	return PaginatedResponse[T]{
		Items:      items,
		Total:      total,
		HasMore:    nextCursor != "",
		NextCursor: nextCursor,
	}
}

// BindPagination extracts pagination from the request with defaults.
func BindPagination(r *http.Request, defaultLimit, maxLimit int) Pagination {
	p := Pagination{
		Page:  QueryInt(r, "page", 1),
		Limit: QueryInt(r, "limit", defaultLimit),
		Sort:  Query(r, "sort"),
		Order: QueryDefault(r, "order", "desc"),
	}

	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = defaultLimit
	}
	if p.Limit > maxLimit {
		p.Limit = maxLimit
	}

	return p
}

// BindPaginationCtx extracts pagination from the Ctx with defaults.
func (c *Ctx) BindPagination(defaultLimit, maxLimit int) Pagination {
	return BindPagination(c.Request, defaultLimit, maxLimit)
}
