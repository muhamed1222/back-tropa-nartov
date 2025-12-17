-- Добавление дополнительных маршрутов
-- Останавливайте Strapi перед выполнением!

-- 2. Маршрут по Чегемскому ущелью
INSERT INTO routes (id, name, slug, description, is_active, created_at, updated_at, published_at)
VALUES (
  2,
  'Маршрут по Чегемскому ущелью',
  'marshrut-po-chegemskomu-usheliyu',
  'Живописный маршрут вдоль Чегемского ущелья к знаменитым Чегемским водопадам. По пути открываются захватывающие виды на ущелье, скальные массивы и горную реку Чегем.',
  1,
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00'
);

INSERT INTO routes_route_type_links (route_id, route_type_id) VALUES (2, 2); -- Авто
INSERT INTO routes_places_links (route_id, place_id, place_order) VALUES (2, 2, 1); -- Чегемские водопады

-- 3. Экскурсия к Голубым озерам
INSERT INTO routes (id, name, slug, description, is_active, created_at, updated_at, published_at)
VALUES (
  3,
  'Экскурсия к Голубым озерам',
  'ekskursiya-k-golubym-ozeram',
  'Увлекательная экскурсия к уникальным карстовым озерам с кристально чистой водой необычного бирюзового цвета. Включает посещение Нижнего голубого озера и смотровых площадок.',
  1,
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00'
);

INSERT INTO routes_route_type_links (route_id, route_type_id) VALUES (3, 2); -- Авто
INSERT INTO routes_places_links (route_id, place_id, place_order) VALUES (3, 3, 1); -- Голубые озера

-- 4. Треккинг по долине Нарзанов
INSERT INTO routes (id, name, slug, description, is_active, created_at, updated_at, published_at)
VALUES (
  4,
  'Треккинг по долине Нарзанов',
  'trekking-po-doline-narzanov',
  'Пеший маршрут через живописную долину с минеральными источниками. Подходит для акклиматизации перед восхождением на Эльбрус. Маршрут проходит по альпийским лугам с видами на главный Кавказский хребет.',
  1,
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00'
);

INSERT INTO routes_route_type_links (route_id, route_type_id) VALUES (4, 1); -- Пеший
INSERT INTO routes_places_links (route_id, place_id, place_order) VALUES (4, 4, 1); -- Долина Нарзанов

-- 5. Большое кольцо Приэльбрусья
INSERT INTO routes (id, name, slug, description, is_active, created_at, updated_at, published_at)
VALUES (
  5,
  'Большое кольцо Приэльбрусья',
  'bolshoe-kolco-prielbrusya',
  'Многодневный комбинированный маршрут по национальному парку Приэльбрусье с посещением ключевых достопримечательностей: Долины Нарзанов, поляны Азау, водопадов и смотровых площадок с видом на Эльбрус.',
  1,
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00'
);

INSERT INTO routes_route_type_links (route_id, route_type_id) VALUES (5, 3); -- Комбинированный
INSERT INTO routes_places_links (route_id, place_id, place_order) VALUES (5, 1, 1); -- Гора Эльбрус
INSERT INTO routes_places_links (route_id, place_id, place_order) VALUES (5, 4, 2); -- Долина Нарзанов
INSERT INTO routes_places_links (route_id, place_id, place_order) VALUES (5, 5, 3); -- Национальный парк

-- 6. Культурно-исторический маршрут
INSERT INTO routes (id, name, slug, description, is_active, created_at, updated_at, published_at)
VALUES (
  6,
  'Культурно-исторический маршрут',
  'kulturno-istoricheskij-marshrut',
  'Автомобильный маршрут по историческим местам Эльбрусского района с посещением средневековых башен, древних селений и музеев местной культуры.',
  1,
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00',
  '2025-11-21 16:30:00'
);

INSERT INTO routes_route_type_links (route_id, route_type_id) VALUES (6, 2); -- Авто
INSERT INTO routes_places_links (route_id, place_id, place_order) VALUES (6, 6, 1); -- Башня Амирхана

-- Проверка
SELECT 'Добавлено маршрутов:' AS message;
SELECT COUNT(*) AS total FROM routes;
SELECT id, name, slug FROM routes ORDER BY id;

