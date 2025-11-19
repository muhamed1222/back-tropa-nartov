# 🚀 Статус Backend API

**Дата проверки:** 2025-11-19  
**Статус:** ✅ **ЗАПУЩЕН И РАБОТАЕТ**

---

## 📡 Основная информация

- **URL:** http://localhost:8001
- **Порт:** 8001
- **Статус:** ✅ Активен

---

## ✅ Доступные эндпоинты

### 🔐 Авторизация
- `POST /auth/register` - Регистрация
- `POST /auth/login` - Вход
- `GET /auth/profile` - Профиль (требует авторизации)
- `PUT /auth/profile` - Обновление профиля (требует авторизации)
- `POST /auth/forgot-password` - Восстановление пароля
- `POST /auth/verify-reset-code` - Проверка кода восстановления
- `POST /auth/reset-password` - Сброс пароля
- `DELETE /auth/delete-account` - Удаление аккаунта (требует авторизации)
- `POST /auth/upload-avatar` - Загрузка аватара (требует авторизации)
- `DELETE /auth/delete-avatar` - Удаление аватара (требует авторизации)
- `PUT /auth/change-password` - Смена пароля (требует авторизации)

### 📍 Места
- `GET /places` - Список мест
- `GET /places/:id` - Детали места
- `POST /places` - Создание места (требует авторизации)
- `PUT /places/:id` - Обновление места (требует авторизации)
- `DELETE /places/:id` - Удаление места (требует авторизации)

### 🗺️ Маршруты
- `GET /routes` - Список маршрутов
- `GET /routes/:id` - Детали маршрута
- `GET /routes/debug/all` - Отладочный эндпоинт

### ⭐ Избранное - Места
- `GET /favorites/places` - Получить избранные места (требует авторизации)
- `POST /favorites/places/:placeId` - Добавить место в избранное (требует авторизации)
- `DELETE /favorites/places/:placeId` - Удалить место из избранного (требует авторизации)
- `GET /favorites/places/:placeId/status` - Проверить статус избранного (требует авторизации)

### ⭐ Избранное - Маршруты
- `GET /favorites/routes` - Получить избранные маршруты (требует авторизации) ✅
- `POST /favorites/routes/:routeId` - Добавить маршрут в избранное (требует авторизации) ✅
- `DELETE /favorites/routes/:routeId` - Удалить маршрут из избранного (требует авторизации) ✅
- `GET /favorites/routes/:routeId/status` - Проверить статус избранного (требует авторизации) ✅

### 📝 Отзывы
- `GET /reviews/place/:placeId` - Отзывы места
- `POST /reviews` - Создание отзыва (требует авторизации)

### 📊 Активность пользователя
- `GET /user/activity` - Общая активность (требует авторизации)
- `GET /user/activity/places` - Посещенные места (требует авторизации)
- `GET /user/activity/routes` - Пройденные маршруты (требует авторизации)
- `POST /user/activity/places/:placeId` - Отметить место как посещенное (требует авторизации)
- `POST /user/activity/routes/:routeId` - Отметить маршрут как пройденный (требует авторизации)
- `DELETE /user/activity/places/:placeId` - Убрать из посещенных (требует авторизации)
- `DELETE /user/activity/routes/:routeId` - Убрать из пройденных (требует авторизации)

### 🧪 Утилиты
- `GET /ping` - Проверка работоспособности

---

## 🧪 Тестовый аккаунт

```
Email: test@test.com
Password: test12345
```

**Использование:**
```bash
# Получение токена
TOKEN=$(curl -s -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test12345"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Использование токена
curl -X GET http://localhost:8001/favorites/routes \
  -H "Authorization: Bearer $TOKEN"
```

---

## ✅ Проверка работоспособности

### Быстрая проверка:
```bash
curl http://localhost:8001/ping
# Ожидаемый ответ: {"message":"pong"}
```

### Проверка с авторизацией:
```bash
# Получить токен
TOKEN=$(curl -s -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test12345"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Проверить эндпоинт
curl -X GET http://localhost:8001/favorites/routes \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📋 Статус функционала

### ✅ Реализовано и протестировано:
- [x] Авторизация (регистрация, вход, восстановление пароля)
- [x] Управление профилем
- [x] Места (CRUD операции)
- [x] Маршруты (чтение)
- [x] Избранное для мест
- [x] Избранное для маршрутов ✅ (новое)
- [x] Отзывы
- [x] Активность пользователя

### 🔄 В разработке:
- [ ] Полный CRUD для маршрутов
- [ ] Загрузка изображений
- [ ] Продвинутые фильтры

---

## 🚀 Команды для запуска

### Запуск сервера:
```bash
cd /Users/kelemetovmuhamed/Documents/тропа\ нартов\ /back
go run ./cmd/api/main.go
```

### Остановка сервера:
```bash
pkill -f "go run.*main.go"
# или
lsof -ti:8001 | xargs kill -9
```

---

**Последнее обновление:** 2025-11-19  
**Версия:** 1.0

