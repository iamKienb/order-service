package module

import (
	"log/slog"
	"net/http"

	orderadapter "order-query-module/internal/adapter/order"

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
	reflector := grpcreflect.NewStaticReflector(orderconnect.OrderQueryName)
	orderQueryServer := orderadapter.NewQueryServer(
		app.GetOrderDetailExecutor,
		app.ListBuyerOrdersExecutor,
		app.ListShopOrdersExecutor,
		app.SearchOrdersExecutor,
	)

	mux.Handle(orderconnect.NewOrderQueryHandler(orderQueryServer, allInterceptors))
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	return &AdapterModule{Mux: mux}
}
