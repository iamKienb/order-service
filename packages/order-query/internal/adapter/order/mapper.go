package order

import (
	"order-query-module/internal/application/service/models"

	api "github.com/iamKienb/api-contract/gen/order"
)

func ToOrderView(order *models.Order) *api.OrderView {
	if order == nil {
		return nil
	}

	return &api.OrderView{
		OrderId:     order.OrderID,
		ShopId:      order.ShopID,
		BuyerId:     order.BuyerID,
		Status:      order.Status,
		GrandTotal:  order.GrandTotal,
		Currency:    order.Currency,
		Reason:      order.Reason,
		CreatedAt:   order.CreatedAt,
		ConfirmedAt: order.ConfirmedAt,
		CancelledAt: order.CancelledAt,
		Items:       ToOrderItemViews(order.Items),
	}
}

func ToOrderViews(orders []models.Order) []*api.OrderView {
	views := make([]*api.OrderView, 0, len(orders))
	for i := range orders {
		views = append(views, ToOrderView(&orders[i]))
	}
	return views
}

func ToOrderItemViews(items []models.OrderItem) []*api.OrderItemView {
	views := make([]*api.OrderItemView, 0, len(items))
	for _, item := range items {
		views = append(views, &api.OrderItemView{
			InventoryId:  item.InventoryID,
			SkuId:        item.SkuID,
			SkuCode:      item.SkuCode,
			ProductId:    item.ProductID,
			ProductName:  item.ProductName,
			Quantity:     item.Quantity,
			BasePrice:    item.BasePrice,
			ItemSubtotal: item.ItemSubtotal,
		})
	}
	return views
}
