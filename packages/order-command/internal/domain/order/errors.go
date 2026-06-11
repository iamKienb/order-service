package order

import "errors"

var (
	ErrOrderNotFound                   = errors.New("ORDER_NOT_FOUND")
	ErrOrderInvalidGrandTotal          = errors.New("ORDER_INVALID_GRAND_TOTAL")
	ErrOrderInvalidStateTransition     = errors.New("ORDER_INVALID_STATE_TRANSITION")
	ErrOrderCannotCancelInCurrentState = errors.New("ORDER_CANNOT_CANCEL_IN_CURRENT_STATE")
	ErrOrderCurrencyMismatch           = errors.New("ORDER_CURRENCY_MISMATCH")
	ErrOrderIdempotencyKeyConflict     = errors.New("ORDER_IDEMPOTENCY_KEY_CONFLICT")

	ErrOrderActorIDRequired  = errors.New("ORDER_ACTOR_ID_REQUIRED")
	ErrOrderInvalidActorType = errors.New("ORDER_INVALID_ACTOR_TYPE")
	ErrOrderBuyerNotAllowed  = errors.New("ORDER_BUYER_NOT_ALLOWED")
	ErrOrderShopNotAllowed   = errors.New("ORDER_SHOP_NOT_ALLOWED")

	ErrCheckoutItemEmpty           = errors.New("ORDER_CHECKOUT_ITEM_EMPTY")
	ErrCheckoutItemInvalid         = errors.New("ORDER_CHECKOUT_ITEM_INVALID")
	ErrCheckoutItemUnavailable     = errors.New("ORDER_CHECKOUT_ITEM_UNAVAILABLE")
	ErrOrderItemInvalidCalculation = errors.New("ORDER_ITEM_INVALID_CALCULATION")
	ErrOrderShopMismatch           = errors.New("ORDER_SHOP_MISMATCH")
	ErrOrderIdempotencyMissing     = errors.New("ORDER_IDEMPOTENCY_KEY_MISSING")
)
