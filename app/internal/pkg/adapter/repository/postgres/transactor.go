package postgres

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/jackc/pgx/v5/pgxpool"
)

var nestedSettings = settings.Must(settings.WithPropagation(trm.PropagationNested))

type Transactor struct {
	manager *trmmanager.Manager
}

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{manager: trmmanager.Must(trmpgx.NewDefaultFactory(pool))}
}

func (t *Transactor) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.manager.Do(ctx, fn)
}

func (t *Transactor) DoInNestedTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.manager.DoWithSettings(ctx, nestedSettings, fn)
}
