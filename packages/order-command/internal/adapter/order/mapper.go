package order

import (
	"strings"

	"order-command-module/internal/application/commands/cancel_order"
	"order-command-module/internal/application/commands/confirm_order"
	"order-command-module/internal/application/commands/place_order"
	"order-command-module/internal/application/commands/preview_checkout"
	domain_order "order-command-module/internal/domain/order"
	"order-command-module/internal/domain/shared"

	orderpb "github.com/iamKienb/api-contract/gen/order"
	"github.com/iamKienb/go-core/app_error"
)

const (
	errCodeBuyerInvalid       = "buyer_invalid"
	errCodeUserAddressInvalid = "user_address_invalid"
	errCodeShopInvalid        = "shop_invalid"
	errCodeSkuInvalid         = "sku_invalid"
	errCodeActorInvalid       = "actor_invalid"
	errCodeIdempotencyInvalid = "idempotency_key_invalid"

	errMsgBuyerInvalid       = "invalid buyer id"
	errMsgUserAddressInvalid = "invalid user address id"
	errMsgShopInvalid        = "invalid shop id"
	errMsgSkuInvalid         = "invalid sku id"
	errMsgActorInvalid       = "invalid actor id"
	errMsgIdempotencyInvalid = "invalid idempotency key"
)

func ToPreviewCheckoutCommand(req *orderpb.PreviewCheckoutRequest) (preview_checkout.Command, error) {
	buyerID, err := parseID[shared.UserID](req.GetBuyerId(), errCodeBuyerInvalid, errMsgBuyerInvalid)
	if err != nil {
		return preview_checkout.Command{}, err
	}

	addressID, err := parseID[shared.UserAddressID](req.GetAddressUserId(), errCodeUserAddressInvalid, errMsgUserAddressInvalid)
	if err != nil {
		return preview_checkout.Command{}, err
	}

	shopID, err := parseID[shared.ShopID](req.GetShopId(), errCodeShopInvalid, errMsgShopInvalid)
	if err != nil {
		return preview_checkout.Command{}, err
	}

	items, err := toPreviewItems(req.GetItems())
	if err != nil {
		return preview_checkout.Command{}, err
	}

	return preview_checkout.Command{
		ShopID:         shopID,
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

func ToPlaceOrderCommand(req *orderpb.PlaceOrderRequest) (place_order.Command, error) {
	buyerID, err := parseID[shared.UserID](req.GetBuyerId(), errCodeBuyerInvalid, errMsgBuyerInvalid)
	if err != nil {
		return place_order.Command{}, err
	}

	addressID, err := parseID[shared.UserAddressID](req.GetAddressUserId(), errCodeUserAddressInvalid, errMsgUserAddressInvalid)
	if err != nil {
		return place_order.Command{}, err
	}

	shopID, err := parseID[shared.ShopID](req.GetShopId(), errCodeShopInvalid, errMsgShopInvalid)
	if err != nil {
		return place_order.Command{}, err
	}

	idempotencyKey := strings.TrimSpace(req.GetIdempotencyKey())
	if idempotencyKey == "" {
		return place_order.Command{}, app_error.New(app_error.KindValidation, errCodeIdempotencyInvalid, errMsgIdempotencyInvalid, nil)
	}

	items := make([]place_order.Item, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		skuID, err := parseID[shared.SkuID](item.GetSkuId(), errCodeSkuInvalid, errMsgSkuInvalid)
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
		ShopID:         shopID,
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

func ToCancelOrderCommand(req *orderpb.CancelOrderRequest) (cancel_order.Command, error) {
	actorID, err := parseID[shared.UserID](req.GetActorId(), errCodeActorInvalid, errMsgActorInvalid)
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

func ToConfirmOrderCommand(req *orderpb.ConfirmOrderRequest) (confirm_order.Command, error) {
	shopID, err := parseID[shared.ShopID](req.GetShopId(), errCodeShopInvalid, errMsgShopInvalid)
	if err != nil {
		return confirm_order.Command{}, err
	}

	actorID, err := parseID[shared.UserID](req.GetActorId(), errCodeActorInvalid, errMsgActorInvalid)
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
		skuID, err := parseID[shared.SkuID](item.GetSkuId(), errCodeSkuInvalid, errMsgSkuInvalid)
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

func parseID[T ~[16]byte](value, code, message string) (T, error) {
	parsed, err := shared.ParseToRawID[T](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, code, message, err)
	}

	return parsed, nil
}
