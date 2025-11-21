package models

import "math"

// PaginationParams параметры пагинации
type PaginationParams struct {
	Page  int `form:"page" json:"page" binding:"min=1"`     // Номер страницы (начиная с 1)
	Limit int `form:"limit" json:"limit" binding:"min=1,max=100"` // Элементов на странице (макс 100)
}

// NewPaginationParams создает параметры пагинации с defaults
func NewPaginationParams(page, limit int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20 // default
	}
	if limit > 100 {
		limit = 100 // max
	}
	return PaginationParams{
		Page:  page,
		Limit: limit,
	}
}

// GetOffset вычисляет offset из page и limit
func (p PaginationParams) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// PaginatedResponse ответ с пагинацией
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`       // Общее количество элементов
	Page       int         `json:"page"`        // Текущая страница
	Limit      int         `json:"limit"`       // Лимит на страницу
	TotalPages int         `json:"total_pages"` // Общее количество страниц
	HasMore    bool        `json:"has_more"`    // Есть ли еще страницы
}

// NewPaginatedResponse создает paginated response
func NewPaginatedResponse(data interface{}, total int64, params PaginationParams) PaginatedResponse {
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	hasMore := params.Page < totalPages

	return PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}
}

