package i18n

const (
	MsgInventoryItemInvalid        = "inventory item payload is invalid"
	MsgInventoryQuantityInvalid    = "inventory quantity must be greater than or equal to zero"
	MsgInventoryItemAlreadyExists  = "inventory item for this sku already exists"
	MsgInventoryItemNotFound       = "inventory item was not found"
	MsgInventoryStockNotEnough     = "inventory stock was not enough"
	MsgInventoryReleaseStockFailed = "inventory failed to release stock"
	MsgInventoryFulfilStockFailed  = "inventory failed to fulfill stock"
	MsgInventoryReservationInvalid = "inventory reservation is invalid"
	MsgInventoryReservationMissing = "inventory reservation was not found"
	MsgInsufficientStock           = "inventory stock is insufficient"

	MsgCheckoutItemInvalid      = "checkout item is invalid"
	MsgCheckoutItemEmpty        = "checkout items are required"
	MsgCheckoutItemUnavailable  = "checkout item is unavailable"
	MsgOrderShopMismatch        = "order does not belong to this shop"
	MsgOrderInvalidTransition   = "order status transition is invalid"
	MsgOrderCannotCancel        = "order cannot be cancelled"
	MsgOrderActorInvalid        = "order actor is invalid"
	MsgOrderActorActionRejected = "order actor is not allowed to perform this action"
	MsgOrderNotFound            = "order not found"
	MsgOrderIdempotencyMissing  = "order idempotency key is required"
	MsgOrderIdempotencyConflict = "order idempotency key is already used by a different request"
)
