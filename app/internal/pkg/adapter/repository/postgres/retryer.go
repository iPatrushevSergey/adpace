package postgres

import (
	"context"
	"errors"
	"io"
	"math"
	"math/rand/v2"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// retryOptions holds the configuration for the retry operation.
type retryOptions struct {
	maxRetries int
	backoff    port.BackoffFunc
}

// RetryOption configures a Retryer.
type RetryOption func(*retryOptions)

// WithMaxRetries sets the maximum number of retry attempts.
func WithMaxRetries(n int) RetryOption {
	return func(rc *retryOptions) { rc.maxRetries = n }
}

// WithExponentialBackoff sets exponential backoff with full jitter.
func WithExponentialBackoff(base, max time.Duration) RetryOption {
	return func(rc *retryOptions) {
		rc.backoff = func(attempt int) time.Duration {
			delay := time.Duration(float64(base) * math.Pow(2, float64(attempt)))
			if delay > max {
				delay = max
			}
			if delay > 0 {
				delay = time.Duration(rand.Int64N(int64(delay)))
			}
			return delay
		}
	}
}

// WithConstantBackoff sets a fixed delay between retries.
func WithConstantBackoff(d time.Duration) RetryOption {
	return func(rc *retryOptions) { rc.backoff = func(_ int) time.Duration { return d } }
}

// WithBackoffFunc sets a custom backoff strategy.
func WithBackoffFunc(fn port.BackoffFunc) RetryOption {
	return func(rc *retryOptions) { rc.backoff = fn }
}

type Retryer struct {
	base retryOptions
}

func NewRetryer(opts ...RetryOption) *Retryer {
	cfg := retryOptions{
		maxRetries: 3,
		backoff:    func(_ int) time.Duration { return 100 * time.Millisecond },
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Retryer{base: cfg}
}

func (r *Retryer) Do(ctx context.Context, op func() error) error {
	var err error
	for attempt := 0; attempt <= r.base.maxRetries; attempt++ {
		err = op()
		if err == nil {
			return nil
		}
		if !IsRetriableError(err) {
			return err
		}
		if attempt == r.base.maxRetries {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(r.base.backoff(attempt)):
		}
	}
	return err
}

func (r *Retryer) With(tuning port.RetryTuning) port.Retryer {
	cfg := r.base
	if tuning.MaxRetries > 0 {
		cfg.maxRetries = tuning.MaxRetries
	}
	if tuning.Backoff != nil {
		cfg.backoff = tuning.Backoff
	}
	return &Retryer{base: cfg}
}

// IsRetriableError checks whether a PostgreSQL or network error is transient and the operation can be retried.
// Reference: https://www.postgresql.org/docs/current/errcodes-appendix.html
func IsRetriableError(err error) bool {
	if err == nil {
		return false
	}

	// Network-level errors (connection refused, timeout, broken pipe, EOF)
	if isNetworkError(err) {
		return true
	}
	return isPostgresError(err)
}

// isNetworkError checks for common transient network errors.
func isNetworkError(err error) bool {
	// EOF — connection was closed by the server
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	// Connection refused, reset, broken pipe
	if errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EPIPE) {
		return true
	}
	// OS-level timeout
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}
	// net.Error with Timeout (e.g. dial timeout, read timeout)
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	return false
}

// isPostgresError checks for transient PostgreSQL-level errors by error code.
func isPostgresError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	code := pgErr.Code

	// Class 08 — Connection Exception (connection lost, broken pipe)
	// Class 25 — Invalid Transaction State (transaction was aborted externally)
	// Class 40 — Transaction Rollback (deadlock, serialization failure)
	// Class 53 — Insufficient Resources (out of memory, too many connections)
	if pgerrcode.IsConnectionException(code) ||
		pgerrcode.IsInvalidTransactionState(code) ||
		pgerrcode.IsTransactionRollback(code) ||
		pgerrcode.IsInsufficientResources(code) {
		return true
	}

	// Class 57 — Operator Intervention (server shutting down, crash recovery)
	switch code {
	case pgerrcode.AdminShutdown,
		pgerrcode.CrashShutdown,
		pgerrcode.CannotConnectNow,
		pgerrcode.DatabaseDropped,
		pgerrcode.IdleInTransactionSessionTimeout:
		return true
	}

	return false
}
