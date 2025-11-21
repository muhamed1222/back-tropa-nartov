-- Скрипт для добавления категорий туризма в Strapi (упрощенная версия)
-- Использование: sqlite3 .tmp/data.db < scripts/add-categories-simple.sql

INSERT OR IGNORE INTO categories (name, slug, "order", is_active, created_at, updated_at) VALUES
('Экотуризм', 'ekotourizm', 1, 1, datetime('now'), datetime('now')),
('Гастротуризм', 'gastrotourizm', 2, 1, datetime('now'), datetime('now')),
('Агротуризм', 'agrotourizm', 3, 1, datetime('now'), datetime('now')),
('Историко-культурный туризм', 'istoriko-kul-turnyj-turizm', 4, 1, datetime('now'), datetime('now')),
('Экстремальный туризм', 'ekstremal-nyj-turizm', 5, 1, datetime('now'), datetime('now')),
('Активный туризм', 'aktivnyj-turizm', 6, 1, datetime('now'), datetime('now')),
('Горнолыжный туризм', 'gornolyzhnyj-turizm', 7, 1, datetime('now'), datetime('now')),
('Курортно-оздоровительный туризм', 'kurortno-ozdorovitel-nyj-turizm', 8, 1, datetime('now'), datetime('now')),
('Урбан-туризм', 'urban-turizm', 9, 1, datetime('now'), datetime('now')),
('Сельский туризм', 'sel-skij-turizm', 10, 1, datetime('now'), datetime('now')),
('Арт-туризм', 'art-turizm', 11, 1, datetime('now'), datetime('now'));

SELECT 'Добавлено категорий: ' || COUNT(*) FROM categories;

