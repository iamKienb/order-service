package order

import "errors"

var (
	ErrCheckoutItemUnavailable = errors.New("ORDER_CHECKOUT_ITEM_UNAVAILABLE")
	ErrCheckoutItemInvalid     = errors.New("ORDER_CHECKOUT_ITEM_INVALID")
	ErrOrderShopMismatch       = errors.New("ORDER_SHOP_MISMATCH")
	ErrOrderIdempotencyMissing = errors.New("ORDER_IDEMPOTENCY_KEY_MISSING")
)
