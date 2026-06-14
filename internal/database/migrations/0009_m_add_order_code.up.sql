ALTER TABLE orders ADD COLUMN order_code VARCHAR(20) DEFAULT NULL;
CREATE INDEX idx_orders_order_code ON orders (order_code);
