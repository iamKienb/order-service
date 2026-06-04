package inventory

import (
	"context"
	"errors"

	"inventory-command-module/db/repository"
	domain_inventory "inventory-command-module/internal/domain/inventory"

	pgx "github.com/iamKienb/go-core/postgres"
	"github.com/jackc/pgx/v5/pgconn"
)

type inventoryRepository struct {
	queries *repository.Queries
}

func NewRepository(service pgx.PGXService) domain_inventory.Repository {
	return &inventoryRepository{
		queries: repository.New(service.GetPool()),
	}
}

func (r *inventoryRepository) getQuerier(ctx context.Context) *repository.Queries {
	if tx := pgx.ExtractTx(ctx); tx != nil {
		return r.queries.WithTx(tx)
	}

	return r.queries
}

func (r *inventoryRepository) isDuplicateSKU(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}
