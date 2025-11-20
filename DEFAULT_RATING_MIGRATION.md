# Миграция: Дефолтный рейтинг 3.5

## 📋 Описание

Эта миграция устанавливает дефолтный рейтинг **3.5** для всех мест и маршрутов.

### Что изменилось:

**До:**
- Новые места и маршруты создавались с рейтингом `0.0`
- Рейтинг обновлялся только после получения отзывов

**После:**
- Новые места и маршруты создаются с рейтингом `3.5`
- При добавлении отзывов рейтинг пересчитывается по формуле `AVG(rating)`
- Рейтинг может как расти (при хороших отзывах), так и падать (при плохих)

## 🔧 Технические изменения

### 1. Модели Go (backend)

**`internal/models/place.go`:**
```go
Rating float32 `gorm:"type:decimal(3,1);default:3.5" json:"rating"`
```

**`internal/models/route.go`:**
```go
Rating float32 `gorm:"type:decimal(3,1);default:3.5"`
```

### 2. База данных PostgreSQL

- Изменен тип колонки `routes.rating` с `decimal(2,1)` на `decimal(3,1)`
- Установлено значение по умолчанию `3.5` для обеих таблиц
- Обновлены существующие записи с рейтингом `0` на `3.5` (только те, у которых нет отзывов)

## 🚀 Применение миграции

### Автоматический способ (рекомендуется):

```bash
cd /Users/kelemetovmuhamed/Documents/тропа\ нартов\ /back
./apply_default_rating.sh
```

### Ручной способ:

```bash
# 1. Убедитесь, что PostgreSQL запущен
docker-compose up -d postgres

# 2. Примените SQL миграцию
psql -h localhost -p 5432 -U postgres -d tropa_nartov -f internal/db/migrations/set_default_rating.sql

# 3. Проверьте результаты
psql -h localhost -p 5432 -U postgres -d tropa_nartov -c "SELECT COUNT(*) FROM places WHERE rating = 3.5;"
psql -h localhost -p 5432 -U postgres -d tropa_nartov -c "SELECT COUNT(*) FROM routes WHERE rating = 3.5;"
```

## 📊 Логика работы рейтинга

### Новые записи:
```
Место/Маршрут создан → rating = 3.5 (по умолчанию)
```

### После получения отзывов:
```
Отзывы: 5, 4, 5, 3
Рейтинг = (5 + 4 + 5 + 3) / 4 = 4.25 ≈ 4.3
```

### Примеры:

| Отзывы | Средняя оценка | Результат |
|--------|----------------|-----------|
| Нет отзывов | - | **3.5** (дефолт) |
| 5, 5, 5 | 5.0 | **5.0** ⬆️ |
| 4, 4, 5 | 4.3 | **4.3** ⬆️ |
| 3, 3, 3 | 3.0 | **3.0** ⬇️ |
| 2, 2, 3 | 2.3 | **2.3** ⬇️ |
| 1, 1, 2 | 1.3 | **1.3** ⬇️ |

## ⚠️ Важно

1. **Миграция безопасна** - она обновляет только записи с рейтингом `0` и без отзывов
2. **Места/маршруты с отзывами не затронуты** - их рейтинг остается рассчитанным по формуле
3. **Автоматическое обновление** - триггеры БД продолжают работать как прежде
4. **Новые записи** - будут автоматически создаваться с рейтингом `3.5`

## 🧪 Тестирование

После применения миграции:

1. **Перезапустите backend:**
   ```bash
   cd /Users/kelemetovmuhamed/Documents/тропа\ нартов\ /back
   go run cmd/api/main.go
   ```

2. **Проверьте в приложении:**
   - Откройте экран "Места" или "Маршруты"
   - Убедитесь, что места без отзывов показывают рейтинг `3.5`
   - Добавьте отзыв и проверьте, что рейтинг пересчитался

3. **Проверьте в базе данных:**
   ```sql
   -- Места с дефолтным рейтингом
   SELECT id, name, rating FROM places WHERE rating = 3.5 LIMIT 10;
   
   -- Маршруты с дефолтным рейтингом
   SELECT id, name, rating FROM routes WHERE rating = 3.5 LIMIT 10;
   
   -- Места с отзывами (рейтинг должен быть рассчитан)
   SELECT p.id, p.name, p.rating, COUNT(r.id) as reviews_count
   FROM places p
   LEFT JOIN reviews r ON r.place_id = p.id AND r.is_active = true
   GROUP BY p.id
   HAVING COUNT(r.id) > 0
   LIMIT 10;
   ```

## 📝 Откат миграции (если нужно)

Если по какой-то причине нужно вернуть рейтинг `0` по умолчанию:

```sql
-- Вернуть дефолт на 0
ALTER TABLE places ALTER COLUMN rating SET DEFAULT 0;
ALTER TABLE routes ALTER COLUMN rating SET DEFAULT 0;

-- Обновить записи обратно на 0 (опционально)
UPDATE places SET rating = 0 WHERE rating = 3.5 AND id NOT IN (
    SELECT DISTINCT place_id FROM reviews WHERE place_id IS NOT NULL AND is_active = true
);
UPDATE routes SET rating = 0 WHERE rating = 3.5 AND id NOT IN (
    SELECT DISTINCT route_id FROM reviews WHERE route_id IS NOT NULL AND is_active = true
);
```

## ✅ Готово!

После применения миграции все новые места и маршруты будут иметь рейтинг **3.5** по умолчанию, который будет изменяться в зависимости от отзывов пользователей.

