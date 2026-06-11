package order

import (
	"context"

	"order-command-module/internal/application/commands/cancel_order"
	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/application/commands/preview_checkout"

	"connectrpc.com/connect"
	orderpb "github.com/iamKienb/api-contract/gen/order"
	"github.com/iamKienb/api-contract/gen/order/orderconnect"
	"github.com/iamKienb/go-core/app_error"
	authx "github.com/iamKienb/go-core/middleware/auth"
)

const errMsgAuthenticationRequired = "authentication required"

type orderServer struct {
	previewCheckoutExecutor preview_checkout.Executor
	placeOrderExecutor      place_order.Executor
	cancelOrderExecutor     cancel_order.Executor
	confirmOrderExecutor    confirm_order.Executor
}

func NewOrderServer(
	previewCheckoutExecutor preview_checkout.Executor,
	placeOrderExecutor place_order.Executor,
	cancelOrderExecutor cancel_order.Executor,
	confirmOrderExecutor confirm_order.Executor,
) *orderServer {
	return &orderServer{
		previewCheckoutExecutor: previewCheckoutExecutor,
		placeOrderExecutor:      placeOrderExecutor,
		cancelOrderExecutor:     cancelOrderExecutor,
		confirmOrderExecutor:    confirmOrderExecutor,
	}
}

func (s *orderServer) PreviewCheckout(ctx context.Context, req *connect.Request[orderpb.PreviewCheckoutRequest]) (*connect.Response[orderpb.PreviewCheckoutResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToPreviewCheckoutCommand(currentUser, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.previewCheckoutExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, toApplicationError(err)
	}

	return connect.NewResponse(ToPreviewCheckoutResponse(cmd.ShopID.String(), result)), nil
}

func (s *orderServer) PlaceOrder(ctx context.Context, req *connect.Request[orderpb.PlaceOrderRequest]) (*connect.Response[orderpb.PlaceOrderResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToPlaceOrderCommand(currentUser, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.placeOrderExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, toApplicationError(err)
	}

	return connect.NewResponse(ToPlaceOrderResponse(result)), nil
}

func (s *orderServer) CancelOrder(ctx context.Context, req *connect.Request[orderpb.CancelOrderRequest]) (*connect.Response[orderpb.CancelOrderResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToCancelOrderCommand(currentUser, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.cancelOrderExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, toApplicationError(err)
	}

	return connect.NewResponse(ToCancelOrderResponse(result)), nil
}

func (s *orderServer) ConfirmOrder(ctx context.Context, req *connect.Request[orderpb.ConfirmOrderRequest]) (*connect.Response[orderpb.ConfirmOrderResponse], error) {
	currentUser, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	cmd, err := ToConfirmOrderCommand(currentUser, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.confirmOrderExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, toApplicationError(err)
	}

	return connect.NewResponse(ToConfirmOrderResponse(result)), nil
}

func requireCurrentUser(ctx context.Context) (string, error) {
	claims := authx.GetUserInfoFromCtx(ctx)
	if claims == nil || claims.UserID == "" {
		return "", app_error.Unauthorized("authentication required")
	}

	return claims.UserID, nil
}

var _ orderconnect.OrderCommandHandler = (*orderServer)(nil)
