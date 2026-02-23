-- file: migrations/000001_initial_schema.down.sql

DROP TRIGGER IF EXISTS update_reservations_updated_at ON reservations;
DROP TRIGGER IF EXISTS update_stock_items_updated_at ON stock_items;
DROP TRIGGER IF EXISTS update_warehouses_updated_at ON warehouses;
DROP TRIGGER IF EXISTS update_product_variants_updated_at ON product_variants;
DROP TRIGGER IF EXISTS update_products_updated_at ON products;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS low_stock_alerts;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS stock_movements;
DROP TABLE IF EXISTS stock_items;
DROP TABLE IF EXISTS warehouses;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;

DROP TYPE IF EXISTS reservation_status;
DROP TYPE IF EXISTS stock_movement_type;

DROP EXTENSION IF EXISTS "uuid-ossp";