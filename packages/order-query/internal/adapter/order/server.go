package order

import (
	"context"

	"order-query-module/internal/application/queries/get_order_detail"
	"order-query-module/internal/application/queries/list_buyer_orders"
	"order-query-module/internal/application/queries/list_shop_orders"
	"order-query-module/internal/application/queries/search_orders"
	"order-query-module/internal/application/service/models"

	"connectrpc.com/connect"
	api "github.com/iamKienb/api-contract/gen/order"
	"github.com/iamKienb/api-contract/gen/order/orderconnect"
)

type queryServer struct {
	orderconnect.UnimplementedOrderQueryHandler
	getOrderDetailExecutor  get_order_detail.Executor
	listBuyerOrdersExecutor list_buyer_orders.Executor
	listShopOrdersExecutor  list_shop_orders.Executor
	searchOrdersExecutor    search_orders.Executor
}

func NewQueryServer(
	getOrderDetailExecutor get_order_detail.Executor,
	listBuyerOrdersExecutor list_buyer_orders.Executor,
	listShopOrdersExecutor list_shop_orders.Executor,
	searchOrdersExecutor search_orders.Executor,
) *queryServer {
	return &queryServer{
		getOrderDetailExecutor:  getOrderDetailExecutor,
		listBuyerOrdersExecutor: listBuyerOrdersExecutor,
		listShopOrdersExecutor:  listShopOrdersExecutor,
		searchOrdersExecutor:    searchOrdersExecutor,
	}
}

func (s *queryServer) GetOrderDetail(ctx context.Context, req *connect.Request[api.GetOrderDetailRequest]) (*connect.Response[api.GetOrderDetailResponse], error) {
	result, err := s.getOrderDetailExecutor.Execute(ctx, get_order_detail.Query{
		OrderID: req.Msg.GetOrderId(),
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.GetOrderDetailResponse{Order: ToOrderView(result.Order)}), nil
}

func (s *queryServer) ListBuyerOrders(ctx context.Context, req *connect.Request[api.ListBuyerOrdersRequest]) (*connect.Response[api.ListBuyerOrdersResponse], error) {
	result, err := s.listBuyerOrdersExecutor.Execute(ctx, list_buyer_orders.Query{
		BuyerID: req.Msg.GetBuyerId(),
		Status:  req.Msg.GetStatus(),
		Page: models.Page{
			Size:  int(req.Msg.GetPageSize()),
			Token: req.Msg.GetPageToken(),
		},
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.ListBuyerOrdersResponse{
		Orders:        ToOrderViews(result.Items),
		Total:         result.Total,
		NextPageToken: result.NextPageToken,
	}), nil
}

func (s *queryServer) ListShopOrders(ctx context.Context, req *connect.Request[api.ListShopOrdersRequest]) (*connect.Response[api.ListShopOrdersResponse], error) {
	result, err := s.listShopOrdersExecutor.Execute(ctx, list_shop_orders.Query{
		ShopID: req.Msg.GetShopId(),
		Status: req.Msg.GetStatus(),
		Page: models.Page{
			Size:  int(req.Msg.GetPageSize()),
			Token: req.Msg.GetPageToken(),
		},
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.ListShopOrdersResponse{
		Orders:        ToOrderViews(result.Items),
		Total:         result.Total,
		NextPageToken: result.NextPageToken,
	}), nil
}

func (s *queryServer) SearchOrders(ctx context.Context, req *connect.Request[api.SearchOrdersRequest]) (*connect.Response[api.SearchOrdersResponse], error) {
	result, err := s.searchOrdersExecutor.Execute(ctx, search_orders.Query{
		ShopID:  req.Msg.GetShopId(),
		BuyerID: req.Msg.GetBuyerId(),
		Status:  req.Msg.GetStatus(),
		Keyword: req.Msg.GetKeyword(),
		Page: models.Page{
			Size:  int(req.Msg.GetPageSize()),
			Token: req.Msg.GetPageToken(),
		},
	})
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&api.SearchOrdersResponse{
		Orders:        ToOrderViews(result.Items),
		Total:         result.Total,
		NextPageToken: result.NextPageToken,
	}), nil
}

var _ orderconnect.OrderQueryHandler = (*queryServer)(nil)
