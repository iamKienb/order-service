-- name: CreateOrderHistoryBatch :exec
INSERT INTO order_history (
    id,
    order_id,
    from_status,
    to_status,
    changed_by,
    actor_type,
    reason,
    created_at
)
SELECT
    unnest(@ids::uuid[]),
    @order_id::text,
    unnest(@from_statuses::text[]),
    unnest(@to_statuses::text[]),
    unnest(@changed_bys::uuid[]),
    unnest(@actor_types::text[]),
    unnest(@reasons::text[]),
    @created_at::timestamptz;

-- name: ListOrderHistoryByOrderID :many
SELECT *
FROM order_history
WHERE order_id = @order_id::text
ORDER BY created_at ASC;
