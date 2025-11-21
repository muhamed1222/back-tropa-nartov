package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	// Общие настройки
	AppName     string `env:"APP_NAME" default:"tropa-nartov"`
	Environment string `env:"ENVIRONMENT" default:"development"`
	Debug       bool   `env:"DEBUG" default:"true"`
	Host        string `env:"HOST" default:"localhost"`
	Port        string `env:"APP_PORT" default:"8001"`

	// Настройки базы данных
	DBHost     string `env:"POSTGRES_HOST" default:"localhost"`
	DBPort     string `env:"POSTGRES_PORT" default:"5432"`
	DBUser     string `env:"POSTGRES_USER" default:"postgres"`
	DBPassword string `env:"POSTGRES_PASSWORD" default:"password"`
	DBName     string `env:"POSTGRES_DB" default:"tropa_nartov"`

	// Яндекс SMTP настройки
	SMTPHost     string `env:"SMTP_HOST" default:"smtp.yandex.ru"`
	SMTPPort     string `env:"SMTP_PORT" default:"465"`
	SMTPUsername string `env:"SMTP_USERNAME" default:"tropanartov@yandex.ru"`
	SMTPPassword string `env:"SMTP_PASSWORD" default:"haavjputgiujiurr"`
	SMTPFrom     string `env:"SMTP_FROM" default:"tropanartov@yandex.ru"`
	SMTPUseTLS   bool   `env:"SMTP_USE_TLS" default:"true"`

	// JWT настройки
	JWTSecret        string `env:"JWT_SECRET_KEY"` // Обязательно: без default для безопасности
	JWTRefreshSecret string `env:"JWT_REFRESH_SECRET"`
	JWTExpiresIn     int    `env:"JWT_EXPIRES_IN" default:"24"`

	// CORS настройки
	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS" default:"http://localhost:3000,http://localhost:8080"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Отладочный вывод SMTP настроек
	// fmt.Printf("🔧 SMTP Configuration:\n")
	// fmt.Printf("   Host: %s\n", cfg.SMTPHost)
	// fmt.Printf("   Port: %s\n", cfg.SMTPPort)
	// fmt.Printf("   Username: %s\n", cfg.SMTPUsername)
	// fmt.Printf("   From: %s\n", cfg.SMTPFrom)
	// fmt.Printf("   UseTLS: %t\n", cfg.SMTPUseTLS)

	return cfg, nil
}
