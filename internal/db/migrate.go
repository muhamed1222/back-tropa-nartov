package db

import (
	"fmt"
	"log"
	"os"
	"tropa-nartov-backend/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Сначала устанавливаем необходимые расширения
	err := db.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error
	if err != nil {
		// log.Printf("Warning: Failed to create PostGIS extension: %v", err)
	}

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS cube").Error
	if err != nil {
		// log.Printf("Warning: Failed to create cube extension: %v", err)
	}

	err = db.Exec("CREATE EXTENSION IF NOT EXISTS earthdistance").Error
	if err != nil {
		// log.Printf("Warning: Failed to create earthdistance extension: %v", err)
	}

	// ВЫПОЛНЯЕМ AUTOMIGRATE ПЕРВЫМ - создаем все таблицы
	// ВАЖНО: Порядок имеет значение! Сначала создаем таблицы без внешних ключей,
	// затем таблицы, которые на них ссылаются.
	// log.Println("Creating database tables...")
	err = db.AutoMigrate(
		// 1. Сначала базовые таблицы без внешних ключей
		&models.User{},
		&models.Type{},      // Сначала типы (на них ссылаются routes)
		&models.Area{},      // Затем районы (на них ссылаются routes и places)
		&models.Category{},  // Категории (many-to-many)
		&models.Tag{},       // Теги (many-to-many)
		// 2. Затем таблицы с внешними ключами
		&models.Image{},     // Изображения
		&models.Route{},     // Маршруты (ссылается на Type и Area)
		&models.Review{},    // Отзывы (ссылается на Place и Route)
		&models.RouteStop{}, // Остановки маршрутов
		// 3. Связующие таблицы (many-to-many и избранное)
		&models.PlaceTag{},
		&models.RouteTag{},
		&models.PlaceCategory{},
		&models.RouteCategory{},
		&models.FavoritePlace{},
		&models.FavoriteRoute{},
		&models.PassedPlace{},
		&models.PassedRoute{},
		// 4. Дополнительные таблицы
		&models.ArticleCategory{},
		&models.Article{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	// Создаем дефолтные типы, если их нет (для корректной работы внешних ключей)
	if err := ensureDefaultTypes(db); err != nil {
		log.Printf("Warning: Failed to ensure default types: %v", err)
	}

	// Исправляем нарушенные внешние ключи в routes (если есть данные с несуществующими type_id)
	if err := fixBrokenForeignKeys(db); err != nil {
		log.Printf("Warning: Failed to fix broken foreign keys: %v", err)
	}

	// Отдельно мигрируем Place с кастомной логикой
	// log.Println("Migrating places table with custom logic...")
	if err := migratePlacesTable(db); err != nil {
		return fmt.Errorf("failed to migrate places table: %w", err)
	}

	// log.Println("Database tables created successfully")

	// ТЕПЕРЬ добавляем колонки, если их нет
	// Добавляем avatar_url, если отсутствует
	if !db.Migrator().HasColumn(&models.User{}, "avatar_url") {
		// log.Println("Adding avatar_url column to users table...")
		if err := db.Exec(`ALTER TABLE "users" ADD COLUMN "avatar_url" VARCHAR(500)`).Error; err != nil {
			log.Printf("Warning: Failed to add avatar_url column: %v", err)
		}
	}

	if !db.Migrator().HasColumn(&models.User{}, "password_hash") {
		// log.Println("Adding password_hash column to users table...")

		// Добавляем колонку без NOT NULL сначала
		if err := db.Exec(`ALTER TABLE "users" ADD COLUMN "password_hash" text`).Error; err != nil {
			return err
		}

		// Заполняем существующие записи временным паролем
		// log.Println("Setting temporary passwords for existing users...")
		if err := db.Exec(`UPDATE "users" SET "password_hash" = ? WHERE "password_hash" IS NULL`, "temporary_password_need_reset").Error; err != nil {
			return err
		}

		// Теперь добавляем ограничение NOT NULL
		if err := db.Exec(`ALTER TABLE "users" ALTER COLUMN "password_hash" SET NOT NULL`).Error; err != nil {
			return err
		}
	}

	// УДАЛЯЕМ старую колонку password если она есть
	if db.Migrator().HasColumn(&models.User{}, "password") {
		// log.Println("Removing old password column...")
		if err := db.Exec(`ALTER TABLE "users" DROP COLUMN "password"`).Error; err != nil {
			log.Printf("Warning: Failed to drop old password column: %v", err)
		}
	}

	// Добавляем avatar_url колонку, если её нет
	if !db.Migrator().HasColumn(&models.User{}, "avatar_url") {
		// log.Println("Adding avatar_url column to users table...")
		if err := db.Exec(`ALTER TABLE "users" ADD COLUMN "avatar_url" VARCHAR(500)`).Error; err != nil {
			log.Printf("Warning: Failed to add avatar_url column: %v", err)
		}
	}

	// Добавляем GIST-индекс для гео-поиска (PostGIS вариант)
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_places_geo ON places USING GIST (ST_Point(longitude, latitude))").Error
	if err != nil {
		// log.Printf("Warning: Failed to create geo index: %v", err)
		// Пробуем альтернативный вариант с earthdistance
		err = db.Exec("CREATE INDEX IF NOT EXISTS idx_places_geo ON places USING GIST (ll_to_earth(latitude, longitude))").Error
		if err != nil {
			// log.Printf("Warning: Failed to create earthdistance index: %v", err)
		}
	}

	// Добавляем индексы для производительности
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_reviews_created_at ON reviews (created_at DESC)").Error
	if err != nil {
		return err
	}

	// Добавляем индекс для images
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_images_place_id ON images(place_id)").Error
	if err != nil {
		return err
	}

	// Добавляем индекс для reviews по place_id
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_reviews_place_id ON reviews(place_id)").Error
	if err != nil {
		return err
	}

	// Выполняем SQL для триггеров (если файл существует)
	if _, err := os.Stat("internal/db/triggers.sql"); err == nil {
		triggersSQL, err := os.ReadFile("internal/db/triggers.sql")
		if err != nil {
			// log.Printf("Warning: Failed to read triggers.sql: %v", err)
		} else {
			err = db.Exec(string(triggersSQL)).Error
			if err != nil {
				// log.Printf("Warning: Failed to execute triggers: %v", err)
			}
		}
	} else {
		// log.Println("No triggers.sql file found, skipping triggers")
	}

	// log.Println("Database migration completed successfully")
	return nil
}

// migratePlacesTable выполняет кастомную миграцию для таблицы places
func migratePlacesTable(db *gorm.DB) error {
	// Создаем временную структуру без новых полей
	type PlaceTemp struct {
		ID           uint    `gorm:"primaryKey"`
		Name         string  `gorm:"type:varchar(200);not null"`
		Description  string  `gorm:"type:text;not null"`
		Overview     string  `gorm:"type:text"`
		History      string  `gorm:"type:text"`
		Address      string  `gorm:"type:varchar(500);not null"`
		Latitude     float64 `gorm:"type:decimal(10,8);not null"`
		Longitude    float64 `gorm:"type:decimal(11,8);not null"`
		OpeningHours string  `gorm:"type:varchar(100)"`
		Contacts     string  `gorm:"type:jsonb"`
		Rating       float32 `gorm:"type:decimal(2,1);default:0"`
		TypeID       uint    `gorm:"not null"`
		AreaID       uint    `gorm:"not null"`
		IsActive     bool    `gorm:"default:true"`
	}

	// Сначала мигрируем базовую структуру
	if err := db.AutoMigrate(&PlaceTemp{}); err != nil {
		return err
	}

	// Теперь добавляем новые колонки по одной с правильной логикой

	// 1. Добавляем type как nullable
	if !db.Migrator().HasColumn(&models.Place{}, "type") {
		log.Println("Adding type column to places table...")
		if err := db.Exec(`ALTER TABLE "places" ADD COLUMN "type" VARCHAR(100)`).Error; err != nil {
			return err
		}

		// Заполняем существующие записи значением по умолчанию
		// log.Println("Setting default values for type column...")
		if err := db.Exec(`UPDATE "places" SET "type" = 'достопримечательность' WHERE "type" IS NULL`).Error; err != nil {
			return err
		}

		// Теперь делаем колонку NOT NULL
		if err := db.Exec(`ALTER TABLE "places" ALTER COLUMN "type" SET NOT NULL`).Error; err != nil {
			return err
		}
	}

	// 2. Добавляем hours
	if !db.Migrator().HasColumn(&models.Place{}, "hours") {
		// log.Println("Adding hours column to places table...")
		if err := db.Exec(`ALTER TABLE "places" ADD COLUMN "hours" VARCHAR(200)`).Error; err != nil {
			return err
		}
	}

	// 3. Добавляем weekend
	if !db.Migrator().HasColumn(&models.Place{}, "weekend") {
		// log.Println("Adding weekend column to places table...")
		if err := db.Exec(`ALTER TABLE "places" ADD COLUMN "weekend" VARCHAR(100)`).Error; err != nil {
			return err
		}
	}

	// 4. Добавляем entry
	if !db.Migrator().HasColumn(&models.Place{}, "entry") {
		// log.Println("Adding entry column to places table...")
		if err := db.Exec(`ALTER TABLE "places" ADD COLUMN "entry" VARCHAR(100)`).Error; err != nil {
			return err
		}
	}

	// 5. Добавляем contacts_email
	if !db.Migrator().HasColumn(&models.Place{}, "contacts_email") {
		// log.Println("Adding contacts_email column to places table...")
		if err := db.Exec(`ALTER TABLE "places" ADD COLUMN "contacts_email" VARCHAR(200)`).Error; err != nil {
			return err
		}
	}

	// 6. Обновляем contacts если нужно (меняем тип с jsonb на varchar)
	if db.Migrator().HasColumn(&models.Place{}, "contacts") {
		// Проверяем тип колонки
		var columnType string
		db.Raw(`SELECT data_type FROM information_schema.columns 
		        WHERE table_name = 'places' AND column_name = 'contacts'`).Scan(&columnType)

		if columnType == "jsonb" {
			// log.Println("Converting contacts column from jsonb to varchar...")
			// Создаем временную колонку
			if err := db.Exec(`ALTER TABLE "places" ADD COLUMN "contacts_temp" VARCHAR(200)`).Error; err != nil {
				return err
			}

			// Копируем данные (извлекаем телефон из JSON)
			if err := db.Exec(`UPDATE "places" SET "contacts_temp" = COALESCE("contacts"->>'phone', '')`).Error; err != nil {
				return err
			}

			// Удаляем старую колонку
			if err := db.Exec(`ALTER TABLE "places" DROP COLUMN "contacts"`).Error; err != nil {
				return err
			}

			// Переименовываем временную колонку
			if err := db.Exec(`ALTER TABLE "places" RENAME COLUMN "contacts_temp" TO "contacts"`).Error; err != nil {
				return err
			}
		}
	} else {
		// log.Println("Adding contacts column to places table...")
		if err := db.Exec(`ALTER TABLE "places" ADD COLUMN "contacts" VARCHAR(200)`).Error; err != nil {
			return err
		}
	}

	return nil
}

// ensureDefaultTypes создает дефолтные типы маршрутов, если их нет
func ensureDefaultTypes(db *gorm.DB) error {
	var count int64
	db.Model(&models.Type{}).Count(&count)
	
	if count == 0 {
		// Создаем дефолтные типы
		defaultTypes := []models.Type{
			{Name: "Пеший поход", Description: "Пешие маршруты по горам и лесам"},
			{Name: "Веломаршрут", Description: "Велосипедные трассы"},
			{Name: "Автомобильный", Description: "Маршруты для автомобильных поездок"},
			{Name: "Комбинированный", Description: "Маршруты с разными видами передвижения"},
		}
		
		for _, t := range defaultTypes {
			if err := db.FirstOrCreate(&t, models.Type{Name: t.Name}).Error; err != nil {
				return err
			}
		}
	}
	
	return nil
}

// fixBrokenForeignKeys исправляет записи с несуществующими внешними ключами
func fixBrokenForeignKeys(db *gorm.DB) error {
	// Получаем первый тип (дефолтный) для замены несуществующих
	var defaultType models.Type
	if err := db.First(&defaultType).Error; err != nil {
		// Если нет типов, пропускаем
		return nil
	}
	
	// Исправляем routes с несуществующими type_id
	result := db.Exec(`
		UPDATE routes 
		SET type_id = ? 
		WHERE type_id NOT IN (SELECT id FROM types)
		AND type_id IS NOT NULL
	`, defaultType.ID)
	
	if result.Error != nil {
		return result.Error
	}
	
	// Получаем первый район (дефолтный) для замены несуществующих
	var defaultArea models.Area
	if err := db.First(&defaultArea).Error; err != nil {
		// Если нет районов, пропускаем
		return nil
	}
	
	// Исправляем routes с несуществующими area_id
	result = db.Exec(`
		UPDATE routes 
		SET area_id = ? 
		WHERE area_id NOT IN (SELECT id FROM areas)
		AND area_id IS NOT NULL
	`, defaultArea.ID)
	
	if result.Error != nil {
		return result.Error
	}
	
	return nil
}
