# 📝 Скрипты для работы с Strapi

## 🚀 Быстрый старт

### 1. Настройка прав доступа (обязательно!)

```bash
cd "/Users/kelemetovmuhamed/Documents/тропа нартов /back/strapi"
node scripts/setup-permissions.js
```

**Или настройте вручную:**
1. Откройте админ-панель: `http://localhost:1337/admin`
2. Settings → Users & Permissions → Roles → Public
3. Найдите "Place" в списке
4. Включите права: `find`, `findOne`, `create`
5. Нажмите "Save"

### 2. Добавление мест

```bash
node scripts/add-places.js
```

Скрипт добавит 6 тестовых мест:
- Парк Атажукинский
- Национальный музей КБР
- Гора Эльбрус
- Чегемские водопады
- Голубое озеро (Церик-Кель)
- Баксанское ущелье

---

## 📋 Файлы

- `setup-permissions.js` - настройка прав доступа для Public роли
- `add-places.js` - автоматическое добавление мест через API
- `add-places-manual.md` - инструкция для ручного добавления через админку

---

## ⚠️ Требования

1. **Strapi должен быть запущен:**
   ```bash
   npm run develop
   ```

2. **Права доступа должны быть настроены** (через скрипт или вручную)

3. **Установлен axios:**
   ```bash
   npm install axios
   ```

---

## 🔍 Проверка

После добавления мест проверьте через API:

```bash
curl http://localhost:1337/api/places?populate=*
```

Или в админ-панели:
- Content Manager → Место

