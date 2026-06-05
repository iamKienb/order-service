package order

import (
	"context"
	"fmt"

	"order-command-module/internal/application/commands/preview_checkout"
)

func (s *orderService) PreviewCheckout(ctx context.Context, cmd preview_checkout.Command, checkoutCtx preview_checkout.CheckoutContext) (*preview_checkout.Result, error) {
	baseItems := make([]checkoutLineInput, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		baseItems = append(baseItems, checkoutLineInput{SkuID: item.SkuID, Quantity: item.Quantity})
	}
	normalizedItems, err := normalizeCheckoutItems(baseItems)
	if err != nil {
		return nil, err
	}

	calcResult, err := s.calculateCheckoutPreview(normalizedItems, checkoutCtx.ProductSkus, checkoutCtx.SkuStocks)
	if err != nil {
		return nil, err
	}

	for _, item := range calcResult.Lines {
		if item.AvailableQuantity < item.Quantity {
			return nil, fmt.Errorf("%w: %s", ErrCheckoutItemUnavailable, item.ProductName)
		}
	}

	previewDetails := make([]preview_checkout.PreviewItemDetail, 0, len(calcResult.Lines))
	for _, item := range calcResult.Lines {
		previewDetails = append(previewDetails, preview_checkout.PreviewItemDetail{
			ShopID:      cmd.ShopID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			SkuID:       item.SkuID,
			SkuCode:     item.SkuCode,
			BasePrice:   item.BasePrice,
			SubTotal:    item.SubTotal,
			Quantity:    item.Quantity,
			ImageURL:    item.ImageURL,
			InventoryID: item.InventoryID,
		})
	}

	return &preview_checkout.Result{
		GrandTotal: calcResult.GrandTotal,
		Items:      previewDetails,
	}, nil
}
