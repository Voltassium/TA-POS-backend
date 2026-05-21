ALTER TABLE order_items DROP CONSTRAINT IF EXISTS chk_served_qty;
ALTER TABLE order_items DROP COLUMN IF EXISTS served_qty;
