package port

import "context"

type UseCase[In, Out any] interface {
	Execute(ctx context.Context, in In) (Out, error)
}
