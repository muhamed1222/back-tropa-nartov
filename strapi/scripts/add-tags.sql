-- Скрипт для добавления тегов мест в Strapi
-- Использование: sqlite3 .tmp/data.db < scripts/add-tags.sql

INSERT OR IGNORE INTO tags (name, slug, "order", is_active, created_at, updated_at) VALUES
('Горы', 'gory', 1, 1, datetime('now'), datetime('now')),
('Водопад', 'vodopad', 2, 1, datetime('now'), datetime('now')),
('Река', 'reka', 3, 1, datetime('now'), datetime('now')),
('Озеро', 'ozero', 4, 1, datetime('now'), datetime('now')),
('Каньон', 'kan-on', 5, 1, datetime('now'), datetime('now')),
('Лес', 'les', 6, 1, datetime('now'), datetime('now')),
('Плато', 'plato', 7, 1, datetime('now'), datetime('now')),
('Долина', 'dolina', 8, 1, datetime('now'), datetime('now')),
('Видовая точка', 'vidovaya-tochka', 9, 1, datetime('now'), datetime('now')),
('Пещера', 'peschera', 10, 1, datetime('now'), datetime('now')),
('Ледник', 'lednik', 11, 1, datetime('now'), datetime('now'));

SELECT 'Добавлено тегов: ' || COUNT(*) FROM tags;

