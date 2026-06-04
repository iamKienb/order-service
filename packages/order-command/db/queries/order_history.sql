-- name: CreateInventoryReservationBatch :exec
INSERT INTO inventory_reservations (
    id,
    inventory_id,
    order_id,
    quantity,
    status,
    expires_at,
    created_at,
    updated_at
) SELECT
    unnest(@ids::uuid[]),
    unnest(@inventory_ids::uuid[]),
    unnest(@order_ids::text[]),
    @status::text,
    @expires_at::timestamptz,
    @created_at::timestamptz,
    @updated_at::timestamptz;


-- name: ReleaseReservationsByOrderID :many
UPDATE inventory_reservations
SET
    status = 'RELEASED',
    updated_at = @updated_at::timestamptz
WHERE order_id = @order_id::text 
    AND status = 'HOLD'
RETURNING *;


-- name: FulfillReservationsByOrderID :many
UPDATE inventory_reservations
SET
    status = 'COMPLETED',
    updated_at = @updated_at::timestamptz
WHERE order_id = @order_id::text 
    AND status = 'HOLD'
RETURNING *;