package models

import "time"

// PlaceListItem легкая версия Place для списков
type PlaceListItem struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Rating      float32   `json:"rating"`
	TypeID      uint      `json:"type_id"`
	AreaID      uint      `json:"area_id"`
	IsFavorite  *bool     `json:"is_favorite,omitempty"` // null если не проверено, иначе bool
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Минимальная информация об изображениях (только первое)
	FirstImage *string `json:"first_image,omitempty"` // URL первого изображения
}

// RouteListItem легкая версия Route для списков
type RouteListItem struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Distance    float64   `json:"distance"`
	Duration    float64   `json:"duration"`
	TypeID      uint      `json:"type_id"`
	AreaID      uint      `json:"area_id"`
	Rating      float32   `json:"rating"`
	IsFavorite  *bool     `json:"is_favorite,omitempty"` // null если не проверено, иначе bool
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Минимальная информация о типе и районе
	TypeName string `json:"type_name,omitempty"`
	AreaName string `json:"area_name,omitempty"`
}

