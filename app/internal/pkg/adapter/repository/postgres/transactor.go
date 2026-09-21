package postgres

import (
	"context"

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

type txKey struct{}

type Transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{pool: pool}
}

// RunInTransaction executes fn inside a transaction. If ctx already carries
// an active transaction, fn reuses it — no nested Begin, no savepoint.
func (t *Transactor) RunInTransaction(
	ctx context.Context,
	retryer port.Retryer,
	fn func(ctx context.Context) error,
) error {
	// For nested transactions.
	if _, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return fn(ctx)
	}

	// Starting a transaction with repeats.
	return retryer.Do(ctx, func() error {
		tx, err := t.pool.Begin(ctx)
		if err != nil {
			return err
		}
		defer tx.Rollback(ctx)

		if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
			return err
		}
		return tx.Commit(ctx)
	})
}

// Returns an existing transaction or connection pool.
func (t *Transactor) GetExecutor(ctx context.Context) Executor {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return t.pool
}
