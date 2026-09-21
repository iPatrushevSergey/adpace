package port

import (
	"context"
	"time"
)

// BackoffFunc calculates the delay for the given attempt.
type BackoffFunc func(attempt int) time.Duration

type RetryTuning struct {
	MaxRetries int
	Backoff    BackoffFunc
}

type Retryer interface {
	Do(ctx context.Context, op func() error) error
	With(tuning RetryTuning) Retryer
}
