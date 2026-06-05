package search_orders

import (
	"context"

	"order-query-module/internal/application/service/models"
)

type Query struct {
	ShopID  string
	BuyerID string
	Status  string
	Keyword string
	Page    models.Page
}

type Result = models.OrderPage

type Executor interface {
	Execute(context.Context, Query) (*Result, error)
}
