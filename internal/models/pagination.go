package models

// PaginationParams параметры пагинации
type PaginationParams struct {
	Limit  int // Максимальное количество элементов (по умолчанию 20, максимум 100)
	Offset int // Смещение (по умолчанию 0)
}

// PaginatedResponse ответ с пагинацией
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`       // Общее количество элементов
	Limit      int         `json:"limit"`       // Лимит на страницу
	Offset     int         `json:"offset"`      // Текущее смещение
	HasMore    bool        `json:"has_more"`    // Есть ли еще элементы
	NextOffset *int        `json:"next_offset"` // Смещение для следующей страницы (null если нет)
}

