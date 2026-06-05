package order

import (
	"time"

	"order-command-module/db/repository"
	domain_order "order-command-module/internal/domain/order"
	"order-command-module/internal/domain/shared"

	"github.com/iamKienb/go-core/postgres/conv"
	"github.com/jackc/pgx/v5/pgtype"
)

func toCreateOrderParams(order *domain_order.Order) repository.CreateOrderParams {
	return repository.CreateOrderParams{
		ID:               order.ID,
		ShopID:           conv.UUID(order.ShopID),
		BuyerID:          conv.UUID(order.BuyerID),
		IdempotencyKey:   order.IdempotencyKey,
		Status:           string(order.Status),
		ShippingName:     order.ShippingName,
		ShippingPhone:    order.ShippingPhone,
		ShippingAddress:  order.ShippingAddress,
		ShippingProvince: order.ShippingProvince,
		ShippingWard:     order.ShippingWard,
		Note:             stringValue(order.Note),
		GrandTotal:       order.GrandTotal,
		Currency:         order.Currency,
		CancelReason:     stringValue(order.CancelReason),
		CancelledBy:      optionalUUID(order.CancelledBy),
		ConfirmedAt:      conv.TimeStampZ(order.ConfirmedAt),
		DeliveredAt:      conv.TimeStampZ(order.DeliveredAt),
		ShippedAt:        conv.TimeStampZ(order.ShippedAt),
		CompletedAt:      conv.TimeStampZ(order.CompletedAt),
		CancelledAt:      conv.TimeStampZ(order.CancelledAt),
		FailedAt:         conv.TimeStampZ(order.FailedAt),
		CreatedAt:        conv.TimeStampZ(&order.CreatedAt),
	}
}

func toCreateOrderItemsParams(order *domain_order.Order) repository.CreateOrderItemsBatchParams {
	params := repository.CreateOrderItemsBatchParams{
		Ids:           make([]pgtype.UUID, 0, len(order.OrderItems)),
		OrderID:       order.ID,
		InventoryIds:  make([]pgtype.UUID, 0, len(order.OrderItems)),
		SkuIds:        make([]pgtype.UUID, 0, len(order.OrderItems)),
		SkuCodes:      make([]string, 0, len(order.OrderItems)),
		ProductIds:    make([]pgtype.UUID, 0, len(order.OrderItems)),
		ProductNames:  make([]string, 0, len(order.OrderItems)),
		ImageUrls:     make([]string, 0, len(order.OrderItems)),
		Quantities:    make([]int64, 0, len(order.OrderItems)),
		BasePrices:    make([]int64, 0, len(order.OrderItems)),
		ItemSubtotals: make([]int64, 0, len(order.OrderItems)),
		Currency:      order.Currency,
		CreatedAt:     conv.TimeStampZ(&order.CreatedAt),
	}

	for _, item := range order.OrderItems {
		params.Ids = append(params.Ids, conv.UUID(item.ID))
		params.InventoryIds = append(params.InventoryIds, conv.UUID(item.InventoryID))
		params.SkuIds = append(params.SkuIds, conv.UUID(item.SkuID))
		params.SkuCodes = append(params.SkuCodes, item.SkuCode)
		params.ProductIds = append(params.ProductIds, conv.UUID(item.ProductID))
		params.ProductNames = append(params.ProductNames, item.ProductName)
		params.ImageUrls = append(params.ImageUrls, item.ImageUrl)
		params.Quantities = append(params.Quantities, item.Quantity)
		params.BasePrices = append(params.BasePrices, item.BasePrice)
		params.ItemSubtotals = append(params.ItemSubtotals, item.ItemSubtotal)
	}

	return params
}

func toCreateOrderHistoryParams(order *domain_order.Order) repository.CreateOrderHistoryBatchParams {
	params := repository.CreateOrderHistoryBatchParams{
		Ids:          make([]pgtype.UUID, 0, len(order.OrderHistory)),
		OrderID:      order.ID,
		FromStatuses: make([]string, 0, len(order.OrderHistory)),
		ToStatuses:   make([]string, 0, len(order.OrderHistory)),
		ChangedBys:   make([]pgtype.UUID, 0, len(order.OrderHistory)),
		ActorTypes:   make([]string, 0, len(order.OrderHistory)),
		Reasons:      make([]string, 0, len(order.OrderHistory)),
		CreatedAt:    conv.TimeStampZ(&order.OrderHistory[0].CreatedAt),
	}

	for _, history := range order.OrderHistory {
		params.Ids = append(params.Ids, conv.UUID(history.ID))
		params.FromStatuses = append(params.FromStatuses, orderStatusValue(history.FromStatus))
		params.ToStatuses = append(params.ToStatuses, string(history.ToStatus))
		params.ChangedBys = append(params.ChangedBys, conv.UUID(history.ChangedBy))
		params.ActorTypes = append(params.ActorTypes, string(history.ActorType))
		params.Reasons = append(params.Reasons, history.Reason)
	}

	return params
}

func toUpdateOrderStatusParams(order *domain_order.Order, expectedStatus domain_order.OrderStatus) repository.UpdateOrderStatusParams {
	return repository.UpdateOrderStatusParams{
		ID:             order.ID,
		ExpectedStatus: string(expectedStatus),
		Status:         string(order.Status),
		CancelReason:   stringValue(order.CancelReason),
		CancelledBy:    optionalUUID(order.CancelledBy),
		ConfirmedAt:    conv.TimeStampZ(order.ConfirmedAt),
		DeliveredAt:    conv.TimeStampZ(order.DeliveredAt),
		ShippedAt:      conv.TimeStampZ(order.ShippedAt),
		CompletedAt:    conv.TimeStampZ(order.CompletedAt),
		CancelledAt:    conv.TimeStampZ(order.CancelledAt),
		FailedAt:       conv.TimeStampZ(order.FailedAt),
	}
}

func toDomainOrder(row repository.Order, items []repository.OrderItem, histories []repository.OrderHistory) *domain_order.Order {
	order := &domain_order.Order{
		ID:               row.ID,
		ShopID:           shared.ShopID(row.ShopID.Bytes),
		BuyerID:          shared.UserID(row.BuyerID.Bytes),
		IdempotencyKey:   row.IdempotencyKey.String,
		Status:           domain_order.OrderStatus(row.Status),
		ShippingName:     row.ShippingName,
		ShippingPhone:    row.ShippingPhone,
		ShippingAddress:  row.ShippingAddress,
		ShippingProvince: row.ShippingProvince,
		ShippingWard:     row.ShippingWard,
		Note:             textPointer(row.Note),
		GrandTotal:       row.GrandTotal,
		Currency:         row.Currency,
		CancelReason:     textPointer(row.CancelReason),
		CancelledBy:      userIDPointer(row.CancelledBy),
		ConfirmedAt:      timePointer(row.ConfirmedAt),
		DeliveredAt:      timePointer(row.DeliveredAt),
		ShippedAt:        timePointer(row.ShippedAt),
		CompletedAt:      timePointer(row.CompletedAt),
		CancelledAt:      timePointer(row.CancelledAt),
		FailedAt:         timePointer(row.FailedAt),
		CreatedAt:        row.CreatedAt.Time,
		OrderItems:       toDomainOrderItems(items),
		OrderHistory:     toDomainOrderHistory(histories),
	}
	return order
}

func toDomainOrderItems(rows []repository.OrderItem) []domain_order.OrderItem {
	items := make([]domain_order.OrderItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, domain_order.OrderItem{
			ID:           shared.OrderItemID(row.ID.Bytes),
			OrderID:      row.OrderID,
			InventoryID:  shared.InventoryID(row.InventoryID.Bytes),
			SkuID:        shared.SkuID(row.SkuID.Bytes),
			SkuCode:      row.SkuCode,
			ProductID:    shared.ProductID(row.ProductID.Bytes),
			ProductName:  row.ProductName,
			ImageUrl:     row.ImageUrl,
			Quantity:     row.Quantity,
			BasePrice:    row.BasePrice,
			ItemSubtotal: row.ItemSubtotal,
			Currency:     row.Currency,
			CreatedAt:    row.CreatedAt.Time,
		})
	}
	return items
}

func toDomainOrderHistory(rows []repository.OrderHistory) []domain_order.OrderHistory {
	histories := make([]domain_order.OrderHistory, 0, len(rows))
	for _, row := range rows {
		histories = append(histories, domain_order.OrderHistory{
			ID:         shared.OrderHistoryID(row.ID.Bytes),
			OrderID:    row.OrderID,
			FromStatus: orderStatusPointer(row.FromStatus),
			ToStatus:   domain_order.OrderStatus(row.ToStatus),
			ChangedBy:  shared.UserID(row.ChangedBy.Bytes),
			ActorType:  domain_order.ActorType(row.ActorType),
			Reason:     row.Reason,
			CreatedAt:  row.CreatedAt.Time,
		})
	}
	return histories
}

func optionalUUID[T ~[16]byte](value *T) pgtype.UUID {
	if value == nil {
		return pgtype.UUID{}
	}
	return conv.UUID(*value)
}

func userIDPointer(value pgtype.UUID) *shared.UserID {
	if !value.Valid {
		return nil
	}
	userID := shared.UserID(value.Bytes)
	return &userID
}

func timePointer(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func textPointer(value pgtype.Text) *string {
	if !value.Valid || value.String == "" {
		return nil
	}
	result := value.String
	return &result
}

func orderStatusPointer(value pgtype.Text) *domain_order.OrderStatus {
	if !value.Valid || value.String == "" {
		return nil
	}
	status := domain_order.OrderStatus(value.String)
	return &status
}

func orderStatusValue(value *domain_order.OrderStatus) string {
	if value == nil {
		return ""
	}
	return string(*value)
}
