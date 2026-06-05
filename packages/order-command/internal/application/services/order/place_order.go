package order

import (
	"context"

	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/application/commands/preview_checkout"
	"order-command-module/internal/application/port"
	domain_order "order-command-module/internal/domain/order"
)

const defaultCurrency = "VND"

func (s *orderService) PlaceOrder(ctx context.Context, cmd place_order.Command, checkoutCtx preview_checkout.CheckoutContext) (*place_order.Result, error) {
	baseItems := make([]checkoutLineInput, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		baseItems = append(baseItems, checkoutLineInput{SkuID: item.SkuID, Quantity: item.Quantity})
	}

	calcResult, err := s.calculateOrderPlacement(baseItems, checkoutCtx.ProductSkus, checkoutCtx.SkuStocks)
	if err != nil {
		return nil, err
	}

	address := checkoutCtx.UserAddress
	newOrder := domain_order.NewOrder(domain_order.NewOrderParams{
		ShopID:           cmd.ShopID,
		BuyerID:          cmd.BuyerID,
		ShippingName:     address.ReceiverName,
		ShippingPhone:    address.PhoneNumber,
		ShippingAddress:  address.AddressLine,
		ShippingProvince: address.ProvinceName,
		ShippingWard:     address.WardName,
		GrandTotal:       calcResult.GrandTotal,
		Currency:         defaultCurrency,
	})

	for _, item := range calcResult.Lines {
		newOrder.AddOrderItem(domain_order.OrderItemsParams{
			ProductID:   item.ProductID,
			InventoryID: item.InventoryID,
			SkuID:       item.SkuID,
			ProductName: item.ProductName,
			SkuCode:     item.SkuCode,
			ImageUrl:    item.ImageURL,
			Quantity:    item.Quantity,
			BasePrice:   item.BasePrice,
		})
	}
	newOrder.MarkAsCreated()

	if err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		if err := s.orderRepo.CreateOrder(ctx, newOrder); err != nil {
			return err
		}
		if events := newOrder.FlushEvents(); len(events) > 0 {
			return s.outboxService.PublishBatch(ctx, []port.OutboxParam{{
				AggregateID:   orderAggregateID(newOrder.ID),
				AggregateType: newOrder.Type(),
				Events:        events,
			}})
		}
		return nil
	}); err != nil {
		return nil, err
	}

	reserveItems := make([]place_order.ReserveItem, 0, len(calcResult.Lines))
	for _, item := range calcResult.Lines {
		reserveItems = append(reserveItems, place_order.ReserveItem{
			InventoryID: item.InventoryID.String(),
			SkuID:       item.SkuID.String(),
			Quantity:    item.Quantity,
		})
	}

	return &place_order.Result{OrderID: newOrder.ID, Status: string(newOrder.Status), ReserveItems: reserveItems}, nil
}
