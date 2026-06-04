package order

import (
	"fmt"
	"order-command-module/internal/application/port"
	domain_shared "order-command-module/internal/domain/shared"
)

type baseItem struct {
	SkuID    domain_shared.SkuID
	Quantity int64
}

type calculatedItem struct {
	ProductID   domain_shared.ProductID
	SkuID       domain_shared.SkuID
	InventoryID domain_shared.InventoryID
	ProductName string
	SkuCode     string
	BasePrice   int64
	Quantity    int64
	SubTotal    int64
	ImageUrl    string
	IsAvailable bool
}

type calculatedCheckoutResult struct {
	Items      []calculatedItem
	GrandTotal int64
}

func (s *orderService) calculateCheckout(items []baseItem, productSkus []*port.ProductSkuDetail, skuStocks []*port.SkuStock) (*calculatedCheckoutResult, error) {
	stockMap := make(map[domain_shared.SkuID]*port.SkuStock)
	for _, skuStock := range skuStocks {
		stockMap[skuStock.SkuID] = skuStock
	}

	productMap := make(map[domain_shared.SkuID]*port.ProductSkuDetail)
	for _, productSku := range productSkus {
		productMap[productSku.SkuID] = productSku
	}

	var grandTotal int64
	calculatedItems := make([]calculatedItem, 0, len(items))

	for _, item := range items {
		productInfo, exists := productMap[item.SkuID]
		if !exists {
			return nil, fmt.Errorf("sản phẩm không hợp lệ")
		}

		stockInfo, exists := stockMap[item.SkuID]
		if !exists {
			return nil, fmt.Errorf("sản phẩm không hợp lệ")
		}

		isAvailable := stockInfo.AvailableQuantity >= item.Quantity
		subtotal := productInfo.Price * item.Quantity

		calculatedItems = append(calculatedItems, calculatedItem{
			SkuID:       item.SkuID,
			ProductID:   productInfo.ProductID,
			InventoryID: stockInfo.InventoryID,
			ProductName: productInfo.ProductName,
			SkuCode:     productInfo.SkuCode,
			BasePrice:   productInfo.Price,
			Quantity:    item.Quantity,
			SubTotal:    subtotal,
			ImageUrl:    productInfo.ImageUrl,

			IsAvailable: isAvailable,
		})

		grandTotal += subtotal

	}

	return &calculatedCheckoutResult{
		Items:      calculatedItems,
		GrandTotal: grandTotal,
	}, nil
}
