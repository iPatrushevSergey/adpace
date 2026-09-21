package port

import "context"

type Transactor interface {
	RunInTransaction(ctx context.Context, retryer Retryer, fn func(ctx context.Context) error) error
}
