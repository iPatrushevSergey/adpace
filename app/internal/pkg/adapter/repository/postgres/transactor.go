package postgres

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transactor struct {
	trManager *trmmanager.Manager
}

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{
		trManager: trmmanager.Must(trmpgx.NewDefaultFactory(pool)),
	}
}

// RunInTransaction executes fn inside a transaction. If ctx already carries
// an active transaction, it's reused (default PropagationRequired) — no
// nested Begin, no savepoint. Use RunInNestedTransaction for that.
func (t *Transactor) RunInTransaction(
	ctx context.Context,
	retryer port.Retryer,
	fn func(ctx context.Context) error,
) error {
	return retryer.Do(ctx, func() error {
		return t.trManager.Do(ctx, fn)
	})
}

// RunInNestedTransaction executes fn in a real nested transaction (SAVEPOINT)
// when ctx already carries an active transaction — the inner part can roll
// back independently of the outer one.
func (t *Transactor) RunInNestedTransaction(
	ctx context.Context,
	retryer port.Retryer,
	fn func(ctx context.Context) error,
) error {
	s := settings.Must(settings.WithPropagation(trm.PropagationNested))
	return retryer.Do(ctx, func() error {
		return t.trManager.DoWithSettings(ctx, s, fn)
	})
}
