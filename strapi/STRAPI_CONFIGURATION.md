# 📋 Руководство по настройке Strapi для проекта "Тропа Нартов"

## ✅ Шаг 1: Создание первого администратора

1. Откройте админ-панель: **http://localhost:1337/admin**
2. Заполните форму:
   - **First name**: Ваше имя
   - **Last name**: Ваша фамилия
   - **Email**: Ваш email
   - **Password**: Надежный пароль
3. Нажмите **"Let's start"**

---

## 🗂️ Шаг 2: Создание Content Types

После создания администратора нужно создать Content Types, соответствующие вашим моделям данных.

### 1. Place (Место)

**Создание:**
- Content-Type Builder → Create new collection type → `Place`

**Поля:**
- `name` - **Text** (Short text) - Required
- `type` - **Text** (Short text)
- `description` - **Rich text** (Long text) - Required
- `overview` - **Rich text** (Long text)
- `history` - **Rich text** (Long text)
- `address` - **Text** (Short text) - Required
- `hours` - **Text** (Short text)
- `weekend` - **Text** (Short text)
- `entry` - **Text** (Short text)
- `contacts` - **Text** (Short text)
- `contacts_email` - **Email**
- `latitude` - **Number** (Decimal) - Required
- `longitude` - **Number** (Decimal) - Required
- `rating` - **Number** (Decimal, default: 0)
- `is_active` - **Boolean** (default: true)
- `type_id` - **Number** (Integer)
- `area_id` - **Number** (Integer)
- `opening_hours` - **Text** (Short text)

**Связи:**
- `images` - **Media** (Multiple media) - Relation to Image
- `reviews` - **Relation** (One-to-Many) - Place has many Reviews

**Настройки:**
- Display name: `Place`
- API ID: `place`
- Draft & Publish: Включить

---

### 2. Route (Маршрут)

**Создание:**
- Content-Type Builder → Create new collection type → `Route`

**Поля:**
- `name` - **Text** (Short text) - Required
- `description` - **Rich text** (Long text) - Required
- `overview` - **Rich text** (Long text)
- `history` - **Rich text** (Long text)
- `distance` - **Number** (Decimal) - Required
- `duration` - **Number** (Decimal)
- `rating` - **Number** (Decimal, default: 0)
- `is_active` - **Boolean** (default: true)
- `type_id` - **Number** (Integer) - Required
- `area_id` - **Number** (Integer) - Required

**Связи:**
- `type` - **Relation** (Many-to-One) - Route belongs to Type
- `area` - **Relation** (Many-to-One) - Route belongs to Area
- `categories` - **Relation** (Many-to-Many) - Route has and belongs to many Categories
- `reviews` - **Relation** (One-to-Many) - Route has many Reviews

**Настройки:**
- Display name: `Route`
- API ID: `route`
- Draft & Publish: Включить

---

### 3. Review (Отзыв)

**Создание:**
- Content-Type Builder → Create new collection type → `Review`

**Поля:**
- `text` - **Rich text** (Long text) - Required
- `rating` - **Number** (Integer, min: 1, max: 5) - Required
- `likes` - **Number** (Integer, default: 0)
- `is_active` - **Boolean** (default: true)

**Связи:**
- `user` - **Relation** (Many-to-One) - Review belongs to User
- `place` - **Relation** (Many-to-One, optional) - Review belongs to Place
- `route` - **Relation** (Many-to-One, optional) - Review belongs to Route

**Настройки:**
- Display name: `Review`
- API ID: `review`
- Draft & Publish: Включить

---

### 4. Category (Категория)

**Создание:**
- Content-Type Builder → Create new collection type → `Category`

**Поля:**
- `name` - **Text** (Short text) - Required, Unique
- `description` - **Rich text** (Long text)

**Связи:**
- `places` - **Relation** (Many-to-Many) - Category has and belongs to many Places
- `routes` - **Relation** (Many-to-Many) - Category has and belongs to many Routes

**Настройки:**
- Display name: `Category`
- API ID: `category`
- Draft & Publish: Включить

---

### 5. Tag (Тег)

**Создание:**
- Content-Type Builder → Create new collection type → `Tag`

**Поля:**
- `name` - **Text** (Short text) - Required, Unique
- `slug` - **UID** (from name)

**Связи:**
- `places` - **Relation** (Many-to-Many) - Tag has and belongs to many Places
- `routes` - **Relation** (Many-to-Many) - Tag has and belongs to many Routes

**Настройки:**
- Display name: `Tag`
- API ID: `tag`
- Draft & Publish: Включить

---

### 6. Type (Тип)

**Создание:**
- Content-Type Builder → Create new collection type → `Type`

**Поля:**
- `name` - **Text** (Short text) - Required, Unique
- `entity_type` - **Enumeration** (values: `place`, `route`) - Required
- `description` - **Rich text** (Long text)

**Связи:**
- `routes` - **Relation** (One-to-Many) - Type has many Routes

**Настройки:**
- Display name: `Type`
- API ID: `type`
- Draft & Publish: Включить

---

### 7. Area (Район)

**Создание:**
- Content-Type Builder → Create new collection type → `Area`

**Поля:**
- `name` - **Text** (Short text) - Required, Unique
- `description` - **Rich text** (Long text)

**Связи:**
- `places` - **Relation** (One-to-Many) - Area has many Places
- `routes` - **Relation** (One-to-Many) - Area has many Routes

**Настройки:**
- Display name: `Area`
- API ID: `area`
- Draft & Publish: Включить

---

### 8. Image (Изображение)

**Создание:**
- Content-Type Builder → Create new collection type → `Image`

**Поля:**
- `url` - **Text** (Short text) - Required
- `alt_text` - **Text** (Short text)

**Связи:**
- `place` - **Relation** (Many-to-One) - Image belongs to Place

**Настройки:**
- Display name: `Image`
- API ID: `image`
- Draft & Publish: Включить

---

## 🔐 Шаг 3: Настройка прав доступа (Permissions)

После создания Content Types нужно настроить права доступа:

1. **Settings** → **Users & Permissions plugin** → **Roles** → **Public**
2. Для каждого Content Type включите:
   - ✅ **find** (GET /api/places)
   - ✅ **findOne** (GET /api/places/:id)
   - ❌ **create** (только для авторизованных)
   - ❌ **update** (только для авторизованных)
   - ❌ **delete** (только для авторизованных)

3. **Settings** → **Users & Permissions plugin** → **Roles** → **Authenticated**
4. Для каждого Content Type включите:
   - ✅ **find**
   - ✅ **findOne**
   - ✅ **create** (создание отзывов)
   - ✅ **update** (обновление своих данных)
   - ✅ **delete** (удаление своих данных)

---

## 🔑 Шаг 4: Создание API Token

1. **Settings** → **API Tokens** → **Create new API Token**
2. Заполните:
   - **Name**: `Mobile App Token`
   - **Token duration**: Unlimited
   - **Token type**: Full access
3. Скопируйте токен (он показывается только один раз!)

---

## 📝 Шаг 5: Настройка медиа библиотеки

1. **Media Library** → Настройте загрузку файлов
2. Убедитесь, что включены форматы:
   - Images: JPG, PNG, WebP
   - Videos: MP4, WebM (если нужно)

---

## 🔄 Шаг 6: Интеграция с Go API (опционально)

Если вы хотите использовать Strapi параллельно с Go API:

### Вариант 1: Strapi только для админки
- Управление контентом через Strapi
- Чтение данных через Go API

### Вариант 2: Синхронизация данных
- Создайте webhook в Strapi для синхронизации
- Или используйте скрипт для миграции данных

### Вариант 3: Полная миграция
- Используйте Strapi API вместо Go API
- Настройте кастомные контроллеры в Strapi

---

## ✅ Чеклист настройки

- [ ] Создан первый администратор
- [ ] Создан Content Type `Place`
- [ ] Создан Content Type `Route`
- [ ] Создан Content Type `Review`
- [ ] Создан Content Type `Category`
- [ ] Создан Content Type `Tag`
- [ ] Создан Content Type `Type`
- [ ] Создан Content Type `Area`
- [ ] Создан Content Type `Image`
- [ ] Настроены права доступа для Public
- [ ] Настроены права доступа для Authenticated
- [ ] Создан API Token
- [ ] Протестирован доступ к API

---

## 🚀 Быстрый старт

1. Откройте: http://localhost:1337/admin
2. Создайте администратора
3. Создайте Content Types (см. выше)
4. Настройте права доступа
5. Создайте API Token
6. Начните добавлять контент!

---

## 📚 Полезные ссылки

- [Документация Strapi](https://docs.strapi.io/)
- [Content-Type Builder](https://docs.strapi.io/dev-docs/backend-customization/models)
- [API Reference](https://docs.strapi.io/dev-docs/api/rest)
- [Permissions](https://docs.strapi.io/dev-docs/plugins/users-permissions)

