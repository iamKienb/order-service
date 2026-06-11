package module

import (
	"log/slog"
	"net/http"

	orderadapter "order-command-module/internal/adapter/order"

	"connectrpc.com/grpcreflect"
	"github.com/iamKienb/api-contract/gen/order/orderconnect"
	observabilityx "github.com/iamKienb/go-core/middleware/observability"
)

type AdapterModule struct {
	Mux *http.ServeMux
}

func NewAdapterModule(app *ApplicationModule, logger *slog.Logger) *AdapterModule {
	allInterceptors := observabilityx.InternalServerOption(logger)

	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(orderconnect.OrderCommandName)

	orderServer := orderadapter.NewOrderServer(
		app.PreviewCheckoutExecutor,
		app.PlaceOrderExecutor,
		app.CancelOrderExecutor,
		app.ConfirmOrderExecutor,
	)

	mux.Handle(orderconnect.NewOrderCommandHandler(orderServer, allInterceptors))
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	return &AdapterModule{Mux: mux}
}
