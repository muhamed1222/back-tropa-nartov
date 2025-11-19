# 🚀 Быстрый старт Strapi

## Шаг 1: Откройте админ-панель

```
http://localhost:1337/admin
```

## Шаг 2: Создайте администратора

Заполните форму регистрации первого администратора.

## Шаг 3: Создайте Content Types

Перейдите в **Content-Type Builder** и создайте следующие типы:

### Минимальный набор (для начала):

1. **Place** - Места
   - name (Text, required)
   - description (Rich text, required)
   - address (Text, required)
   - latitude (Number, decimal)
   - longitude (Number, decimal)
   - rating (Number, decimal)

2. **Route** - Маршруты
   - name (Text, required)
   - description (Rich text, required)
   - distance (Number, decimal)
   - duration (Number, decimal)
   - rating (Number, decimal)

3. **Review** - Отзывы
   - text (Rich text, required)
   - rating (Number, integer, 1-5)
   - place (Relation, optional)
   - route (Relation, optional)

## Шаг 4: Настройте права доступа

**Settings** → **Users & Permissions** → **Roles** → **Public**

Включите для всех Content Types:
- ✅ find
- ✅ findOne

## Шаг 5: Создайте API Token

**Settings** → **API Tokens** → **Create new API Token**

Используйте этот токен для доступа к API из вашего приложения.

## Готово! 🎉

Теперь вы можете:
- Добавлять контент через админ-панель
- Получать данные через API: `http://localhost:1337/api/places`
- Использовать Strapi как админку для управления контентом

