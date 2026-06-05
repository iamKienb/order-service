package order

import (
	"time"

	"order-command-module/internal/domain/shared"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusDelivered OrderStatus = "DELIVERED"
	StatusShipped   OrderStatus = "SHIPPED"

	StatusCompleted OrderStatus = "COMPLETED"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusFailed    OrderStatus = "FAILED"
)

type Order struct {
	ID             string
	ShopID         shared.ShopID
	BuyerID        shared.UserID
	IdempotencyKey string
	Status         OrderStatus

	ShippingName     string
	ShippingPhone    string
	ShippingAddress  string
	ShippingProvince string
	ShippingWard     string

	Note *string

	GrandTotal int64
	Currency   string

	CancelReason *string
	CancelledBy  *shared.UserID

	ConfirmedAt *time.Time
	DeliveredAt *time.Time
	ShippedAt   *time.Time
	CompletedAt *time.Time
	CancelledAt *time.Time
	FailedAt    *time.Time

	CreatedAt    time.Time
	OrderItems   []OrderItem
	OrderHistory []OrderHistory

	shared.EventEntity
}

func NewOrder(params NewOrderParams) *Order {
	orderID := generateOrderID()
	now := time.Now().UTC()

	return &Order{
		ID:             orderID,
		ShopID:         params.ShopID,
		BuyerID:        params.BuyerID,
		IdempotencyKey: params.IdempotencyKey,
		Status:         StatusPending,

		ShippingName:     params.ShippingName,
		ShippingPhone:    params.ShippingPhone,
		ShippingAddress:  params.ShippingAddress,
		ShippingProvince: params.ShippingProvince,
		ShippingWard:     params.ShippingWard,

		Note: params.Note,

		GrandTotal: params.GrandTotal,

		Currency:     params.Currency,
		CancelReason: nil,
		CancelledBy:  nil,

		ConfirmedAt: nil,
		DeliveredAt: nil,
		ShippedAt:   nil,
		CompletedAt: nil,
		CancelledAt: nil,
		FailedAt:    nil,

		CreatedAt: now,
	}
}

func (o *Order) MarkAsCreated() {
	items := make([]OrderCreatedEventItem, 0, len(o.OrderItems))
	for _, item := range o.OrderItems {
		items = append(items, OrderCreatedEventItem{
			InventoryID:  item.InventoryID.String(),
			SkuID:        item.SkuID.String(),
			SkuCode:      item.SkuCode,
			ProductID:    item.ProductID.String(),
			ProductName:  item.ProductName,
			Quantity:     item.Quantity,
			BasePrice:    item.BasePrice,
			ItemSubtotal: item.ItemSubtotal,
		})
	}

	o.AddEvent(OrderCreatedEvent{
		OrderID:    o.ID,
		ShopID:     o.ShopID.String(),
		BuyerID:    o.BuyerID.String(),
		Status:     o.Status,
		GrandTotal: o.GrandTotal,
		Currency:   o.Currency,
		Items:      items,
		CreatedAt:  o.CreatedAt,
	})
}

func (o *Order) AddOrderItem(params OrderItemsParams) {
	orderItemID := shared.NewID[shared.OrderItemID]()

	item := createOrderItem(NewOrderItemParams{
		OrderID:     o.ID,
		OrderItemID: orderItemID,
		ProductID:   params.ProductID,
		InventoryID: params.InventoryID,
		SkuID:       params.SkuID,
		ProductName: params.ProductName,
		SkuCode:     params.SkuCode,
		ImageUrl:    params.ImageUrl,
		Quantity:    params.Quantity,
		BasePrice:   params.BasePrice,
		Currency:    o.Currency,
		CreatedAt:   o.CreatedAt,
	})

	o.OrderItems = append(o.OrderItems, *item)
}

func (o *Order) Confirm(actorID shared.UserID) error {
	if o.Status != StatusPending {
		return ErrOrderInvalidStateTransition
	}

	oldStatus := o.Status
	now := time.Now().UTC()
	o.Status = StatusConfirmed
	o.ConfirmedAt = &now
	o.AddEvent(OrderConfirmedEvent{
		OrderID:     o.ID,
		ShopID:      o.ShopID.String(),
		Status:      o.Status,
		ConfirmedAt: now,
	})

	o.appendStateHistory(NewOrderHistoryParams{
		OrderID:    o.ID,
		FromStatus: &oldStatus,
		ToStatus:   StatusConfirmed,
		ChangedBy:  actorID,
		ActorType:  ActorShop,
		Reason:     "Staff duyệt đơn hàng",
		CreatedAt:  now,
	})

	return nil
}

func (o *Order) Cancel(params CancelParams) error {
	if o.Status == StatusCompleted || o.Status == StatusCancelled || o.Status == StatusDelivered {
		return ErrOrderCannotCancel
	}

	switch params.ActorType {
	case ActorBuyer:
		if o.Status == StatusConfirmed {
			return ErrOrderBuyerNotAllowed
		}

		if params.ActorID == nil {
			return ErrOrderActorIDRequired
		}
	case ActorShop:
		if o.Status == StatusShipped {
			return ErrOrderShopNotAllowed
		}

		if params.ActorID == nil {
			return ErrOrderActorIDRequired
		}
	case ActorSystem:
		if params.Reason == "" {
			params.Reason = "RESERVED_INVENTORY_NOT_SUCCESS"
		}
	default:
		return ErrOrderInvalidActorType
	}

	oldStatus := o.Status
	now := time.Now().UTC()
	o.Status = StatusCancelled
	o.CancelReason = &params.Reason
	o.CancelledAt = &now
	o.AddEvent(OrderCancelledEvent{
		OrderID:     o.ID,
		ShopID:      o.ShopID.String(),
		BuyerID:     o.BuyerID.String(),
		Status:      o.Status,
		Reason:      params.Reason,
		CancelledAt: now,
	})

	var finalActorID shared.UserID
	if params.ActorType == ActorSystem {
		o.CancelledBy = &shared.SystemID
		finalActorID = shared.SystemID
	} else {
		o.CancelledBy = params.ActorID
		finalActorID = *params.ActorID
	}

	o.appendStateHistory(NewOrderHistoryParams{
		OrderID:    o.ID,
		FromStatus: &oldStatus,
		ToStatus:   StatusCancelled,
		ChangedBy:  finalActorID,
		ActorType:  ActorType(params.ActorType),
		Reason:     params.Reason,
		CreatedAt:  now,
	})

	return nil
}

func (o *Order) appendStateHistory(params NewOrderHistoryParams) {
	orderHistoryID := shared.NewID[shared.OrderHistoryID]()
	params.ID = orderHistoryID
	history := NewOrderHistory(params)

	o.OrderHistory = append(o.OrderHistory, *history)
}

func (o *Order) Type() string {
	return "Order"
}

func (o *Order) FlushEvents() []shared.DomainEvent {
	events := o.Flush()
	o.ClearEvent()
	return events
}

func generateOrderID() string {
	datePart := time.Now().Format("060102")
	randomPart := shared.CryptoRandomString(6)
	return "ORD-" + datePart + "-" + randomPart
}
