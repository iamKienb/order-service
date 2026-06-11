package activity

import (
	"context"
	"errors"

	"order-command-module/internal/application/commands/cancel_order"
	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/application/commands/preview_checkout"
	"order-command-module/internal/application/port"
	app_order "order-command-module/internal/application/services/order"
	domain_order "order-command-module/internal/domain/order"
	"order-command-module/internal/domain/shared"

	"golang.org/x/sync/errgroup"
)

const (
	reserveStockFailedReason = "RESERVE_STOCK_FAILED"
	FulfillStockFailedReason = "FULFILL_STOCK_FAILED"
)

type OrderActivity struct {
	service         app_order.Service
	userClient      port.UserClient
	productClient   port.ProductClient
	inventoryClient port.InventoryClient
}

type ReserveStockCommand struct {
	ShopID  string
	OrderID string
	BuyerID string
	Items   []place_order.ReserveItem
}

type SystemCancelOrderCommand struct {
	OrderID string
	Reason  string
}

func NewOrderActivity(
	service app_order.Service,
	userClient port.UserClient,
	productClient port.ProductClient,
	inventoryClient port.InventoryClient,
) *OrderActivity {
	return &OrderActivity{
		service:         service,
		userClient:      userClient,
		productClient:   productClient,
		inventoryClient: inventoryClient,
	}
}

func (a *OrderActivity) CreateOrder(ctx context.Context, cmd place_order.Command) (*place_order.Result, error) {
	existing, err := a.service.FindExistingPlaceOrder(ctx, cmd)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, domain_order.ErrOrderNotFound) {
		return nil, err
	}

	checkoutCtx, err := a.checkoutContext(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return a.service.PlaceOrder(ctx, cmd, *checkoutCtx)
}

func (a *OrderActivity) ConfirmOrder(ctx context.Context, cmd confirm_order.Command) (*confirm_order.Result, error) {
	return a.service.ConfirmOrder(ctx, cmd)
}

func (a *OrderActivity) CancelOrder(ctx context.Context, cmd cancel_order.Command) (*cancel_order.Result, error) {
	return a.service.CancelOrder(ctx, cmd)
}

func (a *OrderActivity) CancelPlacedOrder(ctx context.Context, orderID string) error {
	return a.CancelOrderBySystem(ctx, SystemCancelOrderCommand{
		OrderID: orderID,
		Reason:  reserveStockFailedReason,
	})
}

func (a *OrderActivity) CancelOrderBySystem(ctx context.Context, cmd SystemCancelOrderCommand) error {
	reason := cmd.Reason
	if reason == "" {
		reason = reserveStockFailedReason
	}

	_, err := a.service.CancelOrder(ctx, cancel_order.Command{
		OrderID:   cmd.OrderID,
		ActorID:   shared.SystemID,
		ActorType: domain_order.ActorSystem,
		Reason:    reason,
	})
	return err
}

func (a *OrderActivity) ReserveStock(ctx context.Context, cmd ReserveStockCommand) error {
	items := make([]port.ReserveStockItem, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		items = append(items, port.ReserveStockItem{
			InventoryID: item.InventoryID,
			SkuID:       item.SkuID,
			Quantity:    item.Quantity,
		})
	}

	return a.inventoryClient.ReserveStock(ctx, port.ReserveStockRequest{
		ShopID:  cmd.ShopID,
		OrderID: cmd.OrderID,
		BuyerID: cmd.BuyerID,
		Items:   items,
	})
}

func (a *OrderActivity) ReleaseStock(ctx context.Context, param port.ReleaseAndFullfilStockParam) error {
	return a.inventoryClient.ReleaseStock(ctx, param)
}

func (a *OrderActivity) FulfillStock(ctx context.Context, param port.ReleaseAndFullfilStockParam) error {
	return a.inventoryClient.FulfillStock(ctx, param)
}

func (a *OrderActivity) checkoutContext(ctx context.Context, cmd place_order.Command) (*preview_checkout.CheckoutContext, error) {
	group, ctx := errgroup.WithContext(ctx)

	var productSkus []*port.ProductSkuDetail
	var skuStocks []*port.SkuStock
	var userAddress *port.UserAddress

	skuIDs := make([]string, 0, len(cmd.Items))
	skuSeen := make(map[string]struct{}, len(cmd.Items))
	for _, item := range cmd.Items {
		skuID := item.SkuID.String()
		if _, exists := skuSeen[skuID]; exists {
			continue
		}
		skuSeen[skuID] = struct{}{}
		skuIDs = append(skuIDs, skuID)
	}

	group.Go(func() error {
		var err error
		productSkus, err = a.productClient.GetProductBySkuIDs(ctx, port.GetProductBySkuIdsRequest{
			ShopID: cmd.ShopID.String(),
			SkuIDs: skuIDs,
		})
		return err
	})

	group.Go(func() error {
		var err error
		skuStocks, err = a.inventoryClient.GetStockBySkuIDs(ctx, port.GetStockBySkuIDsRequest{
			ShopID: cmd.ShopID.String(),
			SkuIDs: skuIDs,
		})
		return err
	})

	group.Go(func() error {
		var err error
		userAddress, err = a.userClient.GetAddressByID(ctx, port.GetAddressByIDRequest{
			UserID:        cmd.BuyerID.String(),
			UserAddressID: cmd.BuyerAddressID.String(),
		})
		return err
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return &preview_checkout.CheckoutContext{
		ProductSkus: productSkus,
		SkuStocks:   skuStocks,
		UserAddress: userAddress,
	}, nil
}
