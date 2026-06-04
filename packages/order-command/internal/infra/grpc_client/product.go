package grpc_client

import (
	"order-command-module/internal/application/port"
	"order-command-module/internal/domain/shared"

	"github.com/iamKienb/api-contract/gen/product/productconnect"
)

type productClient struct {
	client productconnect.ProductQueryServiceClient
}

func NewProductClient(client productconnect.ProductQueryServiceClient) port.ProductClient {
	return &productClient{client: client}
}

func (c *productClient) GetProductBySkuIDs(ctx context.Context, req port.GetProductBySkuIdsRequest) ([]*port.ProductSkuDetail, error) {
	grpcReq := connect.NewRequest(&product.GetProductBySkuIdsRequest{
		ShopId: req.ShopID,
		SkuIds: req.SkuIDs,
	})

	resp, err := c.client.GetProductsBySkuIDs(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	result := make([]*port.ProductSkuDetail, 0, len(resp.Msg.Items))
	for _, p := range resp.Msg.Items {
		skuID, err := shared.ParseToRawID[shared.SkuID](p.SkuId)
		if err != nil {
			return nil, err
		}

		productID, err := shared.ParseToRawID[shared.ProductID](p.ProductId)
		if err != nil {
			return nil, err
		}

		shopID, err := shared.ParseToRawID[shared.ShopID](p.ShopId)
		if err != nil {
			return nil, err
		}

		result = append(result, &port.ProductSkuDetail{
			SkuID:       skuID,
			ProductID:   productID,
			ShopID:      shopID,
			SkuCode:     p.SkuCode,
			ProductName: p.ProductName,
			Price:       p.Price,
			ImageUrl:    p.ImageUrl,
			Status:      p.Status,
		})
	}

	return result, nil
}
