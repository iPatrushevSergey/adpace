package postgres

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	trmmanager "github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// It is needed to return a transaction or a separate connection in the repository.
type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

type Transactor struct {
	trManager *trmmanager.Manager
	getter    *trmpgx.CtxGetter
	pool      *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{
		trManager: trmmanager.Must(trmpgx.NewDefaultFactory(pool)),
		getter:    trmpgx.DefaultCtxGetter,
		pool:      pool,
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

// Returns an existing transaction or connection pool.
func (t *Transactor) GetExecutor(ctx context.Context) Executor {
	return t.getter.DefaultTrOrDB(ctx, t.pool)
}
