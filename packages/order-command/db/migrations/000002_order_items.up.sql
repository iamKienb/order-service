CREATE TABLE order_items (
    id              UUID PRIMARY KEY,
    order_id        TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    
    inventory_id    UUID         NOT NULL,
    sku_id          UUID         NOT NULL,
    sku_code        VARCHAR(100) NOT NULL,
    product_id      UUID         NOT NULL,
    product_name    VARCHAR(500) NOT NULL,
    image_url       TEXT NOT NULL,
    quantity        BIGINT          NOT NULL CHECK (quantity > 0),

    base_price      BIGINT NOT NULL,      
    item_subtotal   BIGINT NOT NULL,        
    currency        VARCHAR(10)   NOT NULL DEFAULT 'VND',

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_items_order_id ON order_items(order_id);
