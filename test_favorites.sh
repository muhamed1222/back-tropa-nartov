#!/bin/bash

# Скрипт для тестирования функциональности избранного
# Использование: ./test_favorites.sh

BASE_URL="http://localhost:8001"
TEST_EMAIL="test@test.com"
TEST_PASSWORD="test12345"

echo "🧪 Тестирование функциональности избранного"
echo "=========================================="
echo ""

# Цвета для вывода
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Функция для вывода успешного результата
success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# Функция для вывода ошибки
error() {
    echo -e "${RED}❌ $1${NC}"
}

# Функция для вывода предупреждения
warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Шаг 1: Получаем токен авторизации
echo "1️⃣  Получение токена авторизации..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$TEST_EMAIL\", \"password\": \"$TEST_PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    error "Не удалось получить токен авторизации"
    echo "Ответ сервера: $LOGIN_RESPONSE"
    exit 1
fi

success "Токен получен: ${TOKEN:0:20}..."

# Шаг 2: Получаем список мест для тестирования
echo ""
echo "2️⃣  Получение списка мест..."
PLACES_RESPONSE=$(curl -s -X GET "$BASE_URL/places" \
  -H "Content-Type: application/json")

PLACE_ID=$(echo "$PLACES_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$PLACE_ID" ]; then
    warning "Не найдено мест для тестирования. Используем ID=1"
    PLACE_ID=1
else
    success "Найдено место с ID=$PLACE_ID"
fi

# Шаг 3: Получаем список маршрутов для тестирования
echo ""
echo "3️⃣  Получение списка маршрутов..."
ROUTES_RESPONSE=$(curl -s -X GET "$BASE_URL/routes" \
  -H "Content-Type: application/json")

ROUTE_ID=$(echo "$ROUTES_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)

if [ -z "$ROUTE_ID" ]; then
    warning "Не найдено маршрутов для тестирования. Используем ID=1"
    ROUTE_ID=1
else
    success "Найден маршрут с ID=$ROUTE_ID"
fi

# ========== ТЕСТИРОВАНИЕ МЕСТ ==========
echo ""
echo "🏛️  ТЕСТИРОВАНИЕ ИЗБРАННЫХ МЕСТ"
echo "--------------------------------"

# Проверяем статус избранного места
echo ""
echo "4️⃣  Проверка статуса избранного места (ID=$PLACE_ID)..."
STATUS_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "$BASE_URL/favorites/places/$PLACE_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$STATUS_RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
BODY=$(echo "$STATUS_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    IS_FAVORITE=$(echo "$BODY" | grep -o '"is_favorite":[^,}]*' | cut -d':' -f2)
    success "Статус избранного получен: is_favorite=$IS_FAVORITE"
else
    error "Ошибка проверки статуса. HTTP: $HTTP_CODE"
    echo "Ответ: $BODY"
fi

# Добавляем место в избранное
echo ""
echo "5️⃣  Добавление места (ID=$PLACE_ID) в избранное..."
ADD_PLACE_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/favorites/places/$PLACE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$ADD_PLACE_RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
BODY=$(echo "$ADD_PLACE_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    success "Место добавлено в избранное"
else
    if [ "$HTTP_CODE" = "400" ]; then
        warning "Место уже в избранном (ожидаемое поведение)"
    else
        error "Ошибка добавления места. HTTP: $HTTP_CODE"
        echo "Ответ: $BODY"
    fi
fi

# Получаем список избранных мест
echo ""
echo "6️⃣  Получение списка избранных мест..."
FAV_PLACES_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "$BASE_URL/favorites/places" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$FAV_PLACES_RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
BODY=$(echo "$FAV_PLACES_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    PLACES_COUNT=$(echo "$BODY" | grep -o '"id"' | wc -l | tr -d ' ')
    success "Получен список избранных мест: $PLACES_COUNT мест"
    if [ "$PLACES_COUNT" -gt 0 ]; then
        echo "Первые 200 символов ответа:"
        echo "$BODY" | head -c 200
        echo "..."
    fi
else
    error "Ошибка получения списка мест. HTTP: $HTTP_CODE"
    echo "Ответ: $BODY"
fi

# ========== ТЕСТИРОВАНИЕ МАРШРУТОВ ==========
echo ""
echo "🛣️  ТЕСТИРОВАНИЕ ИЗБРАННЫХ МАРШРУТОВ"
echo "-------------------------------------"

# Проверяем статус избранного маршрута
echo ""
echo "7️⃣  Проверка статуса избранного маршрута (ID=$ROUTE_ID)..."
STATUS_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "$BASE_URL/favorites/routes/$ROUTE_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$STATUS_RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
BODY=$(echo "$STATUS_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    IS_FAVORITE=$(echo "$BODY" | grep -o '"is_favorite":[^,}]*' | cut -d':' -f2)
    success "Статус избранного получен: is_favorite=$IS_FAVORITE"
else
    error "Ошибка проверки статуса. HTTP: $HTTP_CODE"
    echo "Ответ: $BODY"
fi

# Добавляем маршрут в избранное
echo ""
echo "8️⃣  Добавление маршрута (ID=$ROUTE_ID) в избранное..."
ADD_ROUTE_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/favorites/routes/$ROUTE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$ADD_ROUTE_RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
BODY=$(echo "$ADD_ROUTE_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    success "Маршрут добавлен в избранное"
else
    if [ "$HTTP_CODE" = "400" ]; then
        warning "Маршрут уже в избранном (ожидаемое поведение)"
    else
        error "Ошибка добавления маршрута. HTTP: $HTTP_CODE"
        echo "Ответ: $BODY"
    fi
fi

# Получаем список избранных маршрутов
echo ""
echo "9️⃣  Получение списка избранных маршрутов..."
FAV_ROUTES_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "$BASE_URL/favorites/routes" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$FAV_ROUTES_RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
BODY=$(echo "$FAV_ROUTES_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    ROUTES_COUNT=$(echo "$BODY" | grep -o '"id"' | wc -l | tr -d ' ')
    success "Получен список избранных маршрутов: $ROUTES_COUNT маршрутов"
    if [ "$ROUTES_COUNT" -gt 0 ]; then
        echo "Первые 200 символов ответа:"
        echo "$BODY" | head -c 200
        echo "..."
        # Проверяем наличие обязательных полей
        if echo "$BODY" | grep -q '"name"'; then
            success "✅ Поле 'name' присутствует в ответе"
        else
            warning "⚠️  Поле 'name' отсутствует в ответе"
        fi
        if echo "$BODY" | grep -q '"description"'; then
            success "✅ Поле 'description' присутствует в ответе"
        else
            warning "⚠️  Поле 'description' отсутствует в ответе"
        fi
        if echo "$BODY" | grep -q '"type_name"'; then
            success "✅ Поле 'type_name' присутствует в ответе"
        else
            warning "⚠️  Поле 'type_name' отсутствует в ответе"
        fi
        if echo "$BODY" | grep -q '"distance"'; then
            success "✅ Поле 'distance' присутствует в ответе"
        else
            warning "⚠️  Поле 'distance' отсутствует в ответе"
        fi
    else
        warning "Список избранных маршрутов пуст"
    fi
else
    error "Ошибка получения списка маршрутов. HTTP: $HTTP_CODE"
    echo "Ответ: $BODY"
fi

# Удаляем маршрут из избранного (для проверки удаления)
echo ""
echo "🔟 Удаление маршрута (ID=$ROUTE_ID) из избранного..."
DELETE_ROUTE_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE "$BASE_URL/favorites/routes/$ROUTE_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$DELETE_ROUTE_RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
BODY=$(echo "$DELETE_ROUTE_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    success "Маршрут удален из избранного"
else
    error "Ошибка удаления маршрута. HTTP: $HTTP_CODE"
    echo "Ответ: $BODY"
fi

# Проверяем статус после удаления
echo ""
echo "1️⃣1️⃣  Проверка статуса после удаления..."
STATUS_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "$BASE_URL/favorites/routes/$ROUTE_ID/status" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$STATUS_RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
BODY=$(echo "$STATUS_RESPONSE" | sed '/HTTP_CODE/d')

if [ "$HTTP_CODE" = "200" ]; then
    IS_FAVORITE=$(echo "$BODY" | grep -o '"is_favorite":[^,}]*' | cut -d':' -f2)
    if [ "$IS_FAVORITE" = "false" ]; then
        success "Статус корректный: is_favorite=false (удален из избранного)"
    else
        warning "Статус: is_favorite=$IS_FAVORITE (ожидалось false)"
    fi
else
    error "Ошибка проверки статуса. HTTP: $HTTP_CODE"
fi

echo ""
echo "=========================================="
success "Тестирование завершено!"
echo ""

