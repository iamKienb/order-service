package list_shop_orders

import (
	"context"

	"order-query-module/internal/application/service/models"
)

type Query struct {
	ShopID string
	Status string
	Page   models.Page
}

type Result = models.OrderPage

type Executor interface {
	Execute(context.Context, Query) (*Result, error)
}
