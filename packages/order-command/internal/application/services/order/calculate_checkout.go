package order

import (
	"math"

	"order-command-module/internal/application/port"
	domain_shared "order-command-module/internal/domain/shared"
)

const productSkuStatusActive = "ACTIVE"

type checkoutLineInput struct {
	SkuID    domain_shared.SkuID
	Quantity int64
}

type checkoutLine struct {
	ProductID         domain_shared.ProductID
	SkuID             domain_shared.SkuID
	InventoryID       domain_shared.InventoryID
	ProductName       string
	SkuCode           string
	BasePrice         int64
	Quantity          int64
	SubTotal          int64
	ImageURL          string
	AvailableQuantity int64
}

type checkoutLinesResult struct {
	Lines      []checkoutLine
	GrandTotal int64
}

func (r checkoutLinesResult) HasUnavailableLine() bool {
	for _, line := range r.Lines {
		if line.AvailableQuantity < line.Quantity {
			return true
		}
	}
	return false
}

func (s *orderService) calculateCheckoutPreview(shopID domain_shared.ShopID, items []checkoutLineInput, productSkus []*port.ProductSkuDetail, skuStocks []*port.SkuStock) (*checkoutLinesResult, error) {
	return s.buildCheckoutLines(shopID, items, productSkus, skuStocks)
}

func (s *orderService) calculateOrderPlacement(shopID domain_shared.ShopID, items []checkoutLineInput, productSkus []*port.ProductSkuDetail, skuStocks []*port.SkuStock) (*checkoutLinesResult, error) {
	result, err := s.buildCheckoutLines(shopID, items, productSkus, skuStocks)
	if err != nil {
		return nil, err
	}
	if result.HasUnavailableLine() {
		return nil, ErrCheckoutItemUnavailable
	}
	return result, nil
}

func (s *orderService) buildCheckoutLines(shopID domain_shared.ShopID, items []checkoutLineInput, productSkus []*port.ProductSkuDetail, skuStocks []*port.SkuStock) (*checkoutLinesResult, error) {
	stockMap := make(map[domain_shared.SkuID]*port.SkuStock, len(skuStocks))
	for _, stock := range skuStocks {
		if _, exists := stockMap[stock.SkuID]; exists {
			return nil, ErrCheckoutItemInvalid
		}
		stockMap[stock.SkuID] = stock
	}

	productMap := make(map[domain_shared.SkuID]*port.ProductSkuDetail, len(productSkus))
	for _, productSku := range productSkus {
		if productSku.ShopID != shopID || productSku.Status != productSkuStatusActive || productSku.Price < 0 {
			return nil, ErrCheckoutItemInvalid
		}
		if _, exists := productMap[productSku.SkuID]; exists {
			return nil, ErrCheckoutItemInvalid
		}
		productMap[productSku.SkuID] = productSku
	}

	var grandTotal int64
	lines := make([]checkoutLine, 0, len(items))

	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrCheckoutItemInvalid
		}

		productInfo, exists := productMap[item.SkuID]
		if !exists {
			return nil, ErrCheckoutItemInvalid
		}

		stockInfo, exists := stockMap[item.SkuID]
		if !exists {
			return nil, ErrCheckoutItemInvalid
		}
		if productInfo.Price > math.MaxInt64/item.Quantity {
			return nil, ErrCheckoutItemInvalid
		}

		subtotal := productInfo.Price * item.Quantity
		if grandTotal > math.MaxInt64-subtotal {
			return nil, ErrCheckoutItemInvalid
		}
		lines = append(lines, checkoutLine{
			SkuID:             item.SkuID,
			ProductID:         productInfo.ProductID,
			InventoryID:       stockInfo.InventoryID,
			ProductName:       productInfo.ProductName,
			SkuCode:           productInfo.SkuCode,
			BasePrice:         productInfo.Price,
			Quantity:          item.Quantity,
			SubTotal:          subtotal,
			ImageURL:          productInfo.ImageUrl,
			AvailableQuantity: stockInfo.AvailableQuantity,
		})
		grandTotal += subtotal
	}

	return &checkoutLinesResult{Lines: lines, GrandTotal: grandTotal}, nil
}
