package order

import "errors"

var (
	ErrOrderInvalidGrandTotal          = errors.New("ORDER_INVALID_GRAND_TOTAL")
	ErrOrderCannotCancelInCurrentState = errors.New("ORDER_CANNOT_CANCEL_IN_CURRENT_STATE")
	ErrOrderInvalidStateTransition     = errors.New("ORDER_INVALID_STATE_TRANSITION")
	ErrOrderItemsCannotBeEmpty         = errors.New("ORDER_ITEMS_CANNOT_BE_EMPTY")
	ErrOrderItemInvalidCalculation     = errors.New("ORDER_ITEM_INVALID_CALCULATION")
	ErrOrderCurrencyMismatch           = errors.New("ORDER_CURRENCY_MISMATCH")

	ErrOrderCannotCancel     = errors.New("order_cannot_cancel_in_current_status")
	ErrOrderBuyerNotAllowed  = errors.New("buyer_cannot_cancel_once_confirmed")
	ErrOrderShopNotAllowed   = errors.New("shop_cannot_cancel_once_shipped")
	ErrOrderActorIDRequired  = errors.New("actor_id_is_required")
	ErrOrderInvalidActorType = errors.New("invalid_actor_type")
)
