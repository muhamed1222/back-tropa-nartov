# 🧪 Тестовый аккаунт

**Дата создания:** 2025-11-19  
**Статус:** ✅ Активен

---

## 📋 Данные для входа

```
Email: test@test.com
Password: test12345
Имя: Тестовый
ID: 3
Роль: user
```

---

## 🔑 Получение токена

### Через API:

```bash
curl -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@test.com",
    "password": "test12345"
  }'
```

### Ответ:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 3,
    "email": "test@test.com",
    "name": "Тестовый",
    "first_name": "",
    "role": "user",
    "avatar_url": ""
  }
}
```

---

## 🚀 Использование

### Через скрипт:

```bash
cd /Users/kelemetovmuhamed/Documents/тропа\ нартов\ /back
./create_test_user.sh test@test.com test12345 "Тестовый Пользователь"
```

### Вручную через API:

```bash
# Получение токена
TOKEN=$(curl -s -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test12345"}' \
  | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Использование токена
curl -X GET http://localhost:8001/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

---

## ✅ Примеры использования

### Проверка профиля:

```bash
curl -X GET http://localhost:8001/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

### Получение избранных маршрутов:

```bash
curl -X GET http://localhost:8001/favorites/routes \
  -H "Authorization: Bearer $TOKEN"
```

### Добавление маршрута в избранное:

```bash
curl -X POST http://localhost:8001/favorites/routes/1 \
  -H "Authorization: Bearer $TOKEN"
```

### Проверка статуса избранного:

```bash
curl -X GET http://localhost:8001/favorites/routes/1/status \
  -H "Authorization: Bearer $TOKEN"
```

### Удаление из избранного:

```bash
curl -X DELETE http://localhost:8001/favorites/routes/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## 📝 Примечания

- Токен действителен 24 часа
- Аккаунт создан в базе данных с ID = 3
- Пароль хранится в хешированном виде
- Можно использовать для тестирования всех API эндпоинтов

---

**Создано:** 2025-11-19  
**Версия:** 1.0

