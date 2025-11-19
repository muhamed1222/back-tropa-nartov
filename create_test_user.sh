#!/bin/bash

# Скрипт для создания тестового пользователя
# Использование: ./create_test_user.sh [email] [password] [name]

API_URL="${API_URL:-http://localhost:8001}"
EMAIL="${1:-test@test.com}"
PASSWORD="${2:-test12345}"
NAME="${3:-Тестовый Пользователь}"

echo "=== Создание тестового пользователя ==="
echo "Email: $EMAIL"
echo "Password: $PASSWORD"
echo "Name: $NAME"
echo ""

# Регистрация
echo "1. Регистрация..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"first_name\": \"$NAME\",
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

echo "$REGISTER_RESPONSE" | python3 -m json.tool 2>/dev/null || echo "$REGISTER_RESPONSE"
echo ""

# Попытка входа
echo "2. Вход..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL\",
    \"password\": \"$PASSWORD\"
  }")

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ ! -z "$TOKEN" ]; then
  echo "✅ Токен получен успешно!"
  echo "Токен: ${TOKEN:0:50}..."
  echo ""
  
  echo "3. Проверка профиля..."
  curl -s -X GET "$API_URL/auth/profile" \
    -H "Authorization: Bearer $TOKEN" | python3 -m json.tool 2>/dev/null || echo "Ошибка получения профиля"
  echo ""
  
  echo "=== Тестовый аккаунт создан успешно! ==="
  echo "Email: $EMAIL"
  echo "Password: $PASSWORD"
  echo ""
  echo "Для использования в тестах:"
  echo "export TEST_TOKEN=\"$TOKEN\""
  echo "curl -H \"Authorization: Bearer \$TEST_TOKEN\" $API_URL/favorites/routes"
else
  echo "❌ Не удалось получить токен"
  echo "Ответ сервера: $LOGIN_RESPONSE"
  exit 1
fi

