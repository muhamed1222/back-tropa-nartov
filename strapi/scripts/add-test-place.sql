-- Добавление тестового места в Strapi для проверки интеграции с Flutter
-- Использование: sqlite3 .tmp/data.db < scripts/add-test-place.sql

-- Добавляем тестовое место
INSERT INTO places (
  name, 
  slug, 
  history, 
  address, 
  latitude, 
  longitude, 
  working_hours, 
  phone, 
  website, 
  is_active,
  published_at,
  created_at, 
  updated_at
) VALUES (
  'Гора Эльбрус',
  'gora-el-brus',
  '<p>Эльбрус — высочайшая вершина России и Европы, расположенная на Кавказе. Это двувершинный конус потухшего вулкана. Высота западной вершины — 5642 м, восточной — 5621 м.</p>',
  'Эльбрусский район, Кабардино-Балкарская Республика',
  43.3499,
  42.4453,
  'Круглосуточно',
  '+7 (928) 123-45-67',
  'https://elbrus-tourism.ru',
  1,
  datetime('now'),
  datetime('now'),
  datetime('now')
);

-- Получаем ID добавленного места
SELECT 'Добавлено место ID: ' || last_insert_rowid() as result;

-- Связываем с районом (Эльбрусский район - ID 10)
INSERT INTO places_area_links (place_id, area_id)
SELECT last_insert_rowid(), 10;

-- Связываем с категориями (Экотуризм - ID 12, Активный туризм - ID 17)
INSERT INTO places_categories_links (place_id, category_id)
SELECT last_insert_rowid(), 12
UNION ALL
SELECT last_insert_rowid(), 17;

-- Связываем с тегами (Горы - ID 1)
INSERT INTO places_tags_links (place_id, tag_id)
SELECT last_insert_rowid(), 1;

SELECT 'Готово! Место "Гора Эльбрус" добавлено' as status;
SELECT 'Проверьте: http://localhost:1337/api/places' as check_url;

