CREATE TABLE order_history (
    id          UUID PRIMARY KEY,
    order_id    TEXT         NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    from_status VARCHAR(50),        
    to_status   VARCHAR(50) NOT NULL, 
    changed_by  NOT NULL UUID,                
    actor_type  VARCHAR(20) NOT NULL DEFAULT 'SYSTEM',
    reason      NOT NULL TEXT,             
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_order_history_order ON order_history(order_id);
