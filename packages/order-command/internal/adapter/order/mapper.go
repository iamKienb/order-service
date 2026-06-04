package inventory

import (
	"inventory-command-module/internal/application/commands/create_inventories"
	"inventory-command-module/internal/application/commands/delete_inventories"
	"inventory-command-module/internal/application/commands/fulfill_stock"
	"inventory-command-module/internal/application/commands/release_stock"
	"inventory-command-module/internal/application/commands/reserve_stock"
	"inventory-command-module/internal/domain/shared"

	"github.com/iamKienb/api-contract/gen/inventory"
	"github.com/iamKienb/go-core/app_error"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCreateInventoriesCommand(userID string, req *inventory.CreateInventoriesRequest) (create_inventories.Command, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return create_inventories.Command{}, err
	}

	parsedShopID, err := parseShopID(req.GetShopId())
	if err != nil {
		return create_inventories.Command{}, err
	}

	items := make([]create_inventories.Item, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		skuID, err := parseSkuID(item.GetSkuId())
		if err != nil {
			return create_inventories.Command{}, err
		}

		items = append(items, create_inventories.Item{
			SkuID:    skuID,
			Quantity: int64(item.GetQuantity()),
		})
	}

	return create_inventories.Command{
		ShopID:    parsedShopID,
		ActorID:   parsedUserID,
		Inventory: items,
	}, nil
}

func ToCreateInventoriesResponse(result *create_inventories.Result) *inventory.CreateInventoriesResponse {
	return &inventory.CreateInventoriesResponse{Success: result.Success}
}

func ToDeleteInventoriesCommand(userID string, req *inventory.DeleteInventoryRequest) (delete_inventories.Command, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return delete_inventories.Command{}, err
	}

	skuIDs := make([]shared.SkuID, 0, len(req.GetSkuIds()))
	for _, rawSkuID := range req.GetSkuIds() {
		skuID, err := parseSkuID(rawSkuID)
		if err != nil {
			return delete_inventories.Command{}, err
		}

		skuIDs = append(skuIDs, skuID)
	}

	return delete_inventories.Command{
		ActorID: parsedUserID,
		SkuIDs:  skuIDs,
	}, nil
}

func ToDeleteInventoriesResponse(result *delete_inventories.Result) *inventory.DeleteInventoryResponse {
	return &inventory.DeleteInventoryResponse{Success: result.Success}
}

func ToReserveStockCommand(userID string, req *inventory.ReserveStockRequest) (reserve_stock.Command, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return reserve_stock.Command{}, err
	}

	parsedShopID, err := parseShopID(req.GetShopId())
	if err != nil {
		return reserve_stock.Command{}, err
	}

	items := make([]reserve_stock.Item, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		skuID, err := parseSkuID(item.GetSkuId())
		if err != nil {
			return reserve_stock.Command{}, err
		}

		items = append(items, reserve_stock.Item{
			SkuID:    skuID,
			Quantity: int64(item.GetQuantity()),
		})
	}

	return reserve_stock.Command{
		ShopID:  parsedShopID,
		ActorID: parsedUserID,
		OrderID: req.GetOrderId(),
		Items:   items,
	}, nil
}

func ToReserveStockResponse(result *reserve_stock.Result) *inventory.ReserveStockResponse {
	return &inventory.ReserveStockResponse{ExpiresAt: timestamppb.New(result.ExpiresAt)}
}

func ToReleaseStockCommand(userID string, req *inventory.ReleaseStockRequest) (release_stock.Command, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return release_stock.Command{}, err
	}

	return release_stock.Command{
		OrderID: req.GetOrderId(),
		ActorID: parsedUserID,
	}, nil
}

func ToReleaseStockResponse(result *release_stock.Result) *inventory.ReleaseStockResponse {
	return &inventory.ReleaseStockResponse{Success: result.Success}
}

func ToFulfillStockCommand(userID string, req *inventory.FulfillStockRequest) (fulfill_stock.Command, error) {
	parsedUserID, err := parseUserID(userID)
	if err != nil {
		return fulfill_stock.Command{}, err
	}

	return fulfill_stock.Command{
		OrderID: req.GetOrderId(),
		ActorID: parsedUserID,
	}, nil
}

func ToAFulfilStockResponse(result *fulfill_stock.Result) *inventory.FulfillStockResponse {
	return &inventory.FulfillStockResponse{Success: result.Success}
}

func parseUserID(value string) (shared.UserID, error) {
	parsed, err := shared.ParseToRawID[shared.UserID](value)
	if err != nil {
		return parsed, app_error.New(app_error.KindValidation, "user_invalid", "invalid user id", err)
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
