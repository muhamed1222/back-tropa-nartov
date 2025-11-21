package models

import "time"

// Image представляет изображение места
// Поддерживает две таблицы:
// 1. images - для изображений через Go API
// 2. files - для изображений через Strapi (files_related_morphs)
type Image struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URL       string    `gorm:"type:varchar(500);not null" json:"url"`
	PlaceID   uint      `gorm:"not null" json:"place_id"`
	Place     Place     `gorm:"foreignKey:PlaceID" json:"-"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:now()" json:"updated_at"`
}

// StrapiFile представляет файл из Strapi (таблица files)
type StrapiFile struct {
	ID              uint      `gorm:"primaryKey;column:id" json:"id"`
	Name            string    `gorm:"column:name" json:"name"`
	AlternativeText string    `gorm:"column:alternative_text" json:"alternativeText"`
	URL             string    `gorm:"column:url" json:"url"`
	Width           int       `gorm:"column:width" json:"width"`
	Height          int       `gorm:"column:height" json:"height"`
	Mime            string    `gorm:"column:mime" json:"mime"`
	Size            float64   `gorm:"column:size" json:"size"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (StrapiFile) TableName() string {
	return "files"
}

// StrapiFileRelation представляет связь файла с местом (таблица files_related_morphs)
type StrapiFileRelation struct {
	ID          uint    `gorm:"primaryKey"`
	FileID      uint    `gorm:"column:file_id"`
	RelatedID   uint    `gorm:"column:related_id"`
	RelatedType string  `gorm:"column:related_type"`
	Field       string  `gorm:"column:field"`
	Order       float64 `gorm:"column:order"`
}

func (StrapiFileRelation) TableName() string {
	return "files_related_morphs"
}
