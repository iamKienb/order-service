package port

import (
	"context"
	"order-command-module/internal/domain/shared"
)

type GetProductBySkuIdsRequest struct {
	ShopID string
	SkuIDs []string
}

type ProductSkuDetail struct {
	SkuID       shared.SkuID
	ProductID   shared.ProductID
	ShopID      shared.ShopID
	SkuCode     string
	ProductName string
	Price       int64
	ImageUrl    string
	Status      string
}

type SkuPriceItem struct {
	SkuID shared.SkuID
	Price int64
}

type GetPriceSkusByIDsRequest struct {
	ShopID string
	SkuIDs []string
}

type ProductClient interface {
	GetProductBySkuIDs(ctx context.Context, req GetProductBySkuIdsRequest) ([]*ProductSkuDetail, error)
	GetPriceSkusByIDs(ctx context.Context, req GetPriceSkusByIDsRequest) ([]*SkuPriceItem, error)
}
