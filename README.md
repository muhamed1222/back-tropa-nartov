# Backend Documentation - Tropa Nartov 🚀

REST API backend для мобильного приложения "Тропа Нартов", построенный на **Strapi CMS**.

## 📋 Содержание

- [Архитектура](#архитектура)
- [Установка](#установка)
- [Конфигурация](#конфигурация)
- [API Endpoints](#api-endpoints)
- [Деплой](#деплой)

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

## 📝 Документация Strapi

Дополнительная информация:
- [Официальная документация Strapi](https://docs.strapi.io/)
- [REST API документация](https://docs.strapi.io/developer-docs/latest/developer-resources/database-apis-reference/rest-api.html)
- [Filtering и Pagination](https://docs.strapi.io/developer-docs/latest/developer-resources/database-apis-reference/rest/filtering-locale-publication.html)

## 📄 Лицензия

MIT

---

**Статус миграции:** ✅ Миграция с Go на Strapi завершена (2024)
