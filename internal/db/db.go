package db

import (
	"fmt"
	"log"
	"tropa-nartov-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect подключается к PostgreSQL
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	// Настраиваем logger для GORM - не логируем ErrRecordNotFound (это ожидаемое поведение)
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * 1_000_000, // 200ms
			LogLevel:                  logger.Warn,      // Логируем только предупреждения и ошибки
			IgnoreRecordNotFoundError: true,            // Не логируем "record not found" - это нормально
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// log.Println("Database connected successfully")
	return db, nil
}
