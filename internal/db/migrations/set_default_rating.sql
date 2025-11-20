-- Миграция для установки дефолтного рейтинга 3.5
-- Дата: 2025-11-19

-- 1. Изменяем тип колонки rating в таблице routes с decimal(2,1) на decimal(3,1)
ALTER TABLE routes ALTER COLUMN rating TYPE DECIMAL(3,1);

-- 2. Изменяем тип колонки rating в таблице places на decimal(3,1) для единообразия
ALTER TABLE places ALTER COLUMN rating TYPE DECIMAL(3,1);

-- 3. Устанавливаем новое значение по умолчанию для places
ALTER TABLE places ALTER COLUMN rating SET DEFAULT 3.5;

-- 4. Устанавливаем новое значение по умолчанию для routes
ALTER TABLE routes ALTER COLUMN rating SET DEFAULT 3.5;

-- 5. Обновляем существующие записи places с рейтингом 0 на 3.5
-- (только те, у которых нет отзывов)
UPDATE places
SET rating = 3.5
WHERE rating = 0 
  AND id NOT IN (
    SELECT DISTINCT place_id 
    FROM reviews 
    WHERE place_id IS NOT NULL AND is_active = true
  );

-- 6. Обновляем существующие записи routes с рейтингом 0 на 3.5
-- (только те, у которых нет отзывов)
UPDATE routes
SET rating = 3.5
WHERE rating = 0 
  AND id NOT IN (
    SELECT DISTINCT route_id 
    FROM reviews 
    WHERE route_id IS NOT NULL AND is_active = true
  );

-- Проверка результатов
-- SELECT COUNT(*) as places_with_default_rating FROM places WHERE rating = 3.5;
-- SELECT COUNT(*) as routes_with_default_rating FROM routes WHERE rating = 3.5;

