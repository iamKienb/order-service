package list_buyer_orders

import (
	"context"

	"order-query-module/internal/application/service/models"
)

type Query struct {
	BuyerID string
	Status  string
	Page    models.Page
}

type Result = models.OrderPage

type Executor interface {
	Execute(context.Context, Query) (*Result, error)
}
