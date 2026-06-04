package order

import (
	"inventory-command-module/internal/application/services/inventory/i18n"
	domain_inventory "inventory-command-module/internal/domain/inventory"

	"github.com/iamKienb/go-core/app_error"
)

var inventoryErrorMap = app_error.ServiceErrorMap{
	domain_inventory.ErrInventoryItemInvalid:       {Kind: app_error.KindValidation, Msg: i18n.MsgInventoryItemInvalid},
	domain_inventory.ErrInventoryQuantityInvalid:   {Kind: app_error.KindValidation, Msg: i18n.MsgInventoryQuantityInvalid},
	domain_inventory.ErrInventoryItemAlreadyExists: {Kind: app_error.KindConflict, Msg: i18n.MsgInventoryItemAlreadyExists},
	domain_inventory.ErrInventoryItemNotFound:      {Kind: app_error.KindNotFound, Msg: i18n.MsgInventoryItemNotFound},
	domain_inventory.ErrInventoryStockNotEnough:    {Kind: app_error.KindConflict, Msg: i18n.MsgInventoryStockNotEnough},

	domain_inventory.ErrInventoryReservationInvalid:  {Kind: app_error.KindValidation, Msg: i18n.MsgInventoryReservationInvalid},
	domain_inventory.ErrInventoryReservationNotFound: {Kind: app_error.KindNotFound, Msg: i18n.MsgInventoryReservationMissing},

	domain_inventory.ErrInventoryReleaseStockFailed: {Kind: app_error.KindValidation, Msg: i18n.MsgInventoryReleaseStockFailed},
	domain_inventory.ErrInventoryFulfilStockFailed:  {Kind: app_error.KindValidation, Msg: i18n.MsgInventoryFulfilStockFailed},

	domain_inventory.ErrInsufficientStock: {Kind: app_error.KindConflict, Msg: i18n.MsgInsufficientStock},
}

func (s *orderService) wrapError(err error) error {
	return app_error.WrapError(err, inventoryErrorMap)
}
