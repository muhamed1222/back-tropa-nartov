#!/bin/bash

# Скрипт для применения миграции дефолтного рейтинга 3.5
# Использование: ./apply_default_rating.sh

set -e

echo "🔄 Применение миграции дефолтного рейтинга 3.5..."

# Загружаем переменные окружения
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Параметры подключения к БД
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-tropa_nartov}"

echo "📊 Подключение к базе данных: $DB_NAME на $DB_HOST:$DB_PORT"

# Применяем миграцию
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f internal/db/migrations/set_default_rating.sql

echo "✅ Миграция успешно применена!"
echo ""
echo "📈 Статистика:"
echo "   - Места с дефолтным рейтингом 3.5:"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM places WHERE rating = 3.5;"
echo "   - Маршруты с дефолтным рейтингом 3.5:"
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM routes WHERE rating = 3.5;"
echo ""
echo "🎉 Готово! Теперь все новые места и маршруты будут иметь рейтинг 3.5 по умолчанию."

