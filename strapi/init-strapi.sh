#!/bin/bash

# Скрипт для инициализации Strapi при первом запуске

echo "🚀 Инициализация Strapi..."

# Проверяем, существует ли уже приложение
if [ ! -f "package.json" ]; then
    echo "📦 Создание нового Strapi приложения..."
    npx create-strapi-app@latest . --quickstart --no-run --dbclient=postgres --dbhost=postgres --dbport=5432 --dbname=tropa_nartov --dbusername=postgres --dbpassword=postgres
else
    echo "✅ Strapi приложение уже существует"
fi

echo "✅ Инициализация завершена!"

