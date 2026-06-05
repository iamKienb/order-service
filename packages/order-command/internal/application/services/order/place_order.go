package order

import (
	"context"
	"errors"
	"strings"

	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/application/commands/preview_checkout"
	"order-command-module/internal/application/port"
	domain_order "order-command-module/internal/domain/order"
	domain_shared "order-command-module/internal/domain/shared"
)

const defaultCurrency = "VND"

func (s *orderService) PlaceOrder(ctx context.Context, cmd place_order.Command, checkoutCtx preview_checkout.CheckoutContext) (*place_order.Result, error) {
	idempotencyKey := strings.TrimSpace(cmd.IdempotencyKey)
	if idempotencyKey == "" {
		return nil, ErrOrderIdempotencyMissing
	}

	baseItems := make([]checkoutLineInput, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		baseItems = append(baseItems, checkoutLineInput{SkuID: item.SkuID, Quantity: item.Quantity})
	}
	normalizedItems, err := normalizeCheckoutItems(baseItems)
	if err != nil {
		return nil, err
	}

	existing, err := s.orderRepo.GetOrderByBuyerAndIdempotencyKey(ctx, cmd.BuyerID, idempotencyKey)
	if err == nil {
		if !matchesPlaceOrderRequest(existing, cmd.ShopID, normalizedItems) {
			return nil, domain_order.ErrOrderIdempotencyKeyConflict
		}
		return placeOrderResultFromOrder(existing), nil
	}
	if !errors.Is(err, domain_order.ErrOrderNotFound) {
		return nil, err
	}

	calcResult, err := s.calculateOrderPlacement(cmd.ShopID, normalizedItems, checkoutCtx.ProductSkus, checkoutCtx.SkuStocks)
	if err != nil {
		return nil, err
	}

	address := checkoutCtx.UserAddress
	newOrder := domain_order.NewOrder(domain_order.NewOrderParams{
		ShopID:           cmd.ShopID,
		BuyerID:          cmd.BuyerID,
		IdempotencyKey:   idempotencyKey,
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
		if errors.Is(err, domain_order.ErrOrderIdempotencyKeyConflict) {
			existing, readErr := s.orderRepo.GetOrderByBuyerAndIdempotencyKey(ctx, cmd.BuyerID, idempotencyKey)
			if readErr == nil {
				if !matchesPlaceOrderRequest(existing, cmd.ShopID, normalizedItems) {
					return nil, domain_order.ErrOrderIdempotencyKeyConflict
				}
				return placeOrderResultFromOrder(existing), nil
			}
		}
		return nil, err
	}

	return placeOrderResultFromOrder(newOrder), nil
}

func placeOrderResultFromOrder(order *domain_order.Order) *place_order.Result {
	reserveItems := make([]place_order.ReserveItem, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		reserveItems = append(reserveItems, place_order.ReserveItem{
			InventoryID: item.InventoryID.String(),
			SkuID:       item.SkuID.String(),
			Quantity:    item.Quantity,
		})
	}

	return &place_order.Result{OrderID: order.ID, Status: string(order.Status), ReserveItems: reserveItems}
}

func matchesPlaceOrderRequest(order *domain_order.Order, shopID domain_shared.ShopID, items []checkoutLineInput) bool {
	if order.ShopID != shopID {
		return false
	}

	expected := make(map[string]int64, len(items))
	for _, item := range items {
		expected[item.SkuID.String()] += item.Quantity
	}

	actual := make(map[string]int64, len(order.OrderItems))
	for _, item := range order.OrderItems {
		actual[item.SkuID.String()] += item.Quantity
	}

	if len(expected) != len(actual) {
		return false
	}
	for skuID, quantity := range expected {
		if actual[skuID] != quantity {
			return false
		}
	}
	return true
}
