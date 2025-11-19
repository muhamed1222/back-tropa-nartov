# 🔧 Решение проблем с авторизацией

## ❌ Ошибка 401 при логине

### Возможные причины:

1. **Неправильный email или пароль**
   - Проверьте правильность введенных данных
   - Убедитесь, что используете правильный email и пароль

2. **Rate Limiting (защита от брутфорса)**
   - Если слишком много неудачных попыток входа
   - Подождите несколько минут перед следующей попыткой
   - Или очистите rate limiter (перезапустите сервер)

3. **Пользователь неактивен**
   - Убедитесь, что `is_active = true` в базе данных
   - Проверьте через SQL или API

4. **Проблемы с хешированием пароля**
   - Пароль должен быть правильно захеширован при регистрации
   - Проверьте, что используется bcrypt

### Решение:

#### 1. Проверка тестового аккаунта:

```bash
curl -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@test.com",
    "password": "test12345"
  }'
```

#### 2. Создание нового аккаунта:

```bash
curl -X POST http://localhost:8001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Новый Пользователь",
    "email": "newuser@test.com",
    "password": "password123"
  }'
```

#### 3. Проверка в базе данных:

```sql
-- Проверить пользователя
SELECT id, email, is_active, password_hash FROM users WHERE email = 'test@test.com';

-- Проверить, активен ли пользователь
UPDATE users SET is_active = true WHERE email = 'test@test.com';
```

### Улучшенная обработка ошибок:

Код теперь возвращает более понятные сообщения:
- "Email обязателен для заполнения"
- "Неверный формат email адреса"
- "Пароль должен содержать минимум 8 символов"
- "Пользователь не найден"
- "Неверный пароль"

### Логирование:

Сервер теперь логирует:
- `🔐 [LOGIN] Попытка входа: email=...`
- `✅ [LOGIN] Успешный вход: id=..., email=...`
- `❌ [LOGIN] Ошибка входа: ...`

---

## 🔍 Диагностика проблем

### Проверка работоспособности API:

```bash
# 1. Проверка сервера
curl http://localhost:8001/ping

# 2. Попытка логина
curl -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test12345"}'

# 3. Проверка с неправильным паролем (ожидаем 401)
curl -X POST http://localhost:8001/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"wrongpassword"}'
```

### Проверка rate limiting:

Если получаете 429 (Too Many Requests):
- Подождите 1-2 минуты
- Или перезапустите сервер для сброса лимитов

---

**Последнее обновление:** 2025-11-19

