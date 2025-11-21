# JWT Secret Configuration

## Новый безопасный JWT секрет сгенерирован

Добавьте следующую строку в ваш `.env` файл:

```env
JWT_SECRET_KEY=ylt5FHf/kzQiXoZaL7jaKu37QwMGAH9657OunbloWy1jrj7l6PWp9BbpLUboQEtU
```

## Важно!

1. ✅ `.env` файл добавлен в `.gitignore`
2. ✅ Fallback секреты удалены из кода
3. ⚠️ Приложение теперь **требует** установки `JWT_SECRET_KEY` в `.env`
4. 🔒 Никогда не коммитьте `.env` файл в git!

## Применение изменений

После добавления секрета в `.env`, перезапустите Go API:

```bash
cd /Users/kelemetovmuhamed/Documents/тропа нартов /back
# Остановите текущий процесс (Ctrl+C)
go run ./cmd/api/main.go
```

## Проверка

API больше не будет запускаться без установленного `JWT_SECRET_KEY`.
Это гарантирует безопасность в production.

