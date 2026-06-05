package grpc_client

import (
	"context"
	"net/http"

	"order-command-module/internal/application/port"
	"order-command-module/internal/domain/shared"

	"connectrpc.com/connect"
	"github.com/iamKienb/api-contract/gen/inventory"
	"github.com/iamKienb/api-contract/gen/inventory/inventoryconnect"
)

type inventoryClient struct {
	client inventoryconnect.InventoryCommandServiceClient
}

func NewInventoryClient(httpClient *http.Client, baseURL string) port.InventoryClient {
	return &inventoryClient{client: inventoryconnect.NewInventoryCommandServiceClient(httpClient, baseURL)}
}

func (c *inventoryClient) GetStockBySkuIDs(ctx context.Context, req port.GetStockBySkuIDsRequest) ([]*port.SkuStock, error) {
	resp, err := c.client.GetStockBySkuIDs(ctx, connect.NewRequest(&inventory.GetStockBySkuIDsRequest{
		ShopId: req.ShopID,
		SkuIds: req.SkuIDs,
	}))
	if err != nil {
		return nil, err
	}

	result := make([]*port.SkuStock, 0, len(resp.Msg.GetSkuStock()))
	for _, item := range resp.Msg.GetSkuStock() {
		skuID, err := shared.ParseToRawID[shared.SkuID](item.GetSkuId())
		if err != nil {
			return nil, err
		}
		inventoryID, err := shared.ParseToRawID[shared.InventoryID](item.GetInventoryId())
		if err != nil {
			return nil, err
		}

		result = append(result, &port.SkuStock{
			InventoryID:       inventoryID,
			SkuID:             skuID,
			AvailableQuantity: item.GetAvailableQuantity(),
		})
	}
	return result, nil
}

func (c *inventoryClient) ReserveStock(ctx context.Context, req port.ReserveStockRequest) error {
	items := make([]*inventory.StockReservationItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &inventory.StockReservationItem{
			InventoryId: item.InventoryID,
			SkuId:       item.SkuID,
			Quantity:    int32(item.Quantity),
		})
	}

	_, err := c.client.ReserveStock(ctx, connect.NewRequest(&inventory.ReserveStockRequest{
		ShopId:  req.ShopID,
		OrderId: req.OrderID,
		Items:   items,
	}))
	return err
}

func (c *inventoryClient) ReleaseStock(ctx context.Context, orderID string) error {
	_, err := c.client.ReleaseStock(ctx, connect.NewRequest(&inventory.ReleaseStockRequest{OrderId: orderID}))
	return err
}

func (c *inventoryClient) FulfillStock(ctx context.Context, orderID string) error {
	_, err := c.client.FulfillStock(ctx, connect.NewRequest(&inventory.FulfillStockRequest{OrderId: orderID}))
	return err
}
