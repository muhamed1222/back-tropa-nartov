package places

import (
	"fmt"
	"tropa-nartov-backend/internal/models"

	"gorm.io/gorm"
)

// Service управляет точками
type Service struct {
	db *gorm.DB
}

// NewService создаёт новый сервис для точек
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Create создаёт новую точку
func (s *Service) Create(place *models.Place) error {
	if err := s.db.Create(place).Error; err != nil {
		return fmt.Errorf("ошибка создания точки: %v", err)
	}
	return nil
}

// Update обновляет точку
func (s *Service) Update(id uint, place *models.Place) error {
	var existing models.Place
	if err := s.db.Where("id = ? AND is_active = ?", id, true).First(&existing).Error; err != nil {
		return fmt.Errorf("точка не найдена")
	}

	place.ID = id
	if err := s.db.Save(place).Error; err != nil {
		return fmt.Errorf("ошибка обновления точки: %v", err)
	}
	return nil
}

// Delete выполняет soft delete точки
func (s *Service) Delete(id uint) error {
	var place models.Place
	if err := s.db.Where("id = ? AND is_active = ?", id, true).First(&place).Error; err != nil {
		return fmt.Errorf("точка не найдена")
	}

	place.IsActive = false
	if err := s.db.Save(&place).Error; err != nil {
		return fmt.Errorf("ошибка удаления точки: %v", err)
	}
	return nil
}

// LoadImagesForPlace загружает изображения из обеих таблиц (images и files через Strapi)
func (s *Service) LoadImagesForPlace(place *models.Place) error {
	fmt.Printf("\n[LoadImagesForPlace START] Place ID: %d, Name: %s\n", place.ID, place.Name)
	
	// 1. Загружаем из таблицы images (наши изображения)
	var images []models.Image
	result := s.db.Where("place_id = ? AND is_active = ?", place.ID, true).Find(&images)
	fmt.Printf("[LoadImagesForPlace] Place ID: %d, images из таблицы images: %d, ошибка: %v\n", place.ID, len(images), result.Error)
	
	// 2. Загружаем из таблицы files (изображения из Strapi)
	var strapiFiles []models.StrapiFile
	result2 := s.db.Table("files").
		Joins("JOIN files_related_morphs ON files_related_morphs.file_id = files.id").
		Where("files_related_morphs.related_id = ?", place.ID).
		Where("files_related_morphs.related_type = ?", "api::place.place").
		Where("files_related_morphs.field = ?", "images").
		Find(&strapiFiles)
	fmt.Printf("[LoadImagesForPlace] Place ID: %d, файлов из Strapi: %d, ошибка: %v\n", place.ID, len(strapiFiles), result2.Error)
	
	// 3. Конвертируем Strapi файлы в формат Image
	strapiBaseURL := "http://localhost:1337" // TODO: вынести в config
	for _, file := range strapiFiles {
		// Если URL начинается с /, добавляем базовый URL Strapi
		imageURL := file.URL
		if len(imageURL) > 0 && imageURL[0] == '/' {
			imageURL = strapiBaseURL + imageURL
		}
		
		images = append(images, models.Image{
			ID:        file.ID,
			URL:       imageURL,
			PlaceID:   place.ID,
			IsActive:  true,
			CreatedAt: file.CreatedAt,
			UpdatedAt: file.UpdatedAt,
		})
	}
	
	place.Images = images
	fmt.Printf("[LoadImagesForPlace END] Place ID: %d, TOTAL images установлено: %d\n\n", place.ID, len(place.Images))
	return nil
}

// GetByID получает точку по ID с отзывами
func (s *Service) GetByID(id uint) (*models.Place, error) {
	var place models.Place
	if err := s.db.Preload("Reviews").Preload("Reviews.User").
		Where("id = ? AND is_active = ?", id, true).First(&place).Error; err != nil {
		return nil, fmt.Errorf("точка не найдена")
	}
	
	// Загружаем изображения из обеих источников
	s.LoadImagesForPlace(&place)
	
	return &place, nil
}

// List возвращает список точек с фильтрами и пагинацией
func (s *Service) List(categoryIDs, typeIDs, areaIDs, tagIDs []uint, pagination models.PaginationParams) ([]models.Place, int64, error) {
	limit := pagination.Limit
	offset := pagination.GetOffset()
	var places []models.Place
	var total int64

	// Базовый запрос для подсчета общего количества
	countQuery := s.db.Model(&models.Place{}).Where("is_active = ?", true)

	// Базовый запрос для получения данных (без Preload для производительности в списках)
	query := s.db.Where("is_active = ?", true)

	// Фильтр по типам (используем новое поле Type вместо TypeID)
	if len(typeIDs) > 0 {
		query = query.Where("type_id IN ?", typeIDs)
		countQuery = countQuery.Where("type_id IN ?", typeIDs)
	}

	// Фильтр по районам
	if len(areaIDs) > 0 {
		query = query.Where("area_id IN ?", areaIDs)
		countQuery = countQuery.Where("area_id IN ?", areaIDs)
	}

	// Фильтр по категориям (через связь many-to-many)
	if len(categoryIDs) > 0 {
		joinQuery := "JOIN place_categories ON place_categories.place_id = places.id"
		query = query.Joins(joinQuery).Where("place_categories.category_id IN ?", categoryIDs)
		countQuery = countQuery.Joins(joinQuery).Where("place_categories.category_id IN ?", categoryIDs)
	}

	// Фильтр по тегам (если используется)
	if len(tagIDs) > 0 {
		joinQuery := "JOIN place_tags ON place_tags.place_id = places.id"
		query = query.Joins(joinQuery).Where("place_tags.tag_id IN ?", tagIDs)
		countQuery = countQuery.Joins(joinQuery).Where("place_tags.tag_id IN ?", tagIDs)
	}

	// Подсчитываем общее количество
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("ошибка подсчета точек: %v", err)
	}

	// Применяем пагинацию
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	// Выполняем запрос (без Preload для списков - загружаем только базовые данные)
	if err := query.Order("created_at DESC").Find(&places).Error; err != nil {
		return nil, 0, fmt.Errorf("ошибка получения списка точек: %v", err)
	}

	return places, total, nil
}

// ListFull возвращает полный список точек с Preload (для обратной совместимости)
func (s *Service) ListFull(categoryIDs, typeIDs, areaIDs, tagIDs []uint) ([]models.Place, error) {
	places, _, err := s.List(categoryIDs, typeIDs, areaIDs, tagIDs, models.PaginationParams{Page: 1, Limit: 0})
	if err != nil {
		return nil, err
	}
	
	// Загружаем связанные данные для каждого места
	for i := range places {
		// Загружаем reviews
		s.db.Preload("Reviews").Preload("Reviews.User").First(&places[i], places[i].ID)
		
		// Загружаем изображения из обеих источников
		s.LoadImagesForPlace(&places[i])
	}
	
	return places, nil
}
