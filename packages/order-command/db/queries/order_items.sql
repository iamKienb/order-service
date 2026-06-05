-- name: CreateOrderItemsBatch :exec
INSERT INTO order_items (
    id,
    order_id,
    inventory_id,
    sku_id,
    sku_code,
    product_id,
    product_name,
    image_url,
    quantity,
    base_price,
    item_subtotal,
    currency,
    created_at
)
SELECT
    unnest(@ids::uuid[]),
    @order_id::text,
    unnest(@inventory_ids::uuid[]),
    unnest(@sku_ids::uuid[]),
    unnest(@sku_codes::text[]),
    unnest(@product_ids::uuid[]),
    unnest(@product_names::text[]),
    unnest(@image_urls::text[]),
    unnest(@quantities::bigint[]),
    unnest(@base_prices::bigint[]),
    unnest(@item_subtotals::bigint[]),
    @currency::text,
    @created_at::timestamptz;

-- name: ListOrderItemsByOrderID :many
SELECT *
FROM order_items
WHERE order_id = @order_id::text
ORDER BY created_at ASC;
