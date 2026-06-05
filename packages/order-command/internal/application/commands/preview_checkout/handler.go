package preview_checkout

import (
	"context"
	"order-command-module/internal/application/port"

	"golang.org/x/sync/errgroup"
)

type CheckoutContext struct {
	ProductSkus []*port.ProductSkuDetail
	SkuStocks   []*port.SkuStock
	UserAddress *port.UserAddress
}

type orderService interface {
	PreviewCheckout(ctx context.Context, cmd Command, checkoutCtx CheckoutContext) (*Result, error)
}

type handler struct {
	userClient      port.UserClient
	productClient   port.ProductClient
	inventoryClient port.InventoryClient
	service         orderService
}

func NewHandler(
	service orderService,
	userClient port.UserClient,
	productClient port.ProductClient,
	inventoryClient port.InventoryClient,
) Executor {
	return &handler{
		service:         service,
		userClient:      userClient,
		productClient:   productClient,
		inventoryClient: inventoryClient,
	}
}

func (h *handler) Execute(ctx context.Context, cmd Command) (*Result, error) {
	group, ctx := errgroup.WithContext(ctx)

	var productSkus []*port.ProductSkuDetail
	var skuStocks []*port.SkuStock
	var userAddress *port.UserAddress

	skuIDs := make([]string, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		skuIDs = append(skuIDs, item.SkuID.String())
	}

	group.Go(func() error {
		var err error
		productSkus, err = h.productClient.GetProductBySkuIDs(ctx, port.GetProductBySkuIdsRequest{
			ShopID: cmd.ShopID.String(),
			SkuIDs: skuIDs,
		})
		return err
	})

	group.Go(func() error {
		var err error
		skuStocks, err = h.inventoryClient.GetStockBySkuIDs(ctx, port.GetStockBySkuIDsRequest{
			ShopID: cmd.ShopID.String(),
			SkuIDs: skuIDs,
		})
		return err
	})

	group.Go(func() error {
		var err error
		userAddress, err = h.userClient.GetAddressByID(ctx, port.GetAddressByIDRequest{
			UserID:        cmd.BuyerID.String(),
			UserAddressID: cmd.BuyerAddressID.String(),
		})
		return err
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	checkoutCtx := CheckoutContext{
		ProductSkus: productSkus,
		SkuStocks:   skuStocks,
		UserAddress: userAddress,
	}

	return h.service.PreviewCheckout(ctx, cmd, checkoutCtx)
}
