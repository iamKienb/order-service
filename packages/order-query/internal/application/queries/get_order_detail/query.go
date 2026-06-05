package get_order_detail

import (
	"context"

	"order-query-module/internal/application/service/models"
)

type Query struct {
	OrderID string
}

type Result struct {
	Order *models.Order
}

type Executor interface {
	Execute(context.Context, Query) (*Result, error)
}
