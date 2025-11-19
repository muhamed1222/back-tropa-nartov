#!/bin/bash

# Скрипт для запуска backend сервера
# Использование: ./start_server.sh

cd "$(dirname "$0")"

echo "🚀 Запуск backend сервера..."
echo ""

# Проверка порта
if lsof -ti:8001 > /dev/null 2>&1; then
    echo "⚠️  Порт 8001 уже занят. Останавливаю старый процесс..."
    kill -9 $(lsof -ti:8001) 2>/dev/null
    sleep 2
fi

# Проверка .env файла
if [ ! -f .env ]; then
    echo "⚠️  Файл .env не найден!"
    echo "Создайте файл .env с необходимыми переменными окружения"
    exit 1
fi

# Проверка go.mod
if [ ! -f go.mod ]; then
    echo "❌ Файл go.mod не найден. Убедитесь, что вы в правильной директории."
    exit 1
fi

echo "✅ Проверки пройдены"
echo "📡 Запуск сервера на http://localhost:8001"
echo ""
echo "Для остановки нажмите Ctrl+C"
echo ""

# Запуск сервера
go run ./cmd/api/main.go

