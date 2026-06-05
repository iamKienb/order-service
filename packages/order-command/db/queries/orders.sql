-- name: CreateOrder :exec
INSERT INTO orders (
    id,
    shop_id,
    buyer_id,
    idempotency_key,
    status,
    shipping_name,
    shipping_phone,
    shipping_address,
    shipping_province,
    shipping_ward,
    note,
    grand_total,
    currency,
    cancel_reason,
    cancelled_by,
    confirmed_at,
    delivered_at,
    shipped_at,
    completed_at,
    cancelled_at,
    failed_at,
    created_at
) VALUES (
    @id::text,
    @shop_id::uuid,
    @buyer_id::uuid,
    @idempotency_key::text,
    @status::text,
    @shipping_name::text,
    @shipping_phone::text,
    @shipping_address::text,
    @shipping_province::text,
    @shipping_ward::text,
    @note::text,
    @grand_total::bigint,
    @currency::text,
    @cancel_reason::text,
    @cancelled_by::uuid,
    @confirmed_at::timestamptz,
    @delivered_at::timestamptz,
    @shipped_at::timestamptz,
    @completed_at::timestamptz,
    @cancelled_at::timestamptz,
    @failed_at::timestamptz,
    @created_at::timestamptz
);

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = @id::text
LIMIT 1;

-- name: GetOrderByBuyerAndIdempotencyKey :one
SELECT *
FROM orders
WHERE buyer_id = @buyer_id::uuid
  AND idempotency_key = @idempotency_key::text
LIMIT 1;

-- name: UpdateOrderStatus :execrows
UPDATE orders
SET
    status = @status::text,
    cancel_reason = @cancel_reason::text,
    cancelled_by = @cancelled_by::uuid,
    confirmed_at = @confirmed_at::timestamptz,
    delivered_at = @delivered_at::timestamptz,
    shipped_at = @shipped_at::timestamptz,
    completed_at = @completed_at::timestamptz,
    cancelled_at = @cancelled_at::timestamptz,
    failed_at = @failed_at::timestamptz
WHERE id = @id::text
  AND status = @expected_status::text;
