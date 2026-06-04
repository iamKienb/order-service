package order

import (
	"order-command-module/internal/domain/shared"
	"time"
)

type ActorType string

const (
	ActorBuyer   = "BUYER"
	ActorShop    = "SHOP"
	ActorSystem  = "SYSTEM"
	SystemUserID = "SYSTEM"
)

type OrderHistory struct {
	ID         shared.OrderHistoryID
	OrderID    string
	FromStatus *OrderStatus
	ToStatus   OrderStatus
	ChangedBy  shared.UserID
	ActorType  ActorType
	Reason     string
	CreatedAt  time.Time
}

func NewOrderHistory(params NewOrderHistoryParams) *OrderHistory {
	return &OrderHistory{
		ID:         params.ID,
		OrderID:    params.OrderID,
		FromStatus: params.FromStatus,
		ToStatus:   params.ToStatus,
		ChangedBy:  params.ChangedBy,
		ActorType:  params.ActorType,
		Reason:     params.Reason,
		CreatedAt:  params.CreatedAt,
	}
}
