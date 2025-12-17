-- Добавление тестового маршрута "Маршрут к Эльбрусу"
-- Убедитесь что Strapi не запущен перед выполнением этого скрипта!

-- Добавляем маршрут
INSERT INTO routes (id, name, slug, description, is_active, created_at, updated_at, published_at)
VALUES (
  1,
  'Маршрут к Эльбрусу',
  'marshrut-k-elbrusu',
  'Увлекательный маршрут к подножию самой высокой горы Европы. Маршрут включает посещение живописных мест и знакомство с культурой местных народов.',
  1,
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00'
);

-- Связываем маршрут с типом "Пеший" (ID=1)
INSERT INTO routes_route_type_links (route_id, route_type_id)
VALUES (1, 1);

-- Связываем маршрут с местом "Гора Эльбрус" (ID=1)
INSERT INTO routes_places_links (route_id, place_id, place_order)
VALUES (1, 1, 1);

-- Проверка результата
SELECT 'Маршрут добавлен:' AS message;
SELECT id, name, slug, description FROM routes WHERE id = 1;

SELECT 'Тип маршрута:' AS message;
SELECT rt.name 
FROM route_types rt
JOIN routes_route_type_links rrtl ON rt.id = rrtl.route_type_id
WHERE rrtl.route_id = 1;

SELECT 'Места в маршруте:' AS message;
SELECT p.name 
FROM places p
JOIN routes_places_links rpl ON p.id = rpl.place_id
WHERE rpl.route_id = 1;
