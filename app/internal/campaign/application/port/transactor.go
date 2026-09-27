package port

import "context"

type Transactor interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
	DoInNestedTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
