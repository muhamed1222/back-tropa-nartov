package routes

import (
	"strconv"
	"tropa-nartov-backend/internal/auth"
	"tropa-nartov-backend/internal/config"
	"tropa-nartov-backend/internal/models"
	"tropa-nartov-backend/internal/places"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PlaceRequest структура для создания/обновления точки
type PlaceRequest struct {
	Name          string  `json:"name" binding:"required"`
	Type          string  `json:"type" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	Overview      string  `json:"overview"`
	History       string  `json:"history"`
	Address       string  `json:"address" binding:"required"`
	Hours         string  `json:"hours"`
	Weekend       string  `json:"weekend"`
	Entry         string  `json:"entry"`
	Contacts      string  `json:"contacts"`
	ContactsEmail string  `json:"contacts_email"`
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
	Rating        float32 `json:"rating"`

	// Старые поля для совместимости
	OpeningHours string `json:"opening_hours"`
	TypeID       uint   `json:"type_id"`
	AreaID       uint   `json:"area_id"`
	CategoryIDs  []uint `json:"category_ids"`
}

// SetupPlaceRoutes настраивает маршруты для точек
func SetupPlaceRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	placeService := places.NewService(db)

	// Группа маршрутов для точек
	placesGroup := r.Group("/places")
	{
		// Список точек с фильтрами и пагинацией
		placesGroup.GET("", func(c *gin.Context) {
			// Парсим фильтры
			categoryIDStrs := c.QueryArray("category_id")
			typeIDStrs := c.QueryArray("type_id")
			areaIDStrs := c.QueryArray("area_id")
			tagIDStrs := c.QueryArray("tag_id")

			var categoryIDs, typeIDs, areaIDs, tagIDs []uint

			for _, str := range categoryIDStrs {
				if id, err := strconv.ParseUint(str, 10, 32); err == nil {
					categoryIDs = append(categoryIDs, uint(id))
				}
			}
			for _, str := range typeIDStrs {
				if id, err := strconv.ParseUint(str, 10, 32); err == nil {
					typeIDs = append(typeIDs, uint(id))
				}
			}
			for _, str := range areaIDStrs {
				if id, err := strconv.ParseUint(str, 10, 32); err == nil {
					areaIDs = append(areaIDs, uint(id))
				}
			}
			for _, str := range tagIDStrs {
				if id, err := strconv.ParseUint(str, 10, 32); err == nil {
					tagIDs = append(tagIDs, uint(id))
				}
			}

		// Парсим параметры пагинации
		page := 1
		if pageStr := c.Query("page"); pageStr != "" {
			if parsed, err := strconv.Atoi(pageStr); err == nil && parsed > 0 {
				page = parsed
			}
		}

		limit := 20 // По умолчанию
		if limitStr := c.Query("limit"); limitStr != "" {
			if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		pagination := models.NewPaginationParams(page, limit)

		// Проверяем, нужен ли легкий формат (DTO)
		useLightDTO := c.Query("light") == "true"

		// Получаем данные с пагинацией
		places, total, err := placeService.List(categoryIDs, typeIDs, areaIDs, tagIDs, pagination)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// Загружаем изображения для каждого места (из images + files Strapi)
			for i := range places {
				placeService.LoadImagesForPlace(&places[i])
			}

			// Формируем ответ
			var data interface{}
			if useLightDTO {
				// Преобразуем в легкие DTO
				items := make([]models.PlaceListItem, len(places))
				for i, place := range places {
					items[i] = models.PlaceListItem{
						ID:          place.ID,
						Name:        place.Name,
						Type:        place.Type,
						Description: place.Description,
						Address:     place.Address,
						Latitude:    place.Latitude,
						Longitude:   place.Longitude,
						Rating:      place.Rating,
						TypeID:      place.TypeID,
						AreaID:      place.AreaID,
						CreatedAt:   place.CreatedAt,
						UpdatedAt:   place.UpdatedAt,
					}
					
					// Добавляем первое изображение если есть
					if len(place.Images) > 0 {
						for _, img := range place.Images {
							if img.IsActive {
								firstImg := img.URL
								items[i].FirstImage = &firstImg
								break
							}
						}
					}
				}
				data = items
		} else {
			// Полный формат (для обратной совместимости)
			data = places
		}

		// Формируем paginated ответ с новым форматом
		response := models.NewPaginatedResponse(data, total, pagination)

			c.JSON(200, response)
		})

		// Получение точки по ID
		placesGroup.GET("/:id", func(c *gin.Context) {
			id, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "недействительный ID"})
				return
			}

			place, err := placeService.GetByID(uint(id))
			if err != nil {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, place)
		})

		// Создание точки (только админ)
		placesGroup.POST("", auth.AuthMiddleware(cfg), auth.AdminMiddleware(), func(c *gin.Context) {
			var req PlaceRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			place := &models.Place{
				Name:          req.Name,
				Type:          req.Type,
				Description:   req.Description,
				Overview:      req.Overview,
				History:       req.History,
				Address:       req.Address,
				Hours:         req.Hours,
				Weekend:       req.Weekend,
				Entry:         req.Entry,
				Contacts:      req.Contacts,
				ContactsEmail: req.ContactsEmail,
				Latitude:      req.Latitude,
				Longitude:     req.Longitude,
				Rating:        req.Rating,

				// Старые поля
				OpeningHours: req.OpeningHours,
				TypeID:       req.TypeID,
				AreaID:       req.AreaID,
				IsActive:     true,
			}

			if err := placeService.Create(place); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// Привязываем категории
			if len(req.CategoryIDs) > 0 {
				var categories []models.Category
				if err := db.Where("id IN ?", req.CategoryIDs).Find(&categories).Error; err == nil {
					db.Model(place).Association("Categories").Append(categories)
				}
			}

			c.JSON(201, place)
		})

		// Обновление точки (только админ)
		placesGroup.PUT("/:id", auth.AuthMiddleware(cfg), auth.AdminMiddleware(), func(c *gin.Context) {
			id, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "недействительный ID"})
				return
			}

			var req PlaceRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			place := &models.Place{
				Name:          req.Name,
				Type:          req.Type,
				Description:   req.Description,
				Overview:      req.Overview,
				History:       req.History,
				Address:       req.Address,
				Hours:         req.Hours,
				Weekend:       req.Weekend,
				Entry:         req.Entry,
				Contacts:      req.Contacts,
				ContactsEmail: req.ContactsEmail,
				Latitude:      req.Latitude,
				Longitude:     req.Longitude,
				Rating:        req.Rating,

				// Старые поля
				OpeningHours: req.OpeningHours,
				TypeID:       req.TypeID,
				AreaID:       req.AreaID,
				IsActive:     true,
			}

			if err := placeService.Update(uint(id), place); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			// Привязываем категории
			if len(req.CategoryIDs) > 0 {
				var categories []models.Category
				if err := db.Where("id IN ?", req.CategoryIDs).Find(&categories).Error; err == nil {
					db.Model(place).Association("Categories").Clear()
					db.Model(place).Association("Categories").Append(categories)
				}
			}

			c.JSON(200, place)
		})

		// Удаление точки (только админ)
		placesGroup.DELETE("/:id", auth.AuthMiddleware(cfg), auth.AdminMiddleware(), func(c *gin.Context) {
			id, err := strconv.ParseUint(c.Param("id"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "недействительный ID"})
				return
			}

			if err := placeService.Delete(uint(id)); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "Точка удалена"})
		})
	}
}
