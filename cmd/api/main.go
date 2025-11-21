package main

import (
	"strings"
	"time"
	"tropa-nartov-backend/internal/auth"
	"tropa-nartov-backend/internal/config"
	"tropa-nartov-backend/internal/db"
	"tropa-nartov-backend/internal/logger"
	"tropa-nartov-backend/internal/middleware"
	"tropa-nartov-backend/internal/routes"
	_ "tropa-nartov-backend/docs" // swagger docs

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Tropa Nartov API
// @version 1.0
// @description REST API для мобильного приложения "Тропа Нартов"
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email tropanartov@yandex.ru

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8001
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer {token}

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		// .env файл опционален в некоторых окружениях
	}

	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Инициализируем logger
	if err := logger.Init(cfg.Environment, cfg.Debug); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	logger.Info("Starting Tropa Nartov API",
		zap.String("environment", cfg.Environment),
		zap.String("version", "1.0.0"),
		zap.String("port", cfg.Port),
	)

	// Подключаемся к БД
	dbConnection, err := db.Connect(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("Database connected successfully")

	// Выполняем миграции
	if err := db.Migrate(dbConnection); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}
	logger.Info("Database migrations completed")

	// Создаем тестовые данные
	// if err := createTestData(dbConnection); err != nil {
	// log.Printf("Warning: Failed to create test data: %v", err)
	// }

	// Инициализируем Gin
	r := gin.Default()

	// Добавляем ETag middleware для кеширования
	r.Use(auth.ETagMiddleware())

	// Добавляем глобальный rate limiting (100 запросов в минуту на IP)
	r.Use(middleware.GlobalRateLimit())
	logger.Info("Rate limiting enabled", zap.Int("limit", 100), zap.String("period", "1 minute"))

	// Настройка CORS из конфигурации
	var allowedOrigins []string
	if cfg.CORSAllowedOrigins != "" {
		// Парсим строку через запятую
		origins := strings.Split(cfg.CORSAllowedOrigins, ",")
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}
	// Проверка: CORS origins должны быть установлены
	if len(allowedOrigins) == 0 {
		// В development режиме используем безопасные defaults
		if cfg.Environment == "development" || cfg.Debug {
			allowedOrigins = []string{
				"http://localhost:3000",
				"http://localhost:8080",
				"http://localhost:8001",
			}
			logger.Warn("CORS_ALLOWED_ORIGINS not set, using development defaults",
				zap.Strings("origins", allowedOrigins),
			)
		} else {
			logger.Fatal("CORS_ALLOWED_ORIGINS must be set in production environment")
		}
	}
	logger.Info("CORS configured", zap.Strings("allowed_origins", allowedOrigins))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "If-None-Match"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		ExposeHeaders:    []string{"ETag", "Content-Length"},
	}))

	// Настраиваем маршруты
	logger.Info("Setting up routes...")
	routes.SetupAuthRoutes(r, dbConnection, cfg)
	routes.SetupPlaceRoutes(r, dbConnection, cfg)
	routes.SetupRouteRoutes(r, dbConnection, cfg)
	routes.SetupReviewRoutes(r, dbConnection, cfg)
	routes.SetupFavoriteRoutes(r, dbConnection, cfg)
	routes.SetupActivityRoutes(r, dbConnection, cfg)
	logger.Info("Routes configured successfully")

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Info("Swagger documentation available at /swagger/index.html")

	// Тестовый эндпоинт
	// @Summary Health check
	// @Description Проверка работоспособности API
	// @Tags health
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /ping [get]
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"timestamp": time.Now().Unix(),
		})
	})

	// Запускаем сервер
	logger.Info("Server starting", 
		zap.String("port", cfg.Port),
		zap.String("host", cfg.Host),
	)
	logger.Info("Available endpoints:",
		zap.Strings("auth", []string{
			"POST /auth/register",
			"POST /auth/login",
			"POST /auth/forgot-password",
			"POST /auth/verify-reset-code",
			"POST /auth/reset-password",
			"DELETE /auth/delete-account (protected)",
			"POST /auth/upload-avatar (protected)",
			"DELETE /auth/delete-avatar (protected)",
		}),
		zap.Strings("public", []string{
			"GET /ping",
			"GET /places",
			"GET /routes",
			"GET /reviews/place/:placeId",
		}),
		zap.Strings("protected", []string{
			"POST /reviews",
			"GET /favorites/places",
			"POST /favorites/places/:placeId",
			"DELETE /favorites/places/:placeId",
			"GET /favorites/places/:placeId/status",
		}),
	)
	// log.Printf("   GET    /favorites/routes (protected)")
	// log.Printf("   POST   /favorites/routes/:routeId (protected)")
	// log.Printf("   DELETE /favorites/routes/:routeId (protected)")
	// log.Printf("   GET    /favorites/routes/:routeId/status (protected)")
	// log.Printf("   GET    /user/activity (protected)")
	// log.Printf("   GET    /user/activity/places (protected)")
	// log.Printf("   GET    /user/activity/routes (protected)")
	// log.Printf("   POST   /user/activity/places/:placeId (protected)")
	// log.Printf("   POST   /user/activity/routes/:routeId (protected)")
	// log.Printf("   DELETE /user/activity/places/:placeId (protected)")
	// log.Printf("   DELETE /user/activity/routes/:routeId (protected)")

	if err := r.Run(":" + cfg.Port); err != nil {
		// log.Fatalf("Failed to run server: %v", err)
	}
}

// createTestData создает тестовые данные для разработки
// func createTestData(db *gorm.DB) error {
// 	// Временно отключаем foreign key constraints
// 	db.Exec("SET session_replication_role = 'replica'")

// 	// Очищаем данные в правильном порядке (сначала зависимые таблицы)
// 	db.Exec("DELETE FROM favorite_places")
// 	db.Exec("DELETE FROM reviews")
// 	db.Exec("DELETE FROM images")
// 	db.Exec("DELETE FROM place_categories")
// 	db.Exec("DELETE FROM places")
// 	db.Exec("DELETE FROM categories")
// 	db.Exec("DELETE FROM areas")
// 	db.Exec("DELETE FROM types")

// 	// Включаем обратно constraints
// 	db.Exec("SET session_replication_role = 'origin'")

// 	// Создаем тестовые типы маршрутов (с проверкой на существование)
// 	types := []models.Type{
// 		{Name: "Пеший поход", Description: "Пешие маршруты по горам и лесам"},
// 		{Name: "Веломаршрут", Description: "Велосипедные трассы"},
// 		{Name: "Автомобильный", Description: "Маршруты для автомобильных поездок"},
// 		{Name: "Комбинированный", Description: "Маршруты с разными видами передвижения"},
// 	}

// 	for i := range types {
// 		if err := db.Create(&types[i]).Error; err != nil {
// 			return err
// 		}
// 	}

// 	// Создаем тестовые районы
// 	areas := []models.Area{
// 		{Name: "Приэльбрусье", Description: "Район вокруг горы Эльбрус"},
// 		{Name: "Домбай", Description: "Домбайский район"},
// 		{Name: "Архыз", Description: "Архызский район"},
// 		{Name: "Чегемское ущелье", Description: "Чегемское ущелье"},
// 		{Name: "Верхняя Балкария", Description: "Верхнебалкарский район"},
// 		{Name: "Карачаевск", Description: "Карачаевский район"},
// 	}

// 	for i := range areas {
// 		if err := db.Create(&areas[i]).Error; err != nil {
// 			return err
// 		}
// 	}

// 	// СОЗДАЕМ ТЕСТОВЫЕ МЕСТА (городские точки в Нальчике)
// 	testPlaces := []models.Place{
// 		{
// 			Name:          `Кафе "Горный Ветерок"`,
// 			Type:          "кафе",
// 			Description:   "Уютное кафе с кавказской кухней и прекрасным видом",
// 			Overview:      "Традиционные блюда кавказской кухни, свежая выпечка, кофе по-восточному",
// 			History:       "Основано в 1995 году, с тех пор является одним из любимых мест жителей и гостей города",
// 			Address:       "г. Нальчик, пр. Ленина, 25",
// 			Hours:         "08:00-23:00",
// 			Weekend:       "Без выходных",
// 			Entry:         "Свободный",
// 			Contacts:      "+78662223344",
// 			ContactsEmail: "gorniy-veterok@mail.ru",
// 			Latitude:      43.4981,
// 			Longitude:     43.7189,
// 			Rating:        4.5,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          `Парк "Атажукинский"`,
// 			Type:          "парк",
// 			Description:   "Крупнейший парк на Северном Кавказе с богатой растительностью",
// 			Overview:      "Площадь 250 гектаров, аллеи, фонтаны, озера, зоопарк и развлечения",
// 			History:       "Заложен в 1847 году как Казенный сад, современное название получил в 1920 году",
// 			Address:       "г. Нальчик, ул. Толстого, 2",
// 			Hours:         "06:00-24:00",
// 			Weekend:       "Без выходных",
// 			Entry:         "Свободный",
// 			Contacts:      "+78662225566",
// 			ContactsEmail: "park@nalchik.ru",
// 			Latitude:      43.4925,
// 			Longitude:     43.6123,
// 			Rating:        4.8,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          "Музей Кабардино-Балкарии",
// 			Type:          "музей",
// 			Description:   "Национальный музей с богатой коллекцией истории и культуры народов КБР",
// 			Overview:      "Экспозиции по археологии, этнографии, искусству и современной истории",
// 			History:       "Основан в 1921 году, в фондах более 150 тысяч экспонатов",
// 			Address:       "г. Нальчик, ул. Горького, 62",
// 			Hours:         "10:00-18:00",
// 			Weekend:       "Понедельник",
// 			Entry:         "Платный",
// 			Contacts:      "+78662227788",
// 			ContactsEmail: "museum@kbr.ru",
// 			Latitude:      43.5050,
// 			Longitude:     43.6250,
// 			Rating:        4.2,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          `Ресторан "Эльбрус"`,
// 			Type:          "ресторан",
// 			Description:   "Ресторан национальной кухни с кавказским гостеприимством",
// 			Overview:      "Блюда кабардинской, балкарской и общекавказской кухни, живая музыка",
// 			History:       "Работает с 1985 года, известен своими шеф-поварами и качеством обслуживания",
// 			Address:       "г. Нальчик, ул. Кабардинская, 45",
// 			Hours:         "11:00-02:00",
// 			Weekend:       "Без выходных",
// 			Entry:         "Свободный",
// 			Contacts:      "+78662229900",
// 			ContactsEmail: "elbrus-rest@yandex.ru",
// 			Latitude:      43.4900,
// 			Longitude:     43.6300,
// 			Rating:        3.9,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          `Сквер "Дружбы"`,
// 			Type:          "парк",
// 			Description:   "Уютный сквер в центре города для прогулок и отдыха",
// 			Overview:      "Благоустроенная территория с фонтаном, скамейками и детской площадкой",
// 			History:       "Создан в 1970-х годах в честь дружбы народов Кавказа",
// 			Address:       "г. Нальчик, пл. Согласия",
// 			Hours:         "Круглосуточно",
// 			Weekend:       "Без выходных",
// 			Entry:         "Свободный",
// 			Contacts:      "",
// 			ContactsEmail: "",
// 			Latitude:      43.5100,
// 			Longitude:     43.6100,
// 			Rating:        4.0,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          `Кофейня "Утро в Горах"`,
// 			Type:          "кафе",
// 			Description:   "Кофейня с авторскими напитками и домашней атмосферой",
// 			Overview:      "Свежая выпечка, десерты, более 20 видов кофе и чая",
// 			History:       "Открыта в 2018 году, быстро стала популярным местом для встреч",
// 			Address:       "г. Нальчик, пр. Шогенцукова, 15",
// 			Hours:         "07:00-22:00",
// 			Weekend:       "Без выходных",
// 			Entry:         "Свободный",
// 			Contacts:      "+78662221122",
// 			ContactsEmail: "utro-v-gorah@coffee.ru",
// 			Latitude:      43.4850,
// 			Longitude:     43.6200,
// 			Rating:        4.7,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          `Театр "Нальчикский"`,
// 			Type:          "театр",
// 			Description:   "Государственный драматический театр им. А. Шогенцукова",
// 			Overview:      "Спектакли на кабардинском и русском языках, классические и современные постановки",
// 			History:       "Основан в 1936 году, носит имя основоположника кабардинской литературы",
// 			Address:       "г. Нальчик, ул. Пушкина, 1",
// 			Hours:         "10:00-19:00",
// 			Weekend:       "Понедельник",
// 			Entry:         "Платный",
// 			Contacts:      "+78662223355",
// 			ContactsEmail: "teatr@nalchik.ru",
// 			Latitude:      43.5000,
// 			Longitude:     43.6050,
// 			Rating:        4.3,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          `Магазин "Сувениры Кавказа"`,
// 			Type:          "магазин",
// 			Description:   "Магазин традиционных кавказских сувениров и подарков",
// 			Overview:      "Оружие, кинжалы, украшения, ковры, национальная одежда и утварь",
// 			History:       "Работает более 20 лет, известен качеством и аутентичностью товаров",
// 			Address:       "г. Нальчик, ул. Советская, 78",
// 			Hours:         "09:00-20:00",
// 			Weekend:       "Воскресенье",
// 			Entry:         "Свободный",
// 			Contacts:      "+78662224433",
// 			ContactsEmail: "souvenir@kavkaz.ru",
// 			Latitude:      43.4950,
// 			Longitude:     43.6400,
// 			Rating:        3.5,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          `Стадион "Спартак"`,
// 			Type:          "стадион",
// 			Description:   "Спортивный комплекс для футбола и легкой атлетики",
// 			Overview:      "Футбольное поле с трибунами, беговые дорожки, спортивные секции",
// 			History:       "Построен в 1960 году, реконструирован в 2010-м",
// 			Address:       "г. Нальчик, ул. Октябрьская, 56",
// 			Hours:         "07:00-22:00",
// 			Weekend:       "Без выходных",
// 			Entry:         "Платный",
// 			Contacts:      "+78662226677",
// 			ContactsEmail: "spartak@stadion.ru",
// 			Latitude:      43.5150,
// 			Longitude:     43.6150,
// 			Rating:        4.1,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 		{
// 			Name:          `Библиотека "Центральная"`,
// 			Type:          "библиотека",
// 			Description:   "Главная библиотека республики с богатым фондом литературы",
// 			Overview:      "Читальные залы, абонемент, электронный каталог, мероприятия",
// 			History:       "Основана в 1922 году, фонд составляет более 500 тысяч изданий. В сердце «Старого города» или «Эски шахар» находится одна из главных достопримечательностей Ташкента – огромный базар Чорсу, известный еще со времен Средневековья. Оказавшись на этом базаре, вы попадаете в восточную сказку. Тут вся история Узбекистана: керамические изделия, тюбетейки, национальные халаты, восточные сладости, специи, фрукты, овощи, глиняные изделия, сувениры ручной работы, книги, подарки, платки из национальных тканей и многое другое. Перечислять разнообразие товаров можно бесконечноВ сердце «Старого города» или «Эски шахар» находится одна из главных достопримечательностей Ташкента – огромный базар Чорсу, известный еще со времен Средневековья. Оказавшись на этом базаре, вы попадаете в восточную сказку. Тут вся история Узбекистана: керамические изделия, тюбетейки, национальные халаты, восточные сладости, специи, фрукты, овощи, глиняные изделия, сувениры ручной работы, книги, подарки, платки из национальных тканей и многое другое. Перечислять разнообразие товаров можно бесконечноВ сердце «Старого города» или «Эски шахар» находится одна из главных достопримечательностей Ташкента – огромный базар Чорсу, известный еще со времен Средневековья. Оказавшись на этом базаре, вы попадаете в восточную сказку. Тут вся история Узбекистана: керамические изделия, тюбетейки, национальные халаты, восточные сладости, специи, фрукты, овощи, глиняные изделия, сувениры ручной работы, книги, подарки, платки из национальных тканей и многое другое. Перечислять разнообразие товаров можно бесконечноВ сердце «Старого города» или «Эски шахар» находится одна из главных достопримечательностей Ташкента – огромный базар Чорсу, известный еще со времен Средневековья. Оказавшись на этом базаре, вы попадаете в восточную сказку. Тут вся история Узбекистана: керамические изделия, тюбетейки, национальные халаты, ",
// 			Address:       "г. Нальчик, ул. Ногмова, 42",
// 			Hours:         "09:00-19:00",
// 			Weekend:       "Воскресенье",
// 			Entry:         "Свободный",
// 			Contacts:      "+78662228899",
// 			ContactsEmail: "library@kbr.ru",
// 			Latitude:      43.4800,
// 			Longitude:     43.6250,
// 			Rating:        4.6,
// 			TypeID:        types[0].ID,
// 			AreaID:        areas[0].ID,
// 			IsActive:      true,
// 		},
// 	}

// 	// Создаем места и сохраняем их ID
// 	var createdPlaces []models.Place
// 	for i := range testPlaces {
// 		if err := db.Create(&testPlaces[i]).Error; err != nil {
// 			// log.Printf("Error creating place %d: %v", i, err)
// 			return err
// 		}
// 		createdPlaces = append(createdPlaces, testPlaces[i])
// 	}

// 	// Создаем тестовые изображения для ВСЕХ мест
// 	testImages := []models.Image{
// 		{URL: "https://picsum.photos/400/300?random=1", PlaceID: createdPlaces[0].ID},
// 		{URL: "https://picsum.photos/400/300?random=2", PlaceID: createdPlaces[1].ID},
// 		{URL: "https://picsum.photos/400/300?random=3", PlaceID: createdPlaces[2].ID},
// 		{URL: "https://picsum.photos/400/300?random=4", PlaceID: createdPlaces[3].ID},
// 		{URL: "https://picsum.photos/400/300?random=5", PlaceID: createdPlaces[4].ID},
// 		{URL: "https://picsum.photos/400/300?random=6", PlaceID: createdPlaces[5].ID},
// 		{URL: "https://picsum.photos/400/300?random=7", PlaceID: createdPlaces[6].ID},
// 		{URL: "https://picsum.photos/400/300?random=8", PlaceID: createdPlaces[7].ID},
// 		{URL: "https://picsum.photos/400/300?random=9", PlaceID: createdPlaces[8].ID},
// 		{URL: "https://picsum.photos/400/300?random=10", PlaceID: createdPlaces[9].ID},
// 	}

// 	for i := range testImages {
// 		if err := db.Create(&testImages[i]).Error; err != nil {
// 			// log.Printf("Error creating image %d: %v", i, err)
// 			return err
// 		}
// 	}

// 	// ИСПРАВЛЕНИЕ: Получаем существующих пользователей вместо создания новых
// 	var existingUsers []models.User
// 	if err := db.Find(&existingUsers).Error; err != nil {
// 		// log.Printf("Warning: Failed to get existing users: %v", err)
// 	}

// 	// Если пользователей нет, создаем одного тестового
// 	if len(existingUsers) == 0 {
// 		testUser := models.User{
// 			Name:         "Тестовый Пользователь",
// 			Email:        "test@example.com",
// 			PasswordHash: "12345678", // password
// 			Role:         "user",
// 			IsActive:     true,
// 		}
// 		if err := db.Create(&testUser).Error; err != nil {
// 			// log.Printf("Warning: Failed to create test user: %v", err)
// 		} else {
// 			existingUsers = append(existingUsers, testUser)
// 		}
// 	}

// 	// Используем существующих пользователей для отзывов
// 	var createdUsers []models.User
// 	if len(existingUsers) > 0 {
// 		createdUsers = existingUsers
// 	} else {
// 		// Если все же нет пользователей, создаем отзывы без пользователей
// 		// log.Printf("Warning: No users available for reviews")
// 	}

// 	// Создаем по 3 тестовых отзыва для каждого места
// 	var testReviews []models.Review
// 	reviewTexts := []string{
// 		"Отличное место! Очень рекомендую к посещению.",
// 		"Хорошая атмосфера, но цены немного завышены.",
// 		"Прекрасное обслуживание и качественная еда!",
// 		"Красивое место, идеально для фотосессий.",
// 		"Уютная атмосфера, приятная музыка.",
// 		"Быстрое обслуживание, вежливый персонал.",
// 		"Интересная архитектура, богатая история.",
// 		"Чисто и аккуратно, приятно проводить время.",
// 		"Разнообразное меню, есть из чего выбрать.",
// 		"Комфортные условия, удобное расположение.",
// 		"Отличное место для семейного отдыха.",
// 		"Интересные экспозиции, познавательно.",
// 		"Приятные цены, хорошее соотношение цены и качества.",
// 		"Красивый вид, романтическая атмосфера.",
// 		"Много зелени, свежий воздух.",
// 		"Современное оборудование, все в отличном состоянии.",
// 		"Богатый выбор литературы, уютные залы.",
// 		"Вкусный кофе и свежая выпечка.",
// 		"Профессиональные сотрудники, грамотные консультации.",
// 		"Отличная спортивная инфраструктура.",
// 		"Удобное расписание, доступные цены.",
// 		"Интересные мероприятия, разнообразная программа.",
// 		"Качественные товары, приятные продавцы.",
// 		"Чистые помещения, современный дизайн.",
// 		"Быстрый интернет, комфортные рабочие места.",
// 		"Просторные помещения, хорошая вентиляция.",
// 		"Насыщенная культурная программа.",
// 		"Аутентичная атмосфера, чувствуется местный колорит.",
// 		"Удобное расположение, легко добраться.",
// 		"Приветливый персонал, всегда готовы помочь.",
// 	}

// 	// ИСПРАВЛЕНИЕ: Проверяем, что есть пользователи перед созданием отзывов
// 	if len(createdUsers) > 0 {
// 		for placeIndex, place := range createdPlaces {
// 			for i := 0; i < 3; i++ {
// 				userIndex := (placeIndex*3 + i) % len(createdUsers)
// 				textIndex := (placeIndex*3 + i) % len(reviewTexts)
// 				rating := 3 + (placeIndex+i)%3 // Рейтинг от 3 до 5

// 				review := models.Review{
// 					UserID:   createdUsers[userIndex].ID,
// 					PlaceID:  &place.ID,
// 					Text:     reviewTexts[textIndex],
// 					Rating:   rating,
// 					IsActive: true,
// 				}
// 				testReviews = append(testReviews, review)
// 			}
// 		}
// 	} else {
// 		// Создаем отзывы без привязки к пользователям (если пользователей нет)
// 		for placeIndex, place := range createdPlaces {
// 			for i := 0; i < 3; i++ {
// 				textIndex := (placeIndex*3 + i) % len(reviewTexts)
// 				rating := 3 + (placeIndex+i)%3 // Рейтинг от 3 до 5

// 				review := models.Review{
// 					PlaceID:  &place.ID,
// 					Text:     reviewTexts[textIndex],
// 					Rating:   rating,
// 					IsActive: true,
// 				}
// 				testReviews = append(testReviews, review)
// 			}
// 		}
// 	}

// 	for i := range testReviews {
// 		if err := db.Create(&testReviews[i]).Error; err != nil {
// 			// log.Printf("Error creating review %d: %v", i, err)
// 			continue
// 		}
// 	}

// 	// Создаем тестовые маршруты
// 	testRoutes := []models.Route{
// 		{
// 			Name:        "Восхождение на Эльбрус",
// 			Description: "Легендарный маршрут к высочайшей точке Европы через живописные ледники и горные перевалы",
// 			Overview:    "Маршрут начинается от поселка Терскол и проходит через приют Бочки, скалы Пастухова до западной вершины Эльбруса",
// 			History:     "Первое успешное восхождение на Эльбрус было совершено в 1829 году экспедицией Российской академии наук под руководством генерала Г. А. Эммануэля",
// 			Distance:    22.5,
// 			Duration:    48.0,
// 			TypeID:      types[0].ID,
// 			AreaID:      areas[0].ID,
// 			Rating:      4.9,
// 			IsActive:    true,
// 		},
// 		{
// 			Name:        "Чегемские водопады",
// 			Description: "Путь к величественным водопадам в живописном Чегемском ущелье, известному своими суровыми скалами и бурной рекой",
// 			Overview:    "Маршрут проходит вдоль реки Чегем через несколько каскадов водопадов, самый известный из которых - Девичьи косы",
// 			History:     "Чегемское ущелье издавна было заселено балкарцами, о чем свидетельствуют древние склепы и башни",
// 			Distance:    8.2,
// 			Duration:    4.5,
// 			TypeID:      types[0].ID,
// 			AreaID:      areas[3].ID,
// 			Rating:      4.7,
// 			IsActive:    true,
// 		},
// 	}

// 	for i := range testRoutes {
// 		if err := db.Create(&testRoutes[i]).Error; err != nil {
// 			// log.Printf("Error creating route %d: %v", i, err)
// 			continue
// 		}
// 	}

// 	// Проверяем итоговое количество в базе
// 	var placesCount, imagesCount, reviewsCount, routesCount, usersCount int64
// 	if err := db.Model(&models.Place{}).Count(&placesCount).Error; err != nil {
// 		return err
// 	}
// 	if err := db.Model(&models.Image{}).Count(&imagesCount).Error; err != nil {
// 		return err
// 	}
// 	if err := db.Model(&models.Review{}).Count(&reviewsCount).Error; err != nil {
// 		return err
// 	}
// 	if err := db.Model(&models.Route{}).Count(&routesCount).Error; err != nil {
// 		return err
// 	}
// 	if err := db.Model(&models.User{}).Count(&usersCount).Error; err != nil {
// 		return err
// 	}

// 	// log.Printf("✅ Тестовые данные созданы: %d мест, %d изображений, %d отзывов, %d маршрутов, %d пользователей",
// 	// placesCount, imagesCount, reviewsCount, routesCount, usersCount)

// 	return nil
// }
