package postgres

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Querier is what an operation receives: the active transaction or the pool.
type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

// Executor is the repositories' only entry point to the database.
// It never begins or ends transactions — that is Transactor's job.
type Executor struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewExecutor(pool *pgxpool.Pool) *Executor {
	return &Executor{pool: pool, getter: trmpgx.DefaultCtxGetter}
}

// Do runs op against the active transaction if ctx carries one, else the pool.
//
// Retries: outside a transaction a single statement is retried by retryer.
// Inside one, op runs exactly once: a failed statement aborts the whole
// transaction (SQLSTATE 25P02), so retrying it can only fail again. The
// transaction's owner (Transactor.RunInTransaction) retries it as a whole.
func (e *Executor) Do(ctx context.Context, retryer port.Retryer, op func(q Querier) error) error {
	if tx := e.getter.DefaultTrOrDB(ctx, nil); tx != nil {
		return op(tx)
	}
	return retryer.Do(ctx, func() error { return op(e.pool) })
}
