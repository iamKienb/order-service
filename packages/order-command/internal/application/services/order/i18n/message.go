package i18n

const (
	MsgInventoryItemInvalid       = "inventory item payload is invalid"
	MsgInventoryQuantityInvalid   = "inventory quantity must be greater than or equal to zero"
	MsgInventoryItemAlreadyExists = "inventory item for this sku already exists"
	MsgInventoryItemNotFound      = "inventory item was not found"
	MsgInventoryStockNotEnough    = "inventory stock was not enough"

	MsgInventoryReleaseStockFailed = ("inventory failed to release stock")
	MsgInventoryFulfilStockFailed  = ("inventory failed to fulfill stock")

	MsgInventoryReservationInvalid = "inventory reservation is invalid"
	MsgInventoryReservationMissing = "inventory reservation was not found"

	MsgInsufficientStock = "inventory stock is insufficient"
)
