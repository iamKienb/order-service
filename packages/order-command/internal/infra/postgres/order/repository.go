package order

import (
	"context"
	"errors"
	"fmt"

	"order-command-module/db/repository"
	domain_order "order-command-module/internal/domain/order"
	"order-command-module/internal/domain/shared"

	pgx "github.com/iamKienb/go-core/postgres"
	"github.com/iamKienb/go-core/postgres/conv"
	pgxv5 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const uniqueOrderBuyerIdempotencyKey = "idx_orders_buyer_idempotency_key"

type orderRepository struct {
	queries *repository.Queries
}

func NewRepository(service pgx.PGXService) domain_order.Repository {
	return &orderRepository{
		queries: repository.New(service.GetPool()),
	}
}

func (r *orderRepository) getQuerier(ctx context.Context) *repository.Queries {
	if tx := pgx.ExtractTx(ctx); tx != nil {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *domain_order.Order) error {
	if err := r.getQuerier(ctx).CreateOrder(ctx, toCreateOrderParams(order)); err != nil {
		if isUniqueConstraint(err, uniqueOrderBuyerIdempotencyKey) {
			return domain_order.ErrOrderIdempotencyKeyConflict
		}
		return fmt.Errorf("infra: create order: %w", err)
	}

	if err := r.getQuerier(ctx).CreateOrderItemsBatch(ctx, toCreateOrderItemsParams(order)); err != nil {
		return fmt.Errorf("infra: create order items: %w", err)
	}

	if len(order.OrderHistory) > 0 {
		if err := r.getQuerier(ctx).CreateOrderHistoryBatch(ctx, toCreateOrderHistoryParams(order)); err != nil {
			return fmt.Errorf("infra: create order history: %w", err)
		}
	}

	return nil
}

func (r *orderRepository) GetOrderByID(ctx context.Context, orderID string) (*domain_order.Order, error) {
	row, err := r.getQuerier(ctx).GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			return nil, domain_order.ErrOrderNotFound
		}
		return nil, fmt.Errorf("infra: get order by id: %w", err)
	}

	items, err := r.getQuerier(ctx).ListOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("infra: list order items: %w", err)
	}

	histories, err := r.getQuerier(ctx).ListOrderHistoryByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("infra: list order history: %w", err)
	}

	return toDomainOrder(row, items, histories), nil
}

func (r *orderRepository) GetOrderByBuyerAndIdempotencyKey(ctx context.Context, buyerID shared.UserID, idempotencyKey string) (*domain_order.Order, error) {
	row, err := r.getQuerier(ctx).GetOrderByBuyerAndIdempotencyKey(ctx, repository.GetOrderByBuyerAndIdempotencyKeyParams{
		BuyerID:        conv.UUID(buyerID),
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			return nil, domain_order.ErrOrderNotFound
		}
		return nil, fmt.Errorf("infra: get order by idempotency key: %w", err)
	}

	items, err := r.getQuerier(ctx).ListOrderItemsByOrderID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("infra: list order items by idempotency key: %w", err)
	}

	histories, err := r.getQuerier(ctx).ListOrderHistoryByOrderID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("infra: list order history by idempotency key: %w", err)
	}

	return toDomainOrder(row, items, histories), nil
}

func (r *orderRepository) SaveOrderState(ctx context.Context, order *domain_order.Order, expectedStatus domain_order.OrderStatus) error {
	rowsAffected, err := r.getQuerier(ctx).UpdateOrderStatus(ctx, toUpdateOrderStatusParams(order, expectedStatus))
	if err != nil {
		return fmt.Errorf("infra: update order state: %w", err)
	}
	if rowsAffected == 0 {
		return domain_order.ErrOrderInvalidStateTransition
	}

	if len(order.OrderHistory) > 0 {
		last := *order
		last.OrderHistory = order.OrderHistory[len(order.OrderHistory)-1:]
		if err := r.getQuerier(ctx).CreateOrderHistoryBatch(ctx, toCreateOrderHistoryParams(&last)); err != nil {
			return fmt.Errorf("infra: create order state history: %w", err)
		}
	}

	return nil
}

func isUniqueConstraint(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" && pgErr.ConstraintName == constraintName
	}
	return false
}
