package module

import (
	"order-query-module/internal/application/queries/get_order_detail"
	"order-query-module/internal/application/queries/list_buyer_orders"
	"order-query-module/internal/application/queries/list_shop_orders"
	"order-query-module/internal/application/queries/search_orders"
	"order-query-module/internal/application/service"
)

type ApplicationModule struct {
	GetOrderDetailExecutor  get_order_detail.Executor
	ListBuyerOrdersExecutor list_buyer_orders.Executor
	ListShopOrdersExecutor  list_shop_orders.Executor
	SearchOrdersExecutor    search_orders.Executor
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	orderQueryService := service.NewQueryService(infra.ESService)

	return &ApplicationModule{
		GetOrderDetailExecutor:  get_order_detail.NewHandler(orderQueryService),
		ListBuyerOrdersExecutor: list_buyer_orders.NewHandler(orderQueryService),
		ListShopOrdersExecutor:  list_shop_orders.NewHandler(orderQueryService),
		SearchOrdersExecutor:    search_orders.NewHandler(orderQueryService),
	}
}
