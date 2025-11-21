package routes

import (
	"errors"
	"strconv"
	"tropa-nartov-backend/internal/auth"
	"tropa-nartov-backend/internal/config"
	"tropa-nartov-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupFavoriteRoutes настраивает маршруты для избранного
func SetupFavoriteRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
	favoritesGroup := r.Group("/favorites")
	favoritesGroup.Use(auth.AuthMiddleware(cfg))
	{
		// Получить избранные места пользователя
		favoritesGroup.GET("/places", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			var favoritePlaces []models.FavoritePlace
			if err := db.Preload("Place").Preload("Place.Images").
				Where("user_id = ?", userID).
				Find(&favoritePlaces).Error; err != nil {
				c.JSON(500, gin.H{"error": "Ошибка получения избранных мест"})
				return
			}

			// Преобразуем в список мест
			places := make([]models.Place, len(favoritePlaces))
			for i, fp := range favoritePlaces {
				places[i] = fp.Place
			}

			c.JSON(200, places)
		})

		// Добавить место в избранное
		favoritesGroup.POST("/places/:placeId", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			placeID, err := strconv.ParseUint(c.Param("placeId"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Неверный ID места"})
				return
			}

			// Проверяем, существует ли место
			var place models.Place
			if err := db.Where("id = ? AND is_active = ?", placeID, true).First(&place).Error; err != nil {
				c.JSON(404, gin.H{"error": "Место не найдено"})
				return
			}

			// Проверяем, не добавлено ли уже в избранное
			var existingFavorite models.FavoritePlace
			// Используем Take() вместо First() - он не логирует ошибку если запись не найдена
			if err := db.Where("user_id = ? AND place_id = ?", userID, placeID).Take(&existingFavorite).Error; err == nil {
				c.JSON(400, gin.H{"error": "Место уже в избранном"})
				return
			}

			// Добавляем в избранное
			favorite := models.FavoritePlace{
				UserID:  userID.(uint),
				PlaceID: uint(placeID),
			}

			if err := db.Create(&favorite).Error; err != nil {
				c.JSON(500, gin.H{"error": "Ошибка добавления в избранное"})
				return
			}

			c.JSON(200, gin.H{"message": "Место добавлено в избранное"})
		})

		// Удалить место из избранного
		favoritesGroup.DELETE("/places/:placeId", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			placeID, err := strconv.ParseUint(c.Param("placeId"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Неверный ID места"})
				return
			}

			// Удаляем из избранного
			if err := db.Where("user_id = ? AND place_id = ?", userID, placeID).
				Delete(&models.FavoritePlace{}).Error; err != nil {
				c.JSON(500, gin.H{"error": "Ошибка удаления из избранного"})
				return
			}

			c.JSON(200, gin.H{"message": "Место удалено из избранного"})
		})

		// Проверить, находится ли место в избранном
		favoritesGroup.GET("/places/:placeId/status", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			placeID, err := strconv.ParseUint(c.Param("placeId"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Неверный ID места"})
				return
			}

			var favorite models.FavoritePlace
			err = db.Where("user_id = ? AND place_id = ?", userID, placeID).First(&favorite).Error
			if err != nil {
				// Если запись не найдена, это нормально - возвращаем false
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(200, gin.H{"is_favorite": false})
					return
				}
				// Другие ошибки - возвращаем 500
				c.JSON(500, gin.H{"error": "Ошибка проверки статуса избранного"})
				return
			}

			c.JSON(200, gin.H{"is_favorite": true})
		})

		// === МАРШРУТЫ ===

		// Получить избранные маршруты пользователя
		favoritesGroup.GET("/routes", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			var favoriteRoutes []models.FavoriteRoute
			if err := db.Preload("Route").Preload("Route.Type").Preload("Route.Area").
				Where("user_id = ?", userID).
				Find(&favoriteRoutes).Error; err != nil {
				c.JSON(500, gin.H{"error": "Ошибка получения избранных маршрутов"})
				return
			}

			// Преобразуем в JSON ответ с полными данными
			var response []gin.H
			for _, fr := range favoriteRoutes {
				route := fr.Route
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

				// Добавляем данные о типе если есть
				if route.Type.ID != 0 {
					routeData["type_name"] = route.Type.Name
				} else {
					routeData["type_name"] = "Пеший поход" // значение по умолчанию
				}

				// Добавляем данные о районе если есть
				if route.Area.ID != 0 {
					routeData["area_name"] = route.Area.Name
				} else {
					routeData["area_name"] = "Приэльбрусье" // значение по умолчанию
				}

				response = append(response, routeData)
			}

			c.JSON(200, response)
		})

		// Добавить маршрут в избранное
		favoritesGroup.POST("/routes/:routeId", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			routeID, err := strconv.ParseUint(c.Param("routeId"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Неверный ID маршрута"})
				return
			}

			// Проверяем, существует ли маршрут
			var route models.Route
			if err := db.Where("id = ? AND is_active = ?", routeID, true).First(&route).Error; err != nil {
				c.JSON(404, gin.H{"error": "Маршрут не найден"})
				return
			}

			// Проверяем, не добавлен ли уже в избранное
			var existingFavorite models.FavoriteRoute
			// Используем Take() вместо First() - он не логирует ошибку если запись не найдена
			err = db.Where("user_id = ? AND route_id = ?", userID, routeID).Take(&existingFavorite).Error
			if err == nil {
				c.JSON(400, gin.H{"error": "Маршрут уже в избранном"})
				return
			}
			// Если ошибка - это нормально (запись не найдена), продолжаем

			// Добавляем в избранное
			favorite := models.FavoriteRoute{
				UserID:  userID.(uint),
				RouteID: uint(routeID),
			}

			if err := db.Create(&favorite).Error; err != nil {
				c.JSON(500, gin.H{"error": "Ошибка добавления в избранное"})
				return
			}

			c.JSON(200, gin.H{"message": "Маршрут добавлен в избранное"})
		})

		// Удалить маршрут из избранного
		favoritesGroup.DELETE("/routes/:routeId", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			routeID, err := strconv.ParseUint(c.Param("routeId"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Неверный ID маршрута"})
				return
			}

			// Удаляем из избранного
			if err := db.Where("user_id = ? AND route_id = ?", userID, routeID).
				Delete(&models.FavoriteRoute{}).Error; err != nil {
				c.JSON(500, gin.H{"error": "Ошибка удаления из избранного"})
				return
			}

			c.JSON(200, gin.H{"message": "Маршрут удален из избранного"})
		})

		// Проверить, находится ли маршрут в избранном
		favoritesGroup.GET("/routes/:routeId/status", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			routeID, err := strconv.ParseUint(c.Param("routeId"), 10, 32)
			if err != nil {
				c.JSON(400, gin.H{"error": "Неверный ID маршрута"})
				return
			}

			var favorite models.FavoriteRoute
			err = db.Where("user_id = ? AND route_id = ?", userID, routeID).First(&favorite).Error
			if err != nil {
				// Если запись не найдена, это нормально - возвращаем false
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(200, gin.H{"is_favorite": false})
					return
				}
				// Другие ошибки - возвращаем 500
				c.JSON(500, gin.H{"error": "Ошибка проверки статуса избранного"})
				return
			}

			c.JSON(200, gin.H{"is_favorite": true})
		})

		// Массовая проверка статусов маршрутов в избранном
		favoritesGroup.POST("/routes/statuses", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(401, gin.H{"error": "Не авторизован"})
				return
			}

			var request struct {
				RouteIDs []uint `json:"route_ids"`
			}
			
			if err := c.ShouldBindJSON(&request); err != nil {
				c.JSON(400, gin.H{"error": "Неверный формат данных"})
				return
			}

			// Если список пустой, возвращаем пустой объект
			if len(request.RouteIDs) == 0 {
				c.JSON(200, gin.H{})
				return
			}

			// Получаем все избранные маршруты пользователя из указанного списка
			var favorites []models.FavoriteRoute
			if err := db.Where("user_id = ? AND route_id IN ?", userID, request.RouteIDs).
				Find(&favorites).Error; err != nil {
				c.JSON(500, gin.H{"error": "Ошибка проверки статусов избранного"})
				return
			}

			// Создаем map с результатами
			result := make(map[string]bool)
			favoriteMap := make(map[uint]bool)
			
			for _, fav := range favorites {
				favoriteMap[fav.RouteID] = true
			}
			
			for _, routeID := range request.RouteIDs {
				result[strconv.FormatUint(uint64(routeID), 10)] = favoriteMap[routeID]
			}

			c.JSON(200, result)
		})
	}
}
