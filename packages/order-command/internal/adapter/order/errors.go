package order

import (
	"order-command-module/internal/application/services/order/i18n"
	domain_order "order-command-module/internal/domain/order"

	"github.com/iamKienb/go-core/app_error"
)

var orderErrorMap = app_error.ServiceErrorMap{
	domain_order.ErrCheckoutItemEmpty:               {Kind: app_error.KindValidation, Msg: i18n.MsgCheckoutItemEmpty},
	domain_order.ErrCheckoutItemInvalid:             {Kind: app_error.KindValidation, Msg: i18n.MsgCheckoutItemInvalid},
	domain_order.ErrCheckoutItemUnavailable:         {Kind: app_error.KindConflict, Msg: i18n.MsgCheckoutItemUnavailable},
	domain_order.ErrOrderItemInvalidCalculation:     {Kind: app_error.KindValidation, Msg: i18n.MsgCheckoutItemInvalid},
	domain_order.ErrOrderCurrencyMismatch:           {Kind: app_error.KindValidation, Msg: i18n.MsgCheckoutItemInvalid},
	domain_order.ErrOrderShopMismatch:               {Kind: app_error.KindValidation, Msg: i18n.MsgOrderShopMismatch},
	domain_order.ErrOrderInvalidGrandTotal:          {Kind: app_error.KindValidation, Msg: i18n.MsgCheckoutItemInvalid},
	domain_order.ErrOrderIdempotencyMissing:         {Kind: app_error.KindValidation, Msg: i18n.MsgOrderIdempotencyMissing},
	domain_order.ErrOrderIdempotencyKeyConflict:     {Kind: app_error.KindConflict, Msg: i18n.MsgOrderIdempotencyConflict},
	domain_order.ErrOrderNotFound:                   {Kind: app_error.KindNotFound, Msg: i18n.MsgOrderNotFound},
	domain_order.ErrOrderInvalidStateTransition:     {Kind: app_error.KindValidation, Msg: i18n.MsgOrderInvalidTransition},
	domain_order.ErrOrderCannotCancelInCurrentState: {Kind: app_error.KindConflict, Msg: i18n.MsgOrderCannotCancel},
	domain_order.ErrOrderActorIDRequired:            {Kind: app_error.KindValidation, Msg: i18n.MsgOrderActorInvalid},
	domain_order.ErrOrderInvalidActorType:           {Kind: app_error.KindValidation, Msg: i18n.MsgOrderActorInvalid},
	domain_order.ErrOrderBuyerNotAllowed:            {Kind: app_error.KindForbidden, Msg: i18n.MsgOrderActorActionRejected},
	domain_order.ErrOrderShopNotAllowed:             {Kind: app_error.KindForbidden, Msg: i18n.MsgOrderActorActionRejected},
}

func toApplicationError(err error) error {
	return app_error.WrapError(err, orderErrorMap)
}
