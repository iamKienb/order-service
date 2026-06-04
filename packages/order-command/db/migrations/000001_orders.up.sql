CREATE TABLE orders (
    id                  Text PRIMARY KEY,
    shop_id             UUID NOT NULL,
    buyer_id            UUID NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'PENDING', 

    shipping_name       VARCHAR(255) NOT NULL,
    shipping_phone      VARCHAR(20)  NOT NULL,
    shipping_address    TEXT         NOT NULL,
    shipping_province   VARCHAR(100) NOT NULL,
    shipping_ward       VARCHAR(100) NOT NULL,
    note                TEXT, -- Ghi chú của khách khi đặt hàng

    grand_total         BIGINT NOT NULL DEFAULT 0,       -- Số tiền cuối cùng khách phải trả
    currency            VARCHAR(10)  NOT NULL DEFAULT 'VND',

    cancel_reason       TEXT,                            
    cancelled_by        UUID,                            
    confirmed_at        TIMESTAMPTZ,
    delivered_at        TIMESTAMPTZ,
    shipped_at          TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,
    cancelled_at        TIMESTAMPTZ,
    failed_at           TIMESTAMPTZ,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_shop_buyer     ON orders(shop_id, buyer_id);
CREATE INDEX idx_orders_status_created ON orders(status, created_at DESC);
