package outbox

import (
	"context"

	"order-command-module/db/repository"
	"order-command-module/internal/application/port"

	pgx "github.com/iamKienb/go-core/postgres"
)

type repositoryImpl struct {
	queries *repository.Queries
}

func NewRepository(service pgx.PGXService) port.OutboxRepository {
	return &repositoryImpl{
		queries: repository.New(service.GetPool()),
	}
}

func (r *repositoryImpl) getQuerier(ctx context.Context) *repository.Queries {
	if tx := pgx.ExtractTx(ctx); tx != nil {
		return r.queries.WithTx(tx)
	}

	return r.queries
}
