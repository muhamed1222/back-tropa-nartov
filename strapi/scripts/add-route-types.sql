-- Скрипт для добавления типов маршрутов в Strapi
-- Использование: sqlite3 .tmp/data.db < scripts/add-route-types.sql

INSERT OR IGNORE INTO route_types (name, slug, "order", is_active, created_at, updated_at) VALUES
('Пеший', 'peshij', 1, 1, datetime('now'), datetime('now')),
('Авто', 'avto', 2, 1, datetime('now'), datetime('now'));

SELECT 'Добавлено типов маршрутов: ' || COUNT(*) FROM route_types;

