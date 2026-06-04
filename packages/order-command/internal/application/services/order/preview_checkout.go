package order

import (
	"context"
	"fmt"
	"order-command-module/internal/application/commands/preview_checkout"
)

func (s *orderService) PreviewCheckout(ctx context.Context, cmd preview_checkout.Command, checkoutCtx preview_checkout.CheckoutContext) (*preview_checkout.Result, error) {
	baseItems := make([]baseItem, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		baseItems = append(baseItems, baseItem{
			SkuID:    item.SkuID,
			Quantity: item.Quantity,
		})
	}

	calcResult, err := s.calculateCheckout(baseItems, checkoutCtx.ProductSkus, checkoutCtx.SkuStocks)
	if err != nil {
		return nil, s.wrapError(err)
	}

	for _, calItem := range calcResult.Items {
		if !calItem.IsAvailable {
			return nil, fmt.Errorf("sản phẩm %s đã hết hàng, không thể hoàn tất đặt hàng", calItem.ProductName)
		}
	}

	previewDetails := make([]preview_checkout.PreviewItemDetail, 0, len(cmd.Items))

	for _, calItem := range calcResult.Items {
		previewDetails = append(previewDetails, preview_checkout.PreviewItemDetail{
			ShopID:      cmd.ShopID,
			ProductID:   calItem.ProductID,
			ProductName: calItem.ProductName,
			SkuID:       calItem.SkuID,
			SkuCode:     calItem.SkuCode,
			SubTotal:    calItem.SubTotal,
			Quantity:    calItem.Quantity,
			ImageURL:    calItem.ImageUrl,
			InventoryID: calItem.InventoryID,
		})
	}

	return &preview_checkout.Result{
		GrandTotal: calcResult.GrandTotal,
		Items:      previewDetails,
	}, nil
}
