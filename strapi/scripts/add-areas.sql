-- Скрипт для добавления районов Кабардино-Балкарии в Strapi
-- Использование: sqlite3 .tmp/data.db < scripts/add-areas.sql

INSERT OR IGNORE INTO areas (name, slug, "order", is_active, created_at, updated_at) VALUES
('Баксанский район', 'baksanskij-rajon', 1, 1, datetime('now'), datetime('now')),
('Зольский район', 'zol-skij-rajon', 2, 1, datetime('now'), datetime('now')),
('Лескенский район', 'leskenskij-rajon', 3, 1, datetime('now'), datetime('now')),
('Майский район', 'majskij-rajon', 4, 1, datetime('now'), datetime('now')),
('Прохладненский район', 'prohladnenskij-rajon', 5, 1, datetime('now'), datetime('now')),
('Терский район', 'terskij-rajon', 6, 1, datetime('now'), datetime('now')),
('Урванский район', 'urvanskij-rajon', 7, 1, datetime('now'), datetime('now')),
('Чегемский район', 'chegemskij-rajon', 8, 1, datetime('now'), datetime('now')),
('Черекский район', 'cherekskij-rajon', 9, 1, datetime('now'), datetime('now')),
('Эльбрусский район', 'el-brusskij-rajon', 10, 1, datetime('now'), datetime('now'));

SELECT 'Добавлено районов: ' || COUNT(*) FROM areas;

