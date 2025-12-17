-- Добавление маршрутов по достопримечательностям Нальчика в Strapi
-- ВНИМАНИЕ: Остановите Strapi перед выполнением этого скрипта!
-- Использование: sqlite3 .tmp/data.db < scripts/add-nalchik-routes.sql
-- Убедитесь, что сначала выполнены скрипты add-nalchik-places.sql

-- 1. Пешеходный маршрут "Центр Нальчика"
INSERT INTO routes (name, slug, description, is_active, rating, created_at, updated_at, published_at)
VALUES (
  'Пешеходный маршрут "Центр Нальчика"',
  'peshehodnyj-marshrut-centr-nalchika',
  'Увлекательная пешая прогулка по центру Нальчика. Маршрут включает посещение главных парков города, площадей и памятников. Идеально подходит для первого знакомства с городом.',
  1,
  4.5,
  datetime('now'),
  datetime('now'),
  datetime('now')
);

-- Связываем с типом "Пеший" (ID=1)
INSERT INTO routes_route_type_links (route_id, route_type_id)
SELECT last_insert_rowid(), 1;

-- Добавляем места в маршрут
INSERT INTO routes_places_links (route_id, place_id, place_order)
SELECT 
  (SELECT last_insert_rowid()),
  id,
  CASE 
    WHEN name = 'Курортный парк "Долина нарзанов"' THEN 1
    WHEN name = 'Площадь Абхазии' THEN 2
    WHEN name = 'Сквер имени Лермонтова' THEN 3
    WHEN name = 'Памятник "Вечная слава"' THEN 4
    WHEN name = 'Атажукинский сад' THEN 5
  END
FROM places
WHERE name IN (
  'Курортный парк "Долина нарзанов"',
  'Площадь Абхазии',
  'Сквер имени Лермонтова',
  'Памятник "Вечная слава"',
  'Атажукинский сад'
);

-- 2. Культурный маршрут "Искусство и культура Нальчика"
INSERT INTO routes (name, slug, description, is_active, rating, created_at, updated_at, published_at)
VALUES (
  'Культурный маршрут "Искусство и культура Нальчика"',
  'kulturnyj-marshrut-iskusstvo-i-kultura-nalchika',
  'Маршрут для любителей искусства и культуры. Посетите главные культурные достопримечательности Нальчика: театр, музеи, мечеть. Погрузитесь в культурную жизнь столицы Кабардино-Балкарии.',
  1,
  4.6,
  datetime('now'),
  datetime('now'),
  datetime('now')
);

-- Связываем с типом "Авто" (ID=2)
INSERT INTO routes_route_type_links (route_id, route_type_id)
SELECT last_insert_rowid(), 2;

-- Добавляем места в маршрут
INSERT INTO routes_places_links (route_id, place_id, place_order)
SELECT 
  (SELECT last_insert_rowid()),
  id,
  CASE 
    WHEN name = 'Кабардинский государственный драматический театр им. А. Шогенцукова' THEN 1
    WHEN name = 'Музей изобразительных искусств Кабардино-Балкарской Республики' THEN 2
    WHEN name = 'Соборная мечеть Нальчика' THEN 3
  END
FROM places
WHERE name IN (
  'Кабардинский государственный драматический театр им. А. Шогенцукова',
  'Музей изобразительных искусств Кабардино-Балкарской Республики',
  'Соборная мечеть Нальчика'
);

-- 3. Исторический маршрут "Память поколений"
INSERT INTO routes (name, slug, description, is_active, rating, created_at, updated_at, published_at)
VALUES (
  'Исторический маршрут "Память поколений"',
  'istoricheskij-marshrut-pamyat-pokolenij',
  'Маршрут по историческим и мемориальным местам Нальчика. Посетите памятники, посвященные истории города и героям Великой Отечественной войны.',
  1,
  4.7,
  datetime('now'),
  datetime('now'),
  datetime('now')
);

-- Связываем с типом "Пеший" (ID=1)
INSERT INTO routes_route_type_links (route_id, route_type_id)
SELECT last_insert_rowid(), 1;

-- Добавляем места в маршрут
INSERT INTO routes_places_links (route_id, place_id, place_order)
SELECT 
  (SELECT last_insert_rowid()),
  id,
  CASE 
    WHEN name = 'Памятник "Вечная слава"' THEN 1
    WHEN name = 'Памятник Ленину' THEN 2
    WHEN name = 'Сквер имени Лермонтова' THEN 3
  END
FROM places
WHERE name IN (
  'Памятник "Вечная слава"',
  'Памятник Ленину',
  'Сквер имени Лермонтова'
);

-- 4. Парковый маршрут "Зеленые легкие Нальчика"
INSERT INTO routes (name, slug, description, is_active, rating, created_at, updated_at, published_at)
VALUES (
  'Парковый маршрут "Зеленые легкие Нальчика"',
  'parkovyj-marshrut-zelenye-legkie-nalchika',
  'Маршрут по всем паркам и скверам Нальчика. Идеально для любителей природы и спокойного отдыха. Прогуляйтесь по самым красивым зеленым зонам города.',
  1,
  4.4,
  datetime('now'),
  datetime('now'),
  datetime('now')
);

-- Связываем с типом "Пеший" (ID=1)
INSERT INTO routes_route_type_links (route_id, route_type_id)
SELECT last_insert_rowid(), 1;

-- Добавляем места в маршрут
INSERT INTO routes_places_links (route_id, place_id, place_order)
SELECT 
  (SELECT last_insert_rowid()),
  id,
  CASE 
    WHEN name = 'Курортный парк "Долина нарзанов"' THEN 1
    WHEN name = 'Атажукинский сад' THEN 2
    WHEN name = 'Сквер имени Лермонтова' THEN 3
    WHEN name = 'Городской парк культуры и отдыха' THEN 4
  END
FROM places
WHERE name IN (
  'Курортный парк "Долина нарзанов"',
  'Атажукинский сад',
  'Сквер имени Лермонтова',
  'Городской парк культуры и отдыха'
);

-- 5. Комплексный маршрут "Знакомство с Нальчиком"
INSERT INTO routes (name, slug, description, is_active, rating, created_at, updated_at, published_at)
VALUES (
  'Комплексный маршрут "Знакомство с Нальчиком"',
  'kompleksnyj-marshrut-znakomstvo-s-nalchikom',
  'Полный маршрут для знакомства с Нальчиком за один день. Включает посещение главных достопримечательностей: парков, культурных объектов, памятников и площадей. Рекомендуется для туристов, впервые посещающих город.',
  1,
  4.6,
  datetime('now'),
  datetime('now'),
  datetime('now')
);

-- Связываем с типом "Авто" (ID=2)
INSERT INTO routes_route_type_links (route_id, route_type_id)
SELECT last_insert_rowid(), 2;

-- Добавляем места в маршрут (топ-7 достопримечательностей)
INSERT INTO routes_places_links (route_id, place_id, place_order)
SELECT 
  (SELECT last_insert_rowid()),
  id,
  CASE 
    WHEN name = 'Курортный парк "Долина нарзанов"' THEN 1
    WHEN name = 'Соборная мечеть Нальчика' THEN 2
    WHEN name = 'Кабардинский государственный драматический театр им. А. Шогенцукова' THEN 3
    WHEN name = 'Площадь Абхазии' THEN 4
    WHEN name = 'Памятник "Вечная слава"' THEN 5
    WHEN name = 'Музей изобразительных искусств Кабардино-Балкарской Республики' THEN 6
    WHEN name = 'Атажукинский сад' THEN 7
  END
FROM places
WHERE name IN (
  'Курортный парк "Долина нарзанов"',
  'Соборная мечеть Нальчика',
  'Кабардинский государственный драматический театр им. А. Шогенцукова',
  'Площадь Абхазии',
  'Памятник "Вечная слава"',
  'Музей изобразительных искусств Кабардино-Балкарской Республики',
  'Атажукинский сад'
);

-- Проверка результатов
SELECT '═══════════════════════════════════════' AS separator;
SELECT 'Добавлено маршрутов по Нальчику:' AS message;
SELECT COUNT(*) AS total 
FROM routes 
WHERE slug LIKE '%nalchik%' OR description LIKE '%Нальчик%';

SELECT '═══════════════════════════════════════' AS separator;
SELECT 'Список добавленных маршрутов:' AS message;
SELECT id, name, slug, rating 
FROM routes 
WHERE slug LIKE '%nalchik%' OR description LIKE '%Нальчик%'
ORDER BY id;

SELECT '═══════════════════════════════════════' AS separator;
SELECT 'Места в маршрутах:' AS message;
SELECT 
  r.name AS route_name,
  p.name AS place_name,
  rpl.place_order
FROM routes r
JOIN routes_places_links rpl ON r.id = rpl.route_id
JOIN places p ON rpl.place_id = p.id
WHERE r.slug LIKE '%nalchik%' OR r.description LIKE '%Нальчик%'
ORDER BY r.id, rpl.place_order;

