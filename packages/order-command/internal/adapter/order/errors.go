package order

import (
	"errors"

	app_order "order-command-module/internal/application/services/order"
	domain_order "order-command-module/internal/domain/order"

	"github.com/iamKienb/go-core/app_error"
)

const (
	errCodeCheckoutItemInvalid      = "checkout_item_invalid"
	errCodeCheckoutItemUnavailable  = "checkout_item_unavailable"
	errCodeOrderShopMismatch        = "order_shop_mismatch"
	errCodeOrderInvalidTransition   = "order_invalid_transition"
	errCodeOrderCannotCancel        = "order_cannot_cancel"
	errCodeOrderActorInvalid        = "order_actor_invalid"
	errCodeOrderActorActionRejected = "order_actor_action_rejected"

	errMsgCheckoutItemInvalid      = "checkout item is invalid"
	errMsgCheckoutItemUnavailable  = "checkout item is unavailable"
	errMsgOrderShopMismatch        = "order does not belong to this shop"
	errMsgOrderInvalidTransition   = "order status transition is invalid"
	errMsgOrderCannotCancel        = "order cannot be cancelled"
	errMsgOrderActorInvalid        = "order actor is invalid"
	errMsgOrderActorActionRejected = "order actor is not allowed to perform this action"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, app_order.ErrCheckoutItemInvalid):
		return app_error.New(app_error.KindValidation, errCodeCheckoutItemInvalid, errMsgCheckoutItemInvalid, err)
	case errors.Is(err, app_order.ErrCheckoutItemUnavailable):
		return app_error.New(app_error.KindConflict, errCodeCheckoutItemUnavailable, errMsgCheckoutItemUnavailable, err)
	case errors.Is(err, app_order.ErrOrderShopMismatch):
		return app_error.New(app_error.KindForbidden, errCodeOrderShopMismatch, errMsgOrderShopMismatch, err)
	case errors.Is(err, domain_order.ErrOrderInvalidStateTransition):
		return app_error.New(app_error.KindConflict, errCodeOrderInvalidTransition, errMsgOrderInvalidTransition, err)
	case errors.Is(err, domain_order.ErrOrderCannotCancel):
		return app_error.New(app_error.KindConflict, errCodeOrderCannotCancel, errMsgOrderCannotCancel, err)
	case errors.Is(err, domain_order.ErrOrderInvalidActorType), errors.Is(err, domain_order.ErrOrderActorIDRequired):
		return app_error.New(app_error.KindValidation, errCodeOrderActorInvalid, errMsgOrderActorInvalid, err)
	case errors.Is(err, domain_order.ErrOrderBuyerNotAllowed), errors.Is(err, domain_order.ErrOrderShopNotAllowed):
		return app_error.New(app_error.KindForbidden, errCodeOrderActorActionRejected, errMsgOrderActorActionRejected, err)
	default:
		return err
	}
}
