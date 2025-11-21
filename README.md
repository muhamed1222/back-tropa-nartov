# Backend Documentation - Tropa Nartov 🚀

REST API backend для мобильного приложения "Тропа Нартов", построенный на Go и Strapi.

## 📋 Содержание

- [Архитектура](#архитектура)
- [Установка](#установка)
- [Конфигурация](#конфигурация)
- [API Endpoints](#api-endpoints)
- [База данных](#база-данных)
- [Тестирование](#тестирование)
- [Деплой](#деплой)

## 🏗️ Архитектура

### Компоненты

```
back/
├── cmd/api/              # Точка входа приложения
│   └── main.go          # Инициализация сервера, middleware, routes
├── internal/            # Внутренняя бизнес-логика
│   ├── auth/           # JWT аутентификация
│   ├── config/         # Конфигурация из .env
│   ├── db/             # Database connection & migrations
│   ├── logger/         # Zap structured logging
│   ├── middleware/     # Rate limiting, ETag, etc.
│   ├── models/         # GORM модели
│   ├── places/         # Places service
│   ├── routes/         # Routes handlers
│   ├── reviews/        # Reviews service
│   ├── sync/           # Strapi synchronization
│   ├── validation/     # Input validators
│   └── metrics/        # Prometheus metrics
├── strapi/             # Strapi CMS
│   ├── src/api/       # Content types (Place, Route, etc.)
│   └── config/        # Strapi configuration
├── docs/               # Swagger documentation (auto-generated)
├── scripts/            # Utility scripts (backup, etc.)
└── docker-compose.yaml # Development environment
```

### Технологии

- **Go 1.22** - основной язык
- **Gin** - HTTP framework с высокой производительностью
- **GORM** - ORM для работы с PostgreSQL
- **PostgreSQL 15 + PostGIS** - база данных с поддержкой геоданных
- **Strapi v4** - headless CMS для управления контентом
- **Redis** - кэширование (опционально)
- **JWT** - stateless аутентификация
- **Zap** - структурированное логирование
- **Prometheus** - метрики

## 🚀 Установка

### Предварительные требования

- Go 1.22+
- Docker & Docker Compose
- PostgreSQL 15+ (если запуск без Docker)
- Node.js 18+ (для Strapi)

### Быстрый старт с Docker

```bash
# 1. Клонируйте репозиторий
git clone https://github.com/yourusername/tropa-nartov.git
cd tropa-nartov/back

# 2. Создайте .env файл
cp .env.example .env
nano .env  # Настройте переменные

# 3. Запустите все сервисы
docker-compose up -d

# 4. Проверьте статус
docker-compose ps

# 5. Выполните миграции (если нужно)
docker exec -it tropa-go-api go run cmd/api/main.go migrate
```

### Локальная установка (без Docker)

```bash
# 1. Установите зависимости Go
go mod download

# 2. Настройте PostgreSQL
createdb tropa_nartov
psql -d tropa_nartov -c "CREATE EXTENSION IF NOT EXISTS postgis;"

# 3. Настройте .env
cp .env.example .env

# 4. Запустите миграции
go run cmd/api/main.go migrate

# 5. Запустите сервер
go run cmd/api/main.go

# Или скомпилируйте:
go build -o bin/api cmd/api/main.go
./bin/api
```

## ⚙️ Конфигурация

### Переменные окружения (.env)

```bash
# Приложение
APP_NAME=tropa-nartov
ENVIRONMENT=development  # development | production
DEBUG=true
HOST=0.0.0.0
APP_PORT=8001

# База данных
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=tropa_nartov

# JWT (ВАЖНО: установите secure secret в production!)
JWT_SECRET_KEY=your-super-secret-jwt-key-min-32-chars
JWT_EXPIRES_IN=24  # часы

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080

# SMTP (для восстановления пароля)
SMTP_HOST=smtp.yandex.ru
SMTP_PORT=465
SMTP_USERNAME=tropanartov@yandex.ru
SMTP_PASSWORD=your_smtp_password
SMTP_FROM=tropanartov@yandex.ru
SMTP_USE_TLS=true

# Strapi Sync (опционально)
STRAPI_URL=http://localhost:1337
STRAPI_API_TOKEN=your_strapi_token

# Логирование
LOG_LEVEL=debug  # debug | info | warn | error
```

### Генерация JWT Secret

```bash
# Используйте один из методов:
openssl rand -base64 48
# или
head /dev/urandom | tr -dc A-Za-z0-9 | head -c 48
```

## 📡 API Endpoints

### Swagger Documentation

Интерактивная документация доступна по адресу:
```
http://localhost:8001/swagger/index.html
```

### Основные endpoints

#### Аутентификация

```http
POST   /auth/register           # Регистрация нового пользователя
POST   /auth/login              # Вход в систему (получение JWT)
POST   /auth/forgot-password    # Запрос восстановления пароля
POST   /auth/verify-reset-code  # Проверка кода восстановления
POST   /auth/reset-password     # Сброс пароля
GET    /auth/profile            # Получение профиля (требует JWT)
PUT    /auth/profile            # Обновление профиля (требует JWT)
POST   /auth/upload-avatar      # Загрузка аватара (требует JWT)
DELETE /auth/delete-account     # Удаление аккаунта (требует JWT)
```

#### Места (Places)

```http
GET    /places                  # Список мест (с пагинацией и фильтрами)
GET    /places/:id              # Детальная информация о месте
POST   /places                  # Создание места (admin)
PUT    /places/:id              # Обновление места (admin)
DELETE /places/:id              # Удаление места (admin)

# Query параметры для /places:
?page=1                # Номер страницы (default: 1)
?limit=20              # Элементов на страницу (default: 20, max: 100)
?type_id=1,2           # Фильтр по типам
?area_id=3             # Фильтр по району
?category_id=4,5       # Фильтр по категориям
?light=true            # Легкий формат (без relations)
```

#### Маршруты (Routes)

```http
GET    /routes                  # Список маршрутов (с пагинацией)
GET    /routes/:id              # Детальная информация о маршруте
POST   /routes                  # Создание маршрута (admin)
PUT    /routes/:id              # Обновление маршрута (admin)
DELETE /routes/:id              # Удаление маршрута (admin)
```

#### Отзывы (Reviews)

```http
GET    /reviews/place/:placeId  # Отзывы о месте
GET    /reviews/route/:routeId  # Отзывы о маршруте
POST   /reviews                 # Создание отзыва (требует JWT)
PUT    /reviews/:id             # Обновление отзыва (требует JWT)
DELETE /reviews/:id             # Удаление отзыва (требует JWT)
```

#### Избранное (Favorites)

```http
GET    /favorites/places        # Список избранных мест (требует JWT)
POST   /favorites/places/:id    # Добавить в избранное (требует JWT)
DELETE /favorites/places/:id    # Удалить из избранного (требует JWT)
GET    /favorites/places/:id/status  # Проверка статуса (требует JWT)

# Аналогично для маршрутов:
/favorites/routes/*
```

#### Утилиты

```http
GET    /ping                    # Health check
GET    /metrics                 # Prometheus метрики
GET    /swagger/*any            # Swagger UI
```

### Пример запроса с JWT

```bash
# 1. Логин
curl -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'

# Ответ: {"token": "eyJhbGciOiJIUzI1NiIs..."}

# 2. Использование токена
curl -X GET http://localhost:8001/auth/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

### Пагинированный ответ

```json
{
  "data": [...],
  "total": 150,
  "page": 1,
  "limit": 20,
  "total_pages": 8,
  "has_more": true
}
```

## 🗄️ База данных

### Схема

```sql
-- Основные таблицы
users              # Пользователи
places             # Достопримечательности
routes             # Туристические маршруты
reviews            # Отзывы
images             # Изображения
types              # Типы (музей, парк, etc.)
areas              # Районы
categories         # Категории
tags               # Теги

-- Связующие таблицы
favorite_places    # Избранные места
favorite_routes    # Избранные маршруты
passed_routes      # Пройденные маршруты
route_stops        # Остановки маршрутов
place_categories   # Many-to-many: places <-> categories
route_categories   # Many-to-many: routes <-> categories
```

### Миграции

```bash
# Выполнить миграции
go run cmd/api/main.go migrate

# Или через Docker
docker exec -it tropa-go-api go run cmd/api/main.go migrate

# Миграции выполняются автоматически при старте сервера
```

### Backup

```bash
# Ручной backup
./scripts/backup_db.sh

# Автоматический backup (crontab)
0 3 * * * /path/to/scripts/backup_db.sh

# Restore
psql -U postgres -d tropa_nartov < backup_20250121_030000.sql.gz
```

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты
go test ./... -v

# С coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Конкретный пакет
go test ./internal/auth -v

# Benchmark
go test ./... -bench=. -benchmem
```

### Написание тестов

```go
// internal/auth/service_test.go
func TestLogin(t *testing.T) {
    db := setupTestDB()
    service := NewService(db, &config.Config{
        JWTSecret: "test-secret",
    })
    
    // Test logic
    token, err := service.Login("test@example.com", "password")
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
}
```

## 📊 Мониторинг

### Prometheus Metrics

Доступны на `/metrics`:

```
# HTTP метрики
http_requests_total{method="GET",endpoint="/places",status="200"}
http_request_duration_seconds{method="GET",endpoint="/places"}

# Database метрики
db_queries_total{operation="SELECT",table="places"}
db_query_duration_seconds{operation="SELECT",table="places"}
```

### Логирование

Структурированные логи через Zap:

```go
logger.Info("Server starting",
    zap.String("port", cfg.Port),
    zap.String("environment", cfg.Environment),
)
```

## 🚢 Деплой

### Production checklist

- [ ] Установите безопасный `JWT_SECRET_KEY` (min 32 символа)
- [ ] Используйте сильные пароли БД
- [ ] Настройте `CORS_ALLOWED_ORIGINS` (не используйте "*")
- [ ] Установите `ENVIRONMENT=production` и `DEBUG=false`
- [ ] Настройте HTTPS/SSL (через nginx)
- [ ] Настройте автоматический backup БД
- [ ] Настройте мониторинг (Prometheus + Grafana)
- [ ] Настройте логи (aggregation в централизованное хранилище)

### Docker Production

```bash
# Build
docker build -t tropa-nartov-api:latest .

# Run
docker run -d \
  --name tropa-api \
  -p 8001:8001 \
  --env-file .env.production \
  tropa-nartov-api:latest
```

### Systemd Service

```bash
# /etc/systemd/system/tropa-api.service
[Unit]
Description=Tropa Nartov API
After=network.target postgresql.service

[Service]
Type=simple
User=tropa
WorkingDirectory=/opt/tropa-nartov
EnvironmentFile=/opt/tropa-nartov/.env
ExecStart=/opt/tropa-nartov/bin/api
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## 🔧 Troubleshooting

### Проблема: "JWT_SECRET_KEY не установлен"

```bash
# Убедитесь, что .env файл загружен
cat .env | grep JWT_SECRET_KEY

# Сгенерируйте новый секрет
openssl rand -base64 48
```

### Проблема: "Database connection failed"

```bash
# Проверьте статус PostgreSQL
docker-compose ps postgres
# или
systemctl status postgresql

# Проверьте переменные окружения
echo $POSTGRES_HOST $POSTGRES_PORT
```

### Проблема: "Rate limit exceeded"

```bash
# Временно увеличьте лимит в middleware/ratelimit.go
# Или подождите 1 минуту для сброса счетчика
```

## 📚 Дополнительные ресурсы

- [Swagger UI](http://localhost:8001/swagger/index.html)
- [Prometheus Metrics](http://localhost:8001/metrics)
- [Strapi Admin](http://localhost:1337/admin)
- [API Documentation](../docs/API_DOCUMENTATION.md)
- [Architecture Guide](../docs/ARCHITECTURE.md)

## 🤝 Contributing

См. [CONTRIBUTING.md](../CONTRIBUTING.md)

## 📄 License

MIT License - см. [LICENSE](../LICENSE)

