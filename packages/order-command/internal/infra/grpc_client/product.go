package grpc_client

import (
	"context"
	"net/http"

	"order-command-module/internal/application/port"
	"order-command-module/internal/domain/shared"

	"connectrpc.com/connect"
	"github.com/iamKienb/api-contract/gen/product"
	"github.com/iamKienb/api-contract/gen/product/productconnect"
)

type productClient struct {
	client productconnect.ProductQueryServiceClient
}

func NewProductClient(httpClient *http.Client, baseURL string) port.ProductClient {
	return &productClient{client: productconnect.NewProductQueryServiceClient(httpClient, baseURL)}
}

func (c *productClient) GetProductBySkuIDs(ctx context.Context, req port.GetProductBySkuIdsRequest) ([]*port.ProductSkuDetail, error) {
	resp, err := c.client.GetProductsBySkuIDs(ctx, connect.NewRequest(&product.GetProductsBySkuIDsRequest{
		ShopId: req.ShopID,
		SkuIds: req.SkuIDs,
	}))
	if err != nil {
		return nil, err
	}

	result := make([]*port.ProductSkuDetail, 0, len(resp.Msg.GetItems()))
	for _, item := range resp.Msg.GetItems() {
		skuID, err := shared.ParseToRawID[shared.SkuID](item.GetSkuId())
		if err != nil {
			return nil, err
		}
		productID, err := shared.ParseToRawID[shared.ProductID](item.GetProductId())
		if err != nil {
			return nil, err
		}
		shopID, err := shared.ParseToRawID[shared.ShopID](item.GetShopId())
		if err != nil {
			return nil, err
		}

		result = append(result, &port.ProductSkuDetail{
			SkuID:       skuID,
			ProductID:   productID,
			ShopID:      shopID,
			SkuCode:     item.GetSkuCode(),
			ProductName: item.GetProductName(),
			Price:       item.GetPrice(),
			ImageUrl:    item.GetImageUrl(),
			Status:      item.GetStatus(),
		})
	}
	return result, nil
}

func (c *productClient) GetPriceSkusByIDs(ctx context.Context, req port.GetPriceSkusByIDsRequest) ([]*port.SkuPriceItem, error) {
	productSkus, err := c.GetProductBySkuIDs(ctx, port.GetProductBySkuIdsRequest{ShopID: req.ShopID, SkuIDs: req.SkuIDs})
	if err != nil {
		return nil, err
	}
	result := make([]*port.SkuPriceItem, 0, len(productSkus))
	for _, sku := range productSkus {
		result = append(result, &port.SkuPriceItem{SkuID: sku.SkuID, Price: sku.Price})
	}
	return result, nil
}
