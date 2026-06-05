package order

import (
	"errors"

	app_order "order-command-module/internal/application/services/order"
	domain_order "order-command-module/internal/domain/order"

	"github.com/iamKienb/go-core/app_error"
)

const (
	errCodeCheckoutItemInvalid      = "checkout_item_invalid"
	errCodeCheckoutItemEmpty        = "checkout_item_empty"
	errCodeCheckoutItemUnavailable  = "checkout_item_unavailable"
	errCodeOrderShopMismatch        = "order_shop_mismatch"
	errCodeOrderInvalidTransition   = "order_invalid_transition"
	errCodeOrderCannotCancel        = "order_cannot_cancel"
	errCodeOrderActorInvalid        = "order_actor_invalid"
	errCodeOrderActorActionRejected = "order_actor_action_rejected"
	errCodeOrderNotFound            = "order_not_found"
	errCodeOrderIdempotencyMissing  = "order_idempotency_key_missing"
	errCodeOrderIdempotencyConflict = "order_idempotency_key_conflict"

	errMsgCheckoutItemInvalid      = "checkout item is invalid"
	errMsgCheckoutItemEmpty        = "checkout items are required"
	errMsgCheckoutItemUnavailable  = "checkout item is unavailable"
	errMsgOrderShopMismatch        = "order does not belong to this shop"
	errMsgOrderInvalidTransition   = "order status transition is invalid"
	errMsgOrderCannotCancel        = "order cannot be cancelled"
	errMsgOrderActorInvalid        = "order actor is invalid"
	errMsgOrderActorActionRejected = "order actor is not allowed to perform this action"
	errMsgOrderNotFound            = "order not found"
	errMsgOrderIdempotencyMissing  = "order idempotency key is required"
	errMsgOrderIdempotencyConflict = "order idempotency key is already used by a different request"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, app_order.ErrCheckoutItemInvalid):
		return app_error.New(app_error.KindValidation, errCodeCheckoutItemInvalid, errMsgCheckoutItemInvalid, err)
	case errors.Is(err, app_order.ErrCheckoutItemEmpty):
		return app_error.New(app_error.KindValidation, errCodeCheckoutItemEmpty, errMsgCheckoutItemEmpty, err)
	case errors.Is(err, app_order.ErrCheckoutItemUnavailable):
		return app_error.New(app_error.KindConflict, errCodeCheckoutItemUnavailable, errMsgCheckoutItemUnavailable, err)
	case errors.Is(err, app_order.ErrOrderShopMismatch):
		return app_error.New(app_error.KindForbidden, errCodeOrderShopMismatch, errMsgOrderShopMismatch, err)
	case errors.Is(err, app_order.ErrOrderIdempotencyMissing):
		return app_error.New(app_error.KindValidation, errCodeOrderIdempotencyMissing, errMsgOrderIdempotencyMissing, err)
	case errors.Is(err, domain_order.ErrOrderNotFound):
		return app_error.New(app_error.KindNotFound, errCodeOrderNotFound, errMsgOrderNotFound, err)
	case errors.Is(err, domain_order.ErrOrderIdempotencyKeyConflict):
		return app_error.New(app_error.KindConflict, errCodeOrderIdempotencyConflict, errMsgOrderIdempotencyConflict, err)
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
