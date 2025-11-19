# Настройка Strapi для проекта "Тропа Нартов"

## ✅ Что было сделано

1. **Заменен Adminer на Strapi** в `docker-compose.yaml`
2. **Создана структура** для Strapi в папке `./strapi`
3. **Настроено подключение** к существующей PostgreSQL БД

## 🚀 Запуск Strapi

### Первый запуск

```bash
cd "/Users/kelemetovmuhamed/Documents/тропа нартов /back"
docker-compose up -d strapi
```

### Проверка статуса

```bash
docker-compose ps strapi
docker-compose logs -f strapi
```

## 📝 Первоначальная настройка

1. **Откройте админ-панель**: http://localhost:1337/admin

2. **Создайте первого администратора**:
   - Имя пользователя
   - Email
   - Пароль

3. **Начните создавать Content Types** для ваших данных

## 🗂️ Рекомендуемые Content Types

Создайте следующие Content Types в Strapi для управления контентом:

### 1. Place (Место)
- `name` (Text, required)
- `type` (Text)
- `description` (Rich Text)
- `overview` (Rich Text)
- `history` (Rich Text)
- `address` (Text, required)
- `latitude` (Number, decimal)
- `longitude` (Number, decimal)
- `hours` (Text)
- `weekend` (Text)
- `entry` (Text)
- `contacts` (Text)
- `contacts_email` (Email)
- `rating` (Number, decimal)
- `images` (Media, multiple)
- `categories` (Relation, many-to-many)
- `tags` (Relation, many-to-many)

### 2. Route (Маршрут)
- `name` (Text, required)
- `description` (Rich Text)
- `overview` (Rich Text)
- `history` (Rich Text)
- `distance` (Number, decimal)
- `duration` (Number, decimal)
- `rating` (Number, decimal)
- `type` (Relation, many-to-one)
- `area` (Relation, many-to-one)
- `categories` (Relation, many-to-many)

### 3. Review (Отзыв)
- `text` (Rich Text, required)
- `rating` (Number, integer, 1-5)
- `place` (Relation, many-to-one, optional)
- `route` (Relation, many-to-one, optional)
- `user` (Relation, many-to-one)
- `likes` (Number, integer, default: 0)

### 4. Category (Категория)
- `name` (Text, required)
- `description` (Text)
- `slug` (UID, from name)

### 5. Tag (Тег)
- `name` (Text, required)
- `slug` (UID, from name)

## 🔗 Интеграция с существующим API

Strapi будет работать **параллельно** с вашим Go API:

- **Go API**: http://localhost:8001 (основной API)
- **Strapi API**: http://localhost:1337/api (админ API)
- **Strapi Admin**: http://localhost:1337/admin (админ-панель)

### Варианты интеграции:

1. **Использовать Strapi только для админки**
   - Управление контентом через Strapi
   - Чтение данных через Go API

2. **Синхронизация данных**
   - Создавать/обновлять через Strapi
   - Синхронизировать с Go API через webhooks или скрипты

3. **Полная миграция на Strapi**
   - Использовать Strapi API вместо Go API
   - Настроить кастомные контроллеры в Strapi

## 🔧 Настройка переменных окружения

Strapi использует переменные из `.env` файла:
- `POSTGRES_DB` - имя БД
- `POSTGRES_USER` - пользователь БД
- `POSTGRES_PASSWORD` - пароль БД
- `JWT_SECRET_KEY` - для JWT токенов

## 📦 Структура файлов

```
back/
├── strapi/
│   ├── Dockerfile          # (не используется, используем готовый образ)
│   ├── package.json        # Зависимости Strapi
│   ├── config/
│   │   ├── database.js     # Конфигурация БД
│   │   ├── server.js       # Настройки сервера
│   │   ├── admin.js        # Настройки админки
│   │   └── middlewares.js  # Middleware
│   └── src/
│       └── index.js         # Точка входа
└── docker-compose.yaml     # Обновлен с Strapi
```

## 🐛 Решение проблем

### Strapi не запускается

1. Проверьте логи:
   ```bash
   docker-compose logs strapi
   ```

2. Проверьте подключение к БД:
   ```bash
   docker-compose exec strapi npm run strapi -- version
   ```

3. Пересоздайте контейнер:
   ```bash
   docker-compose down strapi
   docker-compose up -d strapi
   ```

### Ошибки подключения к БД

Убедитесь, что PostgreSQL запущен и доступен:
```bash
docker-compose ps postgres
```

## 📚 Дополнительные ресурсы

- [Документация Strapi](https://docs.strapi.io/)
- [Strapi API Reference](https://docs.strapi.io/dev-docs/api/rest)
- [Strapi Content Types](https://docs.strapi.io/dev-docs/backend-customization/models)

## ⚠️ Важные замечания

1. **Strapi создаст свои таблицы** в PostgreSQL с префиксом `strapi_`
2. **Ваши существующие таблицы** не будут затронуты
3. **Strapi и Go API** могут работать с одной БД, но лучше разделить данные
4. **Для production** настройте отдельные секреты для Strapi

