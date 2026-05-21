ALTER TABLE order_items
    ADD COLUMN served_qty INTEGER NOT NULL DEFAULT 0;

ALTER TABLE order_items
    ADD CONSTRAINT chk_served_qty CHECK (served_qty >= 0 AND served_qty <= quantity);
