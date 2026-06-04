-- name: CreateInventoryItemsBatch :exec
INSERT INTO inventories (
    id,
    sku_id,
    shop_id,
    quantity,
    reserved,
    status,
    created_by,
    updated_by,
    created_at,
    updated_at
)
SELECT
    unnest(@ids::uuid[]),
    unnest(@sku_ids::uuid[]),
    @shop_id::uuid,
    unnest(@quantities::bigint[]),
    unnest(@reserved_quantities::bigint[]),
    unnest(@statuses::text[]),
    @created_by::uuid,
    @updated_by::uuid,
    @created_at::timestamptz,
    NULL::timestamptz;

-- name: SoftDeleteInventoryItemsBySkuIDs :execrows
UPDATE inventories
SET
    status = @deleted_status::text,
    updated_by = @actor_id::uuid,
    updated_at = @updated_at::timestamptz
WHERE sku_id = ANY(@sku_ids::uuid[])
  AND status = @active_status::text;

-- name: ListInventoryItemsBySkuIDs :many
SELECT
    id,
    sku_id,
    shop_id,
    quantity,
    reserved,
    status,
    created_by,
    updated_by,
    created_at,
    updated_at
FROM inventories
WHERE sku_id = ANY(@sku_ids::uuid[])
ORDER BY created_at ASC;

-- name: GetActiveInventoryItemByShopAndSku :one
SELECT
    id,
    sku_id,
    shop_id,
    quantity,
    reserved,
    status,
    created_by,
    updated_by,
    created_at,
    updated_at
FROM inventories
WHERE shop_id = @shop_id::uuid
  AND sku_id = @sku_id::uuid
  AND status = @active_status::text
LIMIT 1;


-- name: ReserveInventoryStockBatch :execrows
UPDATE inventories AS i
SET 
    reserved = i.reserved + u.req_quantity,
    updated_by = @actor_id::uuid,
    created_at = @created_at::timestamptz
FROM (
    SELECT * FROM ROWS FROM (
        unnest(@inventory_ids::uuid[]), 
        unnest(@quantities::bigint[])
    ) AS tmp(inventory_id, req_quantity)
) AS u
WHERE i.shop_id = @shop_id::uuid
  AND i.id = u.inventory_id
  AND (i.quantity - i.reserved) >= u.req_quantity;

-- name: ReleaseInventoryStockBatch :execrows
UPDATE inventories AS i
SET
    reserved = i.reserved - u.req_quantity,
    updated_by = @actor_id::uuid,
    updated_at = @updated_at::timestamptz
FROM (
    SELECT * FROM ROWS FROM (
        unnest(@inventory_ids::uuid[]),
        unnest(@quantities::bigint[])
    ) AS tmp(inventory_id, req_quantity)
) AS u
WHERE i.shop_id = @shop_id::uuid
    AND i.id = u.inventory_id
    AND i.reserved >= u.req_quantity;


-- name: FulfillInventoryStockBatch :execrows
UPDATE inventories AS i
SET
    quantity = i.quantity - u.req_quantity,
    reserved = i.reserved - u.req_quantity,
    updated_by = @actor_id::uuid,
    updated_at = @updated_at::timestamptz
FROM (
    SELECT * FROM ROWS FROM (
        unnest(@inventory_ids::uuid[]),
        unnest(@quantities::bigint[])
    ) AS tmp(inventory_id, req_quantity)
) AS u
WHERE i.shop_id = @shop_id::uuid
    AND i.id = u.inventory_id
    AND i.quantity >= u.req_quantity
    AND i.reserved >= u.req_quantity;