package port

import (
	"context"
	"time"
)

// BackoffFunc calculates the delay before the next attempt, 0-based.
type BackoffFunc func(attempt int) time.Duration

type RetryConfig struct {
	// Attempts is the total number of tries, including the first one.
	Attempts int
	Backoff  BackoffFunc
}

type RetryOption func(*RetryConfig)

func WithAttempts(n int) RetryOption {
	return func(rc *RetryConfig) { rc.Attempts = n }
}

func WithBackoffFunc(fn BackoffFunc) RetryOption {
	return func(rc *RetryConfig) { rc.Backoff = fn }
}

type Retryer interface {
	Do(ctx context.Context, op func() error, opts ...RetryOption) error
}
