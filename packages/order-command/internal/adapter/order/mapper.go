package order

import (
	"strings"

	"order-command-module/internal/application/commands/cancel_order"
	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/application/commands/preview_checkout"
	"order-command-module/internal/application/services/order/i18n"
	domain_order "order-command-module/internal/domain/order"
	"order-command-module/internal/domain/shared"

	orderpb "github.com/iamKienb/api-contract/gen/order"
	"github.com/iamKienb/go-core/app_error"
)

func ToPreviewCheckoutCommand(userID string, req *orderpb.PreviewCheckoutRequest) (preview_checkout.Command, error) {
	addressID, err := parseUserAddressID(req.GetAddressUserId())
	if err != nil {
		return preview_checkout.Command{}, err
	}

	buyerID, err := parseUserID(userID)
	if err != nil {
		return preview_checkout.Command{}, err
	}

	parsedShopID, err := parseShopID(req.GetShopId())
	if err != nil {
		return preview_checkout.Command{}, err
	}

	items, err := toPreviewItems(req.GetItems())
	if err != nil {
		return preview_checkout.Command{}, err
	}

	return preview_checkout.Command{
		ShopID:         parsedShopID,
		BuyerID:        buyerID,
		BuyerAddressID: addressID,
		Items:          items,
	}, nil
}

func ToPreviewCheckoutResponse(shopID string, result *preview_checkout.Result) *orderpb.PreviewCheckoutResponse {
	items := make([]*orderpb.PreviewItemDetails, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, &orderpb.PreviewItemDetails{
			SkuId:       item.SkuID.String(),
			InventoryId: item.InventoryID.String(),
			Name:        item.ProductName,
			ThumbUrl:    item.ImageURL,
			Price:       item.BasePrice,
			Quantity:    int32(item.Quantity),
		})
	}

	return &orderpb.PreviewCheckoutResponse{
		ShopId:           shopID,
		TotalItemPrice:   result.GrandTotal,
		TotalShippingFee: 0,
		GrandTotal:       result.GrandTotal,
		Items:            items,
	}
}

func ToPlaceOrderCommand(userID string, req *orderpb.PlaceOrderRequest) (place_order.Command, error) {
	addressID, err := parseUserAddressID(req.GetAddressUserId())
	if err != nil {
		return place_order.Command{}, err
	}

	buyerID, err := parseUserID(userID)
	if err != nil {
		return place_order.Command{}, err
	}

	parsedShopID, err := parseShopID(req.GetShopId())
	if err != nil {
		return place_order.Command{}, err
	}

	idempotencyKey := strings.TrimSpace(req.GetIdempotencyKey())
	if idempotencyKey == "" {
		return place_order.Command{}, app_error.New(app_error.KindValidation, "", i18n.MsgOrderIdempotencyMissing, domain_order.ErrOrderIdempotencyMissing)
	}

	items := make([]place_order.Item, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		if item.GetQuantity() <= 0 {
			return place_order.Command{}, app_error.New(app_error.KindValidation, "", i18n.MsgInventoryQuantityInvalid, nil)
		}

		skuID, err := parseSkuID(item.GetSkuId())
		if err != nil {
			return place_order.Command{}, err
		}

		items = append(items, place_order.Item{
			SkuID:     skuID,
			BasePrice: item.GetBasePrice(),
			Quantity:  int64(item.GetQuantity()),
		})
	}

	return place_order.Command{
		ShopID:         parsedShopID,
		BuyerID:        buyerID,
		BuyerAddressID: addressID,
		IdempotencyKey: idempotencyKey,
		Items:          items,
	}, nil
}

func ToPlaceOrderResponse(result *place_order.Result) *orderpb.PlaceOrderResponse {
	if result == nil {
		return &orderpb.PlaceOrderResponse{Success: false}
	}

	return &orderpb.PlaceOrderResponse{
		Success: true,
		OrderId: result.OrderID,
		Status:  result.Status,
	}
}

func ToCancelOrderCommand(userID string, req *orderpb.CancelOrderRequest) (cancel_order.Command, error) {
	actorID, err := parseUserID(userID)
	if err != nil {
		return cancel_order.Command{}, err
	}

	return cancel_order.Command{
		OrderID:   req.GetOrderId(),
		ActorID:   actorID,
		ActorType: domain_order.ActorBuyer,
		Reason:    req.GetReason(),
	}, nil
}

func ToCancelOrderResponse(result *cancel_order.Result) *orderpb.CancelOrderResponse {
	return &orderpb.CancelOrderResponse{
		OrderId: result.OrderID,
		Status:  result.Status,
		Message: "order cancelled",
	}
}

func ToConfirmOrderCommand(userID string, req *orderpb.ConfirmOrderRequest) (confirm_order.Command, error) {
	actorID, err := parseUserID(userID)
	if err != nil {
		return confirm_order.Command{}, err
	}

	shopID, err := parseShopID(req.GetShopId())
	if err != nil {
		return confirm_order.Command{}, err
	}

	return confirm_order.Command{
		OrderID: req.GetOrderId(),
		ShopID:  shopID,
		ActorID: actorID,
	}, nil
}

func ToConfirmOrderResponse(result *confirm_order.Result) *orderpb.ConfirmOrderResponse {
	return &orderpb.ConfirmOrderResponse{
		OrderId: result.OrderID,
		Status:  result.Status,
	}
}

func toPreviewItems(items []*orderpb.CheckoutItem) ([]preview_checkout.Item, error) {
	result := make([]preview_checkout.Item, 0, len(items))
	for _, item := range items {
		if item.GetQuantity() <= 0 {
			return nil, app_error.New(app_error.KindValidation, "", i18n.MsgInventoryQuantityInvalid, nil)
		}

		skuID, err := parseSkuID(item.GetSkuId())
		if err != nil {
			return nil, err
		}

		result = append(result, preview_checkout.Item{
			SkuID:    skuID,
			Quantity: int64(item.GetQuantity()),
		})
	}

	return result, nil
}

func parseUserID(value string) (shared.UserID, error) {
	parsed, err := shared.ParseToRawID[shared.UserID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "user_invalid", "invalid user id", err)
	}

	return parsed, nil
}

func parseUserAddressID(value string) (shared.UserAddressID, error) {
	parsed, err := shared.ParseToRawID[shared.UserAddressID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "user_address_invalid", "invalid user address id", err)
	}

	return parsed, nil
}

func parseShopID(value string) (shared.ShopID, error) {
	parsed, err := shared.ParseToRawID[shared.ShopID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "shop_invalid", "invalid shop id", err)
	}

	return parsed, nil
}

func parseSkuID(value string) (shared.SkuID, error) {
	parsed, err := shared.ParseToRawID[shared.SkuID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "sku_invalid", "invalid sku id", err)
	}

	return parsed, nil
}

func parseInventoryID(value string) (shared.InventoryID, error) {
	parsed, err := shared.ParseToRawID[shared.InventoryID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "inventory_invalid", "invalid inventory id", err)
	}

	return parsed, nil
}
