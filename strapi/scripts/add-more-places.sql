-- Добавление дополнительных мест Кабардино-Балкарии
-- Останавливайте Strapi перед выполнением!

-- 2. Чегемские водопады
INSERT INTO places (id, name, slug, history, address, latitude, longitude, working_hours, phone, website, is_active, created_at, updated_at, published_at)
VALUES (
  2,
  'Чегемские водопады',
  'chegemskie-vodopady',
  'Чегемские водопады — одна из самых живописных природных достопримечательностей Кабардино-Балкарии. Водопады расположены в Чегемском ущелье и представляют собой каскад замерзающих зимой потоков воды, стекающих с высоких скал. Особенно впечатляющее зрелище открывается зимой, когда водопады превращаются в гигантские ледяные колонны.',
  'Чегемский район, Кабардино-Балкария',
  43.2733,
  42.9217,
  'Круглосуточно',
  '+7 (8662) 40-00-00',
  null,
  1,
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00'
);

INSERT INTO places_categories_links (place_id, category_id) VALUES (2, 12); -- Экотуризм
INSERT INTO places_categories_links (place_id, category_id) VALUES (2, 17); -- Активный туризм
INSERT INTO places_tags_links (place_id, tag_id) VALUES (2, 2); -- Водопад
INSERT INTO places_tags_links (place_id, tag_id) VALUES (2, 6); -- Каньон
INSERT INTO places_area_links (place_id, area_id) VALUES (2, 9); -- Чегемский район

-- 3. Голубые озера
INSERT INTO places (id, name, slug, history, address, latitude, longitude, working_hours, phone, website, is_active, created_at, updated_at, published_at)
VALUES (
  3,
  'Голубые озера',
  'golubye-ozera',
  'Голубые озера — группа из пяти карстовых озер в Черекском районе Кабардино-Балкарии. Главное озеро Церик-Кёль (Нижнее голубое озеро) является одним из глубочайших карстовых озер в мире. Вода в озере имеет необычный бирюзовый цвет из-за присутствия сероводорода и преломления солнечного света. Температура воды круглый год около +9°C.',
  'Черекский район, Кабардино-Балкария',
  43.3333,
  43.5333,
  'Круглосуточно',
  '+7 (8662) 40-00-00',
  null,
  1,
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00'
);

INSERT INTO places_categories_links (place_id, category_id) VALUES (3, 12); -- Экотуризм
INSERT INTO places_categories_links (place_id, category_id) VALUES (3, 19); -- Курортно-оздоровительный
INSERT INTO places_tags_links (place_id, tag_id) VALUES (3, 4); -- Озеро
INSERT INTO places_area_links (place_id, area_id) VALUES (3, 10); -- Черекский район

-- 4. Долина Нарзанов
INSERT INTO places (id, name, slug, history, address, latitude, longitude, working_hours, phone, website, is_active, created_at, updated_at, published_at)
VALUES (
  4,
  'Долина Нарзанов',
  'dolina-narzanov',
  'Долина Нарзанов — уникальное место в Приэльбрусье, где из-под земли бьют более 20 источников минеральной воды, насыщенной железом. Вода в источниках имеет характерный рыжий цвет из-за окисления железа. Место популярно среди туристов и альпинистов как остановка на пути к Эльбрусу. Воздух здесь насыщен ионами, что оказывает благотворное влияние на здоровье.',
  'Эльбрусский район, Кабардино-Балкария',
  43.2167,
  42.5167,
  'Круглосуточно',
  '+7 (8662) 40-00-00',
  null,
  1,
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00'
);

INSERT INTO places_categories_links (place_id, category_id) VALUES (4, 12); -- Экотуризм
INSERT INTO places_categories_links (place_id, category_id) VALUES (4, 19); -- Курортно-оздоровительный
INSERT INTO places_tags_links (place_id, tag_id) VALUES (4, 8); -- Долина
INSERT INTO places_tags_links (place_id, tag_id) VALUES (4, 3); -- Река
INSERT INTO places_area_links (place_id, area_id) VALUES (4, 10); -- Эльбрусский район

-- 5. Национальный парк Приэльбрусье
INSERT INTO places (id, name, slug, history, address, latitude, longitude, working_hours, phone, website, is_active, created_at, updated_at, published_at)
VALUES (
  5,
  'Национальный парк Приэльбрусье',
  'nacionalnyj-park-prielbrusie',
  'Национальный парк Приэльбрусье — особо охраняемая природная территория в центральной части Кавказских гор. Основан в 1986 году для сохранения уникальных природных комплексов Приэльбрусья. На территории парка находятся высочайшие вершины России и Европы, в том числе Эльбрус, древние ледники, альпийские луга, минеральные источники. Это популярное место для альпинизма, горнолыжного спорта и экологического туризма.',
  'Эльбрусский район, Кабардино-Балкария',
  43.3333,
  42.4500,
  '09:00-18:00',
  '+7 (86638) 7-14-38',
  'https://elbrus-ngp.ru',
  1,
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00'
);

INSERT INTO places_categories_links (place_id, category_id) VALUES (5, 12); -- Экотуризм
INSERT INTO places_categories_links (place_id, category_id) VALUES (5, 17); -- Активный туризм
INSERT INTO places_categories_links (place_id, category_id) VALUES (5, 18); -- Горнолыжный
INSERT INTO places_tags_links (place_id, tag_id) VALUES (5, 1); -- Горы
INSERT INTO places_tags_links (place_id, tag_id) VALUES (5, 7); -- Лес
INSERT INTO places_tags_links (place_id, tag_id) VALUES (5, 12); -- Ледник
INSERT INTO places_area_links (place_id, area_id) VALUES (5, 10); -- Эльбрусский район

-- 6. Башня Амирхана
INSERT INTO places (id, name, slug, history, address, latitude, longitude, working_hours, phone, website, is_active, created_at, updated_at, published_at)
VALUES (
  6,
  'Башня Амирхана',
  'bashnya-amirhana',
  'Башня Амирхана — средневековая сторожевая башня в селении Эль-Тюбю, образец балкарской оборонительной архитектуры XV-XVII веков. Башня построена из местного камня без использования раствора и достигает высоты около 15 метров. С башни открывается впечатляющий вид на Эльбрус и окрестные горы. Является объектом культурного наследия и популярным туристическим местом.',
  'с. Эль-Тюбю, Эльбрусский район, Кабардино-Балкария',
  43.2500,
  42.4833,
  'Круглосуточно',
  null,
  null,
  1,
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00',
  '2025-11-21 16:00:00'
);

INSERT INTO places_categories_links (place_id, category_id) VALUES (6, 15); -- Историко-культурный
INSERT INTO places_categories_links (place_id, category_id) VALUES (6, 22); -- Арт-туризм
INSERT INTO places_tags_links (place_id, tag_id) VALUES (6, 9); -- Видовая точка
INSERT INTO places_area_links (place_id, area_id) VALUES (6, 10); -- Эльбрусский район

-- Проверка
SELECT 'Добавлено мест:' AS message;
SELECT COUNT(*) AS total FROM places;
SELECT id, name, slug FROM places ORDER BY id;

