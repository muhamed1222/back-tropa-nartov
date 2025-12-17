# Backend Documentation - Tropa Nartov 🚀

REST API backend для мобильного приложения "Тропа Нартов", построенный на **Strapi CMS**.

## 📋 Содержание

- [Архитектура](#архитектура)
- [Установка](#установка)
- [Конфигурация](#конфигурация)
- [API Endpoints](#api-endpoints)
- [Схема базы данных](#схема-базы-данных)
- [Деплой](#деплой)
- [Git-репозиторий](#git-репозиторий)
- [Документация](#документация)

## 🏗️ Архитектура

### Strapi-Only Architecture

Проект полностью мигрирован на **Strapi CMS** в качестве единственного backend решения. Go API был полностью удален.

```
back/
├── strapi/                  # Strapi CMS (единственный backend)
│   ├── src/
│   │   ├── api/            # Content types
│   │   │   ├── place/      # Места (Place)
│   │   │   ├── route/      # Маршруты (Route)
│   │   │   ├── review/     # Отзывы (Review)
│   │   │   ├── favorite/   # Избранное (Favorite)
│   │   │   ├── visited-place/ # История посещений
│   │   │   ├── category/   # Категории
│   │   │   └── tag/        # Теги
│   │   ├── extensions/     # Пользовательские расширения
│   │   └── middlewares/    # Middleware
│   ├── config/             # Конфигурация Strapi
│   ├── public/             # Статические файлы
│   └── database/           # SQLite база данных (dev) / PostgreSQL (prod)
├── .env                    # Переменные окружения
├── .env.example            # Пример конфигурации
└── README.md               # Этот файл
```

### Технологии

- **Strapi v4** - headless CMS с REST API
- **Node.js 18+** - runtime
- **SQLite** - база данных для разработки
- **PostgreSQL** (опционально) - база данных для production

## 📦 Установка

### Требования

- Node.js 18+ и npm/yarn
- Git

### Быстрый старт

1. **Клонировать репозиторий**

```bash
git clone <repository-url>
cd back/strapi
```

2. **Установить зависимости**

```bash
npm install
# или
yarn install
```

3. **Настроить переменные окружения**

Создайте файл `.env` в корне директории `strapi`:

```bash
cp .env.example .env
```

Основные переменные:

```env
# Server
HOST=0.0.0.0
PORT=1337

# Admin panel
ADMIN_JWT_SECRET=your-admin-jwt-secret

# API tokens
API_TOKEN_SALT=your-api-token-salt
JWT_SECRET=your-jwt-secret

# Database (SQLite по умолчанию, для production используйте PostgreSQL)
DATABASE_CLIENT=sqlite
DATABASE_FILENAME=.tmp/data.db

# Для PostgreSQL (production):
# DATABASE_CLIENT=postgres
# DATABASE_HOST=localhost
# DATABASE_PORT=5432
# DATABASE_NAME=strapi
# DATABASE_USERNAME=strapi
# DATABASE_PASSWORD=strapi
```

4. **Запустить Strapi**

```bash
npm run develop
# или
yarn develop
```

Strapi будет доступен по адресу: `http://localhost:1337`

5. **Создать администратора**

При первом запуске откройте `http://localhost:1337/admin` и создайте учетную запись администратора.

## ⚙️ Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `HOST` | Хост сервера | `0.0.0.0` |
| `PORT` | Порт сервера | `1337` |
| `ADMIN_JWT_SECRET` | Секрет для JWT админ-панели | - |
| `API_TOKEN_SALT` | Salt для API токенов | - |
| `JWT_SECRET` | Секрет для JWT пользователей | - |
| `DATABASE_CLIENT` | Тип БД (sqlite/postgres) | `sqlite` |

### Strapi Admin Panel

Админ-панель доступна по адресу: `http://localhost:1337/admin`

Здесь вы можете:
- Управлять контентом (места, маршруты, отзывы)
- Настраивать роли и разрешения
- Управлять медиафайлами
- Просматривать аналитику

## 🔌 API Endpoints

### Базовый URL

```
http://localhost:1337/api
```

### Аутентификация

#### Регистрация

```http
POST /api/auth/local/register
Content-Type: application/json

{
  "username": "user",
  "email": "user@example.com",
  "password": "password123"
}
```

Ответ:

```json
{
  "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "user",
    "email": "user@example.com"
  }
}
```

#### Вход

```http
POST /api/auth/local
Content-Type: application/json

{
  "identifier": "user@example.com",
  "password": "password123"
}
```

#### Получение профиля

```http
GET /api/users/me
Authorization: Bearer <jwt_token>
```

#### Восстановление пароля

```http
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

```http
POST /api/auth/reset-password
Content-Type: application/json

{
  "code": "reset_code",
  "password": "newpassword123",
  "passwordConfirmation": "newpassword123"
}
```

### Места (Places)

```http
GET /api/places
GET /api/places/:id
GET /api/places?populate=*
GET /api/places?filters[category][id][$eq]=1
```

### Маршруты (Routes)

```http
GET /api/routes
GET /api/routes/:id
GET /api/routes?populate=*
GET /api/routes?filters[route_type][id][$eq]=1
```

### Отзывы (Reviews)

```http
GET /api/reviews
GET /api/reviews?filters[place][id][$eq]=1
GET /api/reviews?filters[route][id][$eq]=1
POST /api/reviews
Authorization: Bearer <jwt_token>
```

### Избранное (Favorites)

```http
GET /api/favorites?filters[user_id][$eq]=<user_id>
POST /api/favorites
Authorization: Bearer <jwt_token>

{
  "data": {
    "user_id": "1",
    "place": 1,  // или "route": 1
  }
}

DELETE /api/favorites/:id
Authorization: Bearer <jwt_token>
```

### История посещений (Visited Places)

```http
GET /api/visited-places?filters[user_id][$eq]=<user_id>
POST /api/visited-places
Authorization: Bearer <jwt_token>

{
  "data": {
    "user_id": "1",
    "place": 1,  // или "route": 1
    "visited_at": "2024-01-01T00:00:00.000Z"
  }
}
```

### Загрузка файлов (Media Library)

```http
POST /api/upload
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data

files: <file>
ref: "api::user.user"
refId: <user_id>
field: "avatar"
```

## 🚀 Деплой

### Production Build

```bash
cd strapi
npm run build
NODE_ENV=production npm start
```

### Docker (опционально)

Создайте `Dockerfile` в директории `strapi`:

```dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .

RUN npm run build

EXPOSE 1337

CMD ["npm", "start"]
```

Запуск:

```bash
docker build -t tropa-nartov-backend .
docker run -p 1337:1337 tropa-nartov-backend
```

## 🗄️ Схема базы данных

### Основные сущности

```
┌─────────────────┐
│     Area        │ (Районы)
│  - id           │
│  - name         │
│  - slug         │
│  - order        │
│  - is_active    │
└────────┬────────┘
         │ 1
         │
         │ N
┌────────▼────────┐
│     Place       │ (Места)
│  - id           │
│  - name         │
│  - slug         │
│  - images       │
│  - history      │
│  - address      │
│  - latitude     │
│  - longitude    │
│  - working_hours│
│  - phone        │
│  - website      │
│  - rating       │
│  - is_active    │
└────────┬────────┘
         │
    ┌────┴────┬──────────────┬──────────────┬──────────────┐
    │         │              │              │              │
    │ N       │ N            │ N            │ N            │ N
    │         │              │              │              │
┌───▼───┐ ┌──▼──────┐  ┌─────▼─────┐  ┌─────▼─────┐  ┌────▼──────┐
│Route  │ │Category │  │   Tag     │  │  Review   │  │ Favorite  │
│       │ │         │  │           │  │           │  │           │
│ - id  │ │ - id    │  │ - id      │  │ - id      │  │ - id      │
│ - name│ │ - name  │  │ - name    │  │ - rating  │  │ - user_id │
│ - slug│ │ - slug  │  │ - slug    │  │ - text    │  │           │
│       │ │         │  │           │  │           │  │           │
└───┬───┘ └─────────┘  └───────────┘  └───────────┘  └───────────┘
    │
    │ N
    │
┌───▼────────┐
│ RouteType  │ (Типы маршрутов)
│  - id      │
│  - name    │
│  - slug    │
│  - order   │
│  - is_active│
└────────────┘

┌─────────────────┐
│  VisitedPlace   │ (История посещений)
│  - id           │
│  - user_id      │
│  - visited_at   │
└─────────────────┘
```

### Связи между сущностями

| Сущность 1 | Связь | Сущность 2 | Описание |
|------------|-------|------------|----------|
| `Area` | 1:N | `Place` | Один район может содержать много мест |
| `Place` | N:M | `Route` | Места могут входить в несколько маршрутов |
| `Place` | N:M | `Category` | Место может иметь несколько категорий |
| `Place` | N:M | `Tag` | Место может иметь несколько тегов |
| `Route` | N:1 | `RouteType` | Маршрут принадлежит одному типу |
| `Place` | 1:N | `Review` | Место может иметь много отзывов |
| `Route` | 1:N | `Review` | Маршрут может иметь много отзывов |
| `Place` | 1:N | `Favorite` | Место может быть в избранном у многих пользователей |
| `Route` | 1:N | `Favorite` | Маршрут может быть в избранном у многих пользователей |
| `Place` | 1:N | `VisitedPlace` | Место может быть посещено многими пользователями |
| `Route` | 1:N | `VisitedPlace` | Маршрут может быть пройден многими пользователями |

### Описание таблиц

#### `places` (Места)
Основная сущность для хранения информации о достопримечательностях и местах.

**Ключевые поля:**
- `name` - название места (обязательное)
- `slug` - уникальный идентификатор для URL
- `latitude`, `longitude` - координаты для карты
- `images` - массив изображений
- `rating` - рейтинг (0.0 - 5.0)

#### `routes` (Маршруты)
Туристические маршруты, состоящие из нескольких мест.

**Ключевые поля:**
- `name` - название маршрута (обязательное)
- `slug` - уникальный идентификатор
- `description` - описание маршрута
- `route_type` - связь с типом маршрута (пеший, авто и т.д.)
- `places` - связь many-to-many с местами

#### `reviews` (Отзывы)
Отзывы пользователей о местах и маршрутах.

**Ключевые поля:**
- `rating` - оценка (1-5)
- `text` - текст отзыва
- `place` или `route` - связь с местом или маршрутом

#### `categories` (Категории)
Категории для фильтрации мест (например: "Музеи", "Парки", "Горы").

#### `tags` (Теги)
Теги для дополнительной фильтрации (например: "горы", "водопад", "река").

#### `areas` (Районы)
Районы Кабардино-Балкарии для географической фильтрации.

#### `route_types` (Типы маршрутов)
Типы маршрутов (например: "Пеший", "Автомобильный", "Велосипедный").

#### `favorites` (Избранное)
Избранные места и маршруты пользователей.

**Ключевые поля:**
- `user_id` - ID пользователя
- `place` или `route` - связь с местом или маршрутом

#### `visited_places` (История посещений)
История посещений мест и прохождения маршрутов.

**Ключевые поля:**
- `user_id` - ID пользователя
- `place` или `route` - связь с местом или маршрутом
- `visited_at` - дата и время посещения

## 🔗 Git-репозиторий

**Репозиторий:** [https://github.com/muhamed1222/back-tropa-nartov.git](https://github.com/muhamed1222/back-tropa-nartov.git)

**Финальная ветка:** `main`

### Клонирование репозитория

```bash
git clone https://github.com/muhamed1222/back-tropa-nartov.git
cd back-tropa-nartov
git checkout main
```

## 📚 Документация

### Официальная документация Strapi

- **[Strapi Documentation](https://docs.strapi.io/)** - Полная документация Strapi CMS
- **[REST API Reference](https://docs.strapi.io/developer-docs/latest/developer-resources/database-apis-reference/rest-api.html)** - Справочник по REST API
- **[Filtering и Pagination](https://docs.strapi.io/developer-docs/latest/developer-resources/database-apis-reference/rest/filtering-locale-publication.html)** - Фильтрация и пагинация данных
- **[Content Types](https://docs.strapi.io/developer-docs/latest/development/backend-customization/models.html)** - Работа с типами контента
- **[Authentication](https://docs.strapi.io/developer-docs/latest/plugins/users-permissions.html)** - Аутентификация и авторизация

### Полезные ссылки

- **[Strapi Admin Panel Guide](https://docs.strapi.io/user-docs/latest/getting-started/introduction.html)** - Руководство по админ-панели
- **[Deployment Guide](https://docs.strapi.io/developer-docs/latest/setup-deployment-guides/deployment.html)** - Руководство по деплою
- **[Database Configuration](https://docs.strapi.io/developer-docs/latest/setup-deployment-guides/configurations/databases.html)** - Настройка баз данных

## 📄 Лицензия

MIT

---

**Статус миграции:** ✅ Миграция с Go на Strapi завершена (2024)
