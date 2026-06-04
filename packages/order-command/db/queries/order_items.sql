-- name: CreateInventoryTransactionBatch :exec
INSERT INTO inventory_transactions (
    id,
    inventory_id,
    type,
    quantity,
    balance_before,
    balance_after,
    reference_type,
    reference_id,
    action_type,
    idempotency_key,
    note,
    created_by,
    created_at
) SELECT
    unnest(@ids::uuid[]),
    unnest(@inventory_ids::uuid[]),
    unnest(@types::text[]),
    unnest(@quantities::bigint[]),
    unnest(@balances_before::bigint[]),
    unnest(@balances_after::bigint[]),
    unnest(@reference_types::text[]),
    unnest(@reference_ids::text[]),
    unnest(@action_types::text[]),
    unnest(@idempotency_keys::text[]),
    unnest(@notes::text[]),
    unnest(@created_bys::uuid[]),
    unnest(@created_ats::timestamptz[]);

