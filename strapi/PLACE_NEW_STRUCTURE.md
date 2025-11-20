# ✅ Новая структура Place - Оптимальная модель

## 🎯 Основные изменения

### ❌ Удалено из Place:

- `rating` - рейтинг теперь считается автоматически из Review
- `description` - заменено на `short_description` (текст для превью)
- `overview` - удалено (не используется)
- `type_name` - удалено (используется только relation `type`)
- `hours`, `opening_hours` - заменено на `working_hours`
- `weekend`, `entry` - удалены (можно добавить в `working_hours`)
- `contacts`, `contacts_email` - заменено на `contacts_phone`, `contacts_website`
- `type_id`, `area_id` - удалены (используются только relations)

### ✅ Добавлено/Изменено в Place:

- `short_description` - краткое описание для превью в списке (text, до 500 символов)
- `working_hours` - часы работы (text)
- `contacts_phone` - телефон (string, до 50 символов)
- `contacts_website` - веб-сайт (string, до 500 символов)
- `images` - теперь Media Library (multiple images) вместо relation
- `slug` - автоматически генерируется из названия

---

## 📋 Итоговая структура Place

### 🔹 1. Основные поля

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| **name** | string (200) | ✅ Да | Название места |
| **slug** | uid | ❌ Нет | URL-фрагмент (автогенерация из name) |
| **short_description** | text (500) | ❌ Нет | Краткое описание для превью в списке |
| **type** | relation (many-to-one) | ❌ Нет | Тип места (связь с Type) |
| **images** | media (multiple) | ❌ Нет | Изображения (Media Library) |

### 🔹 2. Контентные вкладки

#### Вкладка "История"

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| **history** | richtext | ❌ Нет | Историческая информация о месте |

#### Вкладка "Обзор"

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| **address** | string (500) | ❌ Нет | Адрес места |
| **working_hours** | text (200) | ❌ Нет | Часы работы |
| **contacts_phone** | string (50) | ❌ Нет | Телефон |
| **contacts_website** | string (500) | ❌ Нет | Веб-сайт |

### 🔹 3. Геоданные

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| **latitude** | decimal | ✅ Да | Широта (для карты) |
| **longitude** | decimal | ✅ Да | Долгота (для карты) |

### 🔹 4. Связи

| Поле | Тип | Описание |
|------|-----|----------|
| **reviews** | relation (one-to-many) | Отзывы о месте |
| **area** | relation (many-to-one) | Область/район |
| **categories** | relation (many-to-many) | Категории |
| **tags** | relation (many-to-many) | Теги |

### 🔹 5. Системные поля

| Поле | Тип | По умолчанию | Описание |
|------|-----|--------------|----------|
| **is_active** | boolean | `true` | Активно ли место |
| **createdAt** | datetime | авто | Дата создания |
| **updatedAt** | datetime | авто | Дата обновления |
| **publishedAt** | datetime | авто | Дата публикации |

---

## ⭐ Рейтинг - автоматический расчет

### ❌ Place НЕ хранит rating

### ✅ Рейтинг считается из Review

**Collection Type: Review**

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| **place** | relation (many-to-one) | ❌ Нет | Связь с Place |
| **route** | relation (many-to-one) | ❌ Нет | Связь с Route (альтернатива) |
| **user_name** | string (100) | ❌ Нет | Имя пользователя |
| **rating** | integer (1-5) | ✅ Да | Оценка (1-5 звезд) |
| **text** | text (2000) | ❌ Нет | Текст отзыва |
| **date** | datetime | ❌ Нет | Дата отзыва (автозаполнение) |
| **user_id** | integer | ❌ Нет | ID пользователя (опционально) |
| **likes** | integer | `0` | Количество лайков |

### Как считается рейтинг:

```javascript
// При запросе Place через API:
GET /api/places/:id?populate[reviews][fields][0]=rating

// Приложение считает среднее:
averageRating = reviews.reduce((sum, r) => sum + r.rating, 0) / reviews.length;

// Показываем: ⭐ 4.8
```

---

## 📊 Сравнение: ДО и ПОСЛЕ

### ❌ Было (старая структура):

```
Place
├── name
├── description (richtext) ← длинный текст
├── overview (richtext)
├── type_name (string)
├── hours, opening_hours, weekend, entry
├── contacts, contacts_email
├── rating (decimal) ← вручную вводили
├── type_id, area_id ← дублирование
└── images (relation) ← через Image модель
```

### ✅ Стало (новая структура):

```
Place
├── name
├── slug (auto)
├── short_description (text) ← краткий текст для превью
├── type (relation) ← выбор из списка
├── history (richtext)
├── address
├── working_hours
├── contacts_phone
├── contacts_website
├── latitude, longitude
├── images (media multiple) ← Media Library
├── reviews (relation) ← для расчета рейтинга
├── area, categories, tags (relations)
└── is_active
```

---

## 🎯 Преимущества новой структуры

### 1. ✅ Упрощение

- Меньше полей
- Четкая структура (вкладки)
- Нет дублирования

### 2. ✅ Автоматический рейтинг

- Рейтинг считается из отзывов
- Нет необходимости вручную вводить
- Всегда актуальный

### 3. ✅ Удобство редактирования

- Media Library для изображений
- Выбор из списков (type, area)
- Краткие поля для контактов

### 4. ✅ Гибкость

- Можно добавить больше полей в будущем
- Связи легко расширяются
- Структура понятна

---

## 📝 Инструкция по миграции

### Если есть существующие данные:

1. **Рейтинг:**
   - Удалите поле `rating` из существующих Place
   - Рейтинг будет считаться из Review автоматически

2. **Описание:**
   - `description` → `short_description` (обрезать до 500 символов)

3. **Контакты:**
   - `contacts` → `contacts_phone`
   - `contacts_email` → `contacts_website` (если это сайт)

4. **Часы работы:**
   - `hours` / `opening_hours` → `working_hours`

5. **Изображения:**
   - Если использовали relation Image → перенести в Media Library

---

## ✅ Итог

**Place теперь:**
- ✅ Проще и понятнее
- ✅ Без дублирования
- ✅ С автоматическим рейтингом
- ✅ Удобнее редактировать
- ✅ Соответствует UI приложения

**Рейтинг:**
- ✅ Не хранится в Place
- ✅ Считается автоматически из Review
- ✅ Всегда актуален
- ✅ Меньше ошибок

---

## 📚 Документация

- `HOW_TO_ADD_CONTENT.md` - как добавлять контент
- `PLACE_SCHEMA_IMPROVEMENTS.md` - предыдущие улучшения
- `PLACE_NEW_STRUCTURE.md` - это руководство

