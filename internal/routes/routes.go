package routes

import (
	"net/http"
	"strconv"
	"tropa-nartov-backend/internal/auth"
	"tropa-nartov-backend/internal/config"
	"tropa-nartov-backend/internal/models"
	"tropa-nartov-backend/internal/route"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouteRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	routeService := route.NewService(db)
	routeGroup := router.Group("/routes")
	{
		// GET /routes - получить все маршруты с пагинацией
		routeGroup.GET("", func(c *gin.Context) {
			// Парсим фильтры
			typeIDStrs := c.QueryArray("type_id")
			areaIDStrs := c.QueryArray("area_id")
			tagIDStrs := c.QueryArray("tag_id")

			var typeIDs, areaIDs, tagIDs []uint

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
			limit := 20 // По умолчанию
			if limitStr := c.Query("limit"); limitStr != "" {
				if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
					if parsed > 100 {
						limit = 100 // Максимум 100 элементов
					} else {
						limit = parsed
					}
				}
			}

			offset := 0
			if offsetStr := c.Query("offset"); offsetStr != "" {
				if parsed, err := strconv.Atoi(offsetStr); err == nil && parsed >= 0 {
					offset = parsed
				}
			}

			// Проверяем, нужен ли легкий формат (DTO)
			useLightDTO := c.Query("light") == "true"

			// Получаем данные с пагинацией
			routes, total, err := routeService.List(typeIDs, areaIDs, tagIDs, limit, offset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Если маршрутов нет, возвращаем пустой ответ с пагинацией
			if len(routes) == 0 {
				response := models.PaginatedResponse{
					Data:       []interface{}{},
					Total:      0,
					Limit:      limit,
					Offset:     offset,
					HasMore:    false,
					NextOffset: nil,
				}
				c.JSON(http.StatusOK, response)
				return
			}

			// Загружаем связанные данные (Type, Area) для всех маршрутов
			for i := range routes {
				db.Preload("Type").Preload("Area").First(&routes[i], routes[i].ID)
			}

			// Формируем ответ
			var data interface{}
			if useLightDTO {
				// Преобразуем в легкие DTO
				items := make([]models.RouteListItem, len(routes))
				for i, route := range routes {
					items[i] = models.RouteListItem{
						ID:          route.ID,
						Name:        route.Name,
						Description: route.Description,
						Distance:    route.Distance,
						Duration:    route.Duration,
						TypeID:      route.TypeID,
						AreaID:      route.AreaID,
						Rating:      route.Rating,
						CreatedAt:   route.CreatedAt,
						UpdatedAt:   route.UpdatedAt,
					}
					
					// Добавляем имена типа и района
					if route.Type.ID != 0 {
						items[i].TypeName = route.Type.Name
					} else {
						items[i].TypeName = "Пеший поход"
					}
					
					if route.Area.ID != 0 {
						items[i].AreaName = route.Area.Name
					} else {
						items[i].AreaName = "Приэльбрусье"
					}
				}
				data = items
			} else {
				// Полный формат (для обратной совместимости)
			var response []gin.H
			for _, route := range routes {
				routeData := gin.H{
					"id":          route.ID,
					"name":        route.Name,
					"description": route.Description,
					"overview":    route.Overview,
					"history":     route.History,
					"distance":    route.Distance,
					"duration":    route.Duration,
					"type_id":     route.TypeID,
					"area_id":     route.AreaID,
					"rating":      route.Rating,
					"is_active":   route.IsActive,
					"created_at":  route.CreatedAt,
					"updated_at":  route.UpdatedAt,
				}

				if route.Type.ID != 0 {
					routeData["type_name"] = route.Type.Name
				} else {
						routeData["type_name"] = "Пеший поход"
				}

				if route.Area.ID != 0 {
					routeData["area_name"] = route.Area.Name
				} else {
						routeData["area_name"] = "Приэльбрусье"
				}

				response = append(response, routeData)
				}
				data = response
			}

			// Формируем ответ с пагинацией
			hasMore := offset+limit < int(total)
			var nextOffset *int
			if hasMore {
				next := offset + limit
				nextOffset = &next
			}

			response := models.PaginatedResponse{
				Data:       data,
				Total:      total,
				Limit:      limit,
				Offset:     offset,
				HasMore:    hasMore,
				NextOffset: nextOffset,
			}

			c.JSON(http.StatusOK, response)
		})

		// ДОБАВЛЕНО: Отладочный эндпоинт для проверки всех маршрутов
		routeGroup.GET("/debug/all", func(c *gin.Context) {
			var routes []models.Route
			if err := db.Find(&routes).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка при получении маршрутов",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"total_count": len(routes),
				"routes":      routes,
			})
		})

		// GET /routes/:id - получить маршрут по ID с местами
		routeGroup.GET("/:id", func(c *gin.Context) {
			idStr := c.Param("id")
			id, err := strconv.ParseUint(idStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Неверный ID маршрута",
				})
				return
			}

			// Используем service для получения маршрута с остановками
			route, err := routeService.GetByID(uint(id))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Маршрут не найден",
				})
				return
			}

			// Загружаем остановки (route_stops) с точками, упорядоченные по order_num
			var stops []models.RouteStop
			if err := db.Preload("Place").Preload("Place.Images").Where("route_id = ?", route.ID).Order("order_num ASC").Find(&stops).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка загрузки остановок",
				})
				return
			}

			// Загружаем связанные данные (Type, Area, Categories)
			if err := db.Preload("Type").Preload("Area").Preload("Categories").First(route, route.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка загрузки связанных данных",
				})
				return
			}

			// Формируем полный ответ
			response := gin.H{
				"id":          route.ID,
				"name":        route.Name,
				"description": route.Description,
				"overview":    route.Overview,
				"history":     route.History,
				"distance":    route.Distance,
				"duration":    route.Duration,
				"type_id":     route.TypeID,
				"area_id":     route.AreaID,
				"rating":      route.Rating,
				"is_active":   route.IsActive,
				"created_at":  route.CreatedAt,
				"updated_at":  route.UpdatedAt,
			}

			// Добавляем связанные данные
			if route.Type.ID != 0 {
				response["type"] = gin.H{
					"id":   route.Type.ID,
					"name": route.Type.Name,
				}
				response["type_name"] = route.Type.Name
			} else {
				response["type_name"] = "Пеший поход"
			}

			if route.Area.ID != 0 {
				response["area"] = gin.H{
					"id":   route.Area.ID,
					"name": route.Area.Name,
				}
				response["area_name"] = route.Area.Name
			} else {
				response["area_name"] = "Приэльбрусье"
			}

			// Добавляем категории
			if len(route.Categories) > 0 {
				var categories []gin.H
				for _, category := range route.Categories {
					categories = append(categories, gin.H{
						"id":   category.ID,
						"name": category.Name,
					})
				}
				response["categories"] = categories
			}

			// Добавляем места маршрута (stops)
			var stopsResponse []gin.H
			for _, stop := range stops {
				placeData := gin.H{
					"id":          stop.Place.ID,
					"name":        stop.Place.Name,
					"type":        stop.Place.Type,
					"description": stop.Place.Description,
					"address":     stop.Place.Address,
					"latitude":    stop.Place.Latitude,
					"longitude":   stop.Place.Longitude,
					"rating":      stop.Place.Rating,
				}

				// Добавляем изображения места
				if len(stop.Place.Images) > 0 {
					var images []gin.H
					for _, img := range stop.Place.Images {
						if img.IsActive {
							images = append(images, gin.H{
								"id":  img.ID,
								"url": img.URL,
							})
						}
					}
					placeData["images"] = images
				}

				stopsResponse = append(stopsResponse, gin.H{
					"place_id":  stop.PlaceID,
					"order_num": stop.OrderNum,
					"place":     placeData,
				})
			}
			response["stops"] = stopsResponse

			c.JSON(http.StatusOK, response)
		})

		// POST /routes - создать маршрут (требует авторизацию)
		routeGroup.POST("", auth.AuthMiddleware(cfg), func(c *gin.Context) {
			var input struct {
				Name        string  `json:"name" binding:"required"`
				Description string  `json:"description" binding:"required"`
				Overview    string  `json:"overview"`
				History     string  `json:"history"`
				Distance    float64 `json:"distance" binding:"required"`
				Duration    float64 `json:"duration"`
				TypeID      uint    `json:"type_id" binding:"required"`
				AreaID      uint    `json:"area_id" binding:"required"`
				CategoryIDs []uint  `json:"category_ids"`
				Stops       []struct {
					PlaceID  uint `json:"place_id" binding:"required"`
					OrderNum int  `json:"order_num" binding:"required"`
				} `json:"stops" binding:"required"`
			}

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Неверные данные: " + err.Error(),
				})
				return
			}

			// Проверяем, что места существуют
			var placeIDs []uint
			for _, stop := range input.Stops {
				placeIDs = append(placeIDs, stop.PlaceID)
			}

			var placeCount int64
			if err := db.Model(&models.Place{}).Where("id IN ?", placeIDs).Count(&placeCount).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка при проверке мест",
				})
				return
			}

			if placeCount != int64(len(placeIDs)) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Одно или несколько указанных мест не существуют",
				})
				return
			}

			// Создаем маршрут
			route := &models.Route{
				Name:        input.Name,
				Description: input.Description,
				Overview:    input.Overview,
				History:     input.History,
				Distance:    input.Distance,
				Duration:    input.Duration,
				TypeID:      input.TypeID,
				AreaID:      input.AreaID,
				IsActive:    true,
			}

			if err := routeService.Create(route); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка при создании маршрута: " + err.Error(),
				})
				return
			}

			// Добавляем категории если указаны
			if len(input.CategoryIDs) > 0 {
				var categories []models.Category
				if err := db.Where("id IN ?", input.CategoryIDs).Find(&categories).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Ошибка при загрузке категорий",
					})
					return
				}
				if err := db.Model(route).Association("Categories").Append(categories); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Ошибка при добавлении категорий",
					})
					return
				}
			}

			// Создаем места маршрута (RouteStop)
			var stops []models.RouteStop
			for _, stop := range input.Stops {
				stops = append(stops, models.RouteStop{
					RouteID:  route.ID,
					PlaceID:  stop.PlaceID,
					OrderNum: stop.OrderNum,
				})
			}

			if len(stops) > 0 {
				if err := db.Create(&stops).Error; err != nil {
					// Откатываем создание маршрута при ошибке
					db.Delete(route)
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Ошибка при создании мест маршрута: " + err.Error(),
					})
					return
				}
			}

			// Загружаем созданный маршрут с связанными данными
			if err := db.Preload("Type").Preload("Area").Preload("Categories").First(route, route.ID).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка при загрузке созданного маршрута",
				})
				return
			}

			c.JSON(http.StatusCreated, gin.H{
				"message": "Маршрут успешно создан",
				"id":      route.ID,
				"route":   route,
			})
		})

		// PUT /routes/:id - обновить маршрут (требует авторизацию)
		routeGroup.PUT("/:id", func(c *gin.Context) {

			id := c.Param("id")

			var route models.Route
			if err := db.First(&route, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Маршрут не найден",
				})
				return
			}

			var input struct {
				Name        string  `json:"name"`
				Description string  `json:"description"`
				Overview    string  `json:"overview"`
				History     string  `json:"history"`
				Distance    float64 `json:"distance"`
				Duration    float64 `json:"duration"`
				TypeID      uint    `json:"type_id"`
				AreaID      uint    `json:"area_id"`
				IsActive    bool    `json:"is_active"`
				CategoryIDs []uint  `json:"category_ids"`
			}

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Неверные данные: " + err.Error(),
				})
				return
			}

			// Обновляем только переданные поля
			updates := make(map[string]interface{})
			if input.Name != "" {
				updates["name"] = input.Name
			}
			if input.Description != "" {
				updates["description"] = input.Description
			}
			if input.Overview != "" {
				updates["overview"] = input.Overview
			}
			if input.History != "" {
				updates["history"] = input.History
			}
			if input.Distance != 0 {
				updates["distance"] = input.Distance
			}
			if input.Duration != 0 {
				updates["duration"] = input.Duration
			}
			if input.TypeID != 0 {
				updates["type_id"] = input.TypeID
			}
			if input.AreaID != 0 {
				updates["area_id"] = input.AreaID
			}
			updates["is_active"] = input.IsActive

			if err := db.Model(&route).Updates(updates).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка при обновлении маршрута",
				})
				return
			}

			// Обновляем категории если переданы
			if input.CategoryIDs != nil {
				if err := db.Model(&route).Association("Categories").Replace(input.CategoryIDs); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Ошибка при обновлении категорий",
					})
					return
				}
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Маршрут успешно обновлен",
				"route":   route,
			})
		})

		// DELETE /routes/:id - удалить маршрут (требует авторизацию)
		routeGroup.DELETE("/:id", func(c *gin.Context) {

			id := c.Param("id")

			var route models.Route
			if err := db.First(&route, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Маршрут не найден",
				})
				return
			}

			if err := db.Delete(&route).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Ошибка при удалении маршрута",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Маршрут успешно удален",
			})
		})
	}
}
