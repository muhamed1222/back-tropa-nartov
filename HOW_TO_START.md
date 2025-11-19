# 🚀 Как запустить Backend сервер

## ✅ Быстрый запуск

### Вариант 1: Через скрипт (рекомендуется)

```bash
cd "/Users/kelemetovmuhamed/Documents/тропа нартов /back"
./start_server.sh
```

### Вариант 2: Вручную через Go

```bash
cd "/Users/kelemetovmuhamed/Documents/тропа нартов /back"
go run ./cmd/api/main.go
```

### Вариант 3: Сборка и запуск

```bash
cd "/Users/kelemetovmuhamed/Documents/тропа нартов /back"
go build -o bin/api ./cmd/api/main.go
./bin/api
```

---

## 📋 Предварительные требования

### 1. Проверка Go
```bash
go version
# Должно быть: go version go1.25.0 или выше
```

### 2. Проверка зависимостей
```bash
cd "/Users/kelemetovmuhamed/Documents/тропа нартов /back"
go mod download
```

### 3. Проверка .env файла
```bash
ls -la .env
# Файл должен существовать
```

### 4. Проверка базы данных
- Убедитесь, что PostgreSQL запущен
- Проверьте настройки в `.env` файле

---

## 🐳 Если используете Docker Compose

```bash
cd "/Users/kelemetovmuhamed/Documents/тропа нартов /back"
docker-compose up -d postgres
# Затем запустите Go сервер:
go run ./cmd/api/main.go
```

---

## ✅ Проверка, что сервер запустился

После запуска вы должны увидеть:

```
[GIN-debug] Listening and serving HTTP on :8001
```

Или:

```
✅ Server starting on port 8001
```

### Проверка работоспособности:

В другом терминале:
```bash
curl http://localhost:8001/ping
# Должен вернуть: {"message":"pong"}
```

---

## 🛑 Остановка сервера

- **В терминале:** Нажмите `Ctrl+C`
- **Если завис:** 
  ```bash
  lsof -ti:8001 | xargs kill -9
  ```

---

## 🔍 Решение проблем

### Проблема: Порт 8001 занят
```bash
# Найти процесс
lsof -ti:8001

# Убить процесс
kill -9 $(lsof -ti:8001)
```

### Проблема: База данных не доступна
- Проверьте, запущен ли PostgreSQL
- Проверьте настройки в `.env`
- Проверьте подключение:
  ```bash
  psql -h localhost -U postgres -d tropa_nartov
  ```

### Проблема: Ошибки компиляции
```bash
# Очистить кеш
go clean -cache

# Скачать зависимости заново
go mod download

# Попробовать снова
go run ./cmd/api/main.go
```

---

## 📝 Переменные окружения (.env)

Убедитесь, что в `.env` файле указаны:

```env
# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=tropa_nartov

# Server
APP_PORT=8001
HOST=localhost

# JWT
JWT_SECRET_KEY=your-secret-key

# CORS (опционально)
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

---

## 🚀 Автозапуск при старте системы

### macOS (через LaunchAgent):

1. Создайте файл `~/Library/LaunchAgents/com.tropanartov.api.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.tropanartov.api</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/go/bin/go</string>
        <string>run</string>
        <string>/Users/kelemetovmuhamed/Documents/тропа нартов /back/cmd/api/main.go</string>
    </array>
    <key>WorkingDirectory</key>
    <string>/Users/kelemetovmuhamed/Documents/тропа нартов /back</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
```

2. Загрузите:
```bash
launchctl load ~/Library/LaunchAgents/com.tropanartov.api.plist
```

---

**Готово!** Теперь вы можете запустить сервер одним из способов выше. 🎉

