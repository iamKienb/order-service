ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS idempotency_key TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS idx_orders_buyer_idempotency_key
    ON orders(buyer_id, idempotency_key)
    WHERE idempotency_key IS NOT NULL;
