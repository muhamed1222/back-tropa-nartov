# 🔐 Учетные данные администратора Strapi

**Дата создания:** 21 ноября 2025, 08:30 MSK

---

## 📋 Данные для входа

```
Email:    admin@example.com
Password: Admin123!
```

---

## 🌐 Доступ

- **Админ-панель:** http://localhost:1337/admin
- **API:** http://localhost:1337/api

---

## ⚠️ Важно

- При первом входе рекомендуется сменить пароль
- Эти данные только для локальной разработки
- Для production используйте более сложный пароль

---

## 🔄 Если забыли пароль

Создайте нового администратора через Docker:

```bash
docker exec tropa-strapi-dev sh -c 'node /opt/app/node_modules/@strapi/strapi/bin/strapi.js admin:create-user \
  --email new-admin@example.com \
  --password NewPassword123! \
  --firstname Admin \
  --lastname User'
```

---

## 📚 Следующие шаги

1. ✅ Войдите в админ-панель
2. ✅ Смените пароль (Settings → Administration Panel → Users)
3. ✅ Проверьте Content-Type Builder
4. ✅ Создайте тестовый маршрут

---

*Создано: 21.11.2025, 08:30 MSK*

