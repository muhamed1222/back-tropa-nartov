FROM node:20-alpine

# Устанавливаем рабочую директорию
WORKDIR /opt/app

# Устанавливаем зависимости для компиляции нативных модулей
RUN apk add --no-cache python3 make g++

# Копируем package.json
COPY package.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем всю структуру проекта
COPY . .

# Создаем необходимые директории если их нет
RUN mkdir -p src/api src/components src/extensions src/policies src/middlewares public/uploads

# Устанавливаем права
RUN chown -R node:node /opt/app

# Переключаемся на пользователя node
USER node

# Открываем порт
EXPOSE 1337

# Запускаем Strapi в режиме разработки
CMD ["npm", "run", "develop"]

