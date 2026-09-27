package retry

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	retrygo "github.com/avast/retry-go/v4"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// WithExponentialBackoff sets exponential backoff with full jitter.
func ExponentialBackoff(base, max time.Duration) port.BackoffFunc {
	return func(attempt int) time.Duration {
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

// ConstantBackoff returns a port.BackoffFunc with a fixed delay.
func ConstantBackoff(d time.Duration) port.BackoffFunc {
	return func(_ int) time.Duration { return d }
}

type Retryer struct {
	defaultConfig port.RetryConfig
}

func NewRetryer(opts ...port.RetryOption) *Retryer {
	cfg := port.RetryConfig{
		Attempts: 4, // 1 initial try + 3 retries
		Backoff:  ConstantBackoff(100 * time.Millisecond),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return &Retryer{defaultConfig: cfg}
}

func (r *Retryer) Do(ctx context.Context, op func() error, opts ...port.RetryOption) error {
	cfg := r.defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return retrygo.Do(
		op,
		retrygo.Context(ctx),
		retrygo.Attempts(uint(cfg.Attempts)),
		retrygo.DelayType(func(n uint, _ error, _ *retrygo.Config) time.Duration {
			return cfg.Backoff(int(n))
		}),
		retrygo.RetryIf(IsRetriableError),
	)
}

// IsRetriableError checks whether an error is transient and the operation can be retried.
func IsRetriableError(err error) bool {
	if err == nil {
		return false
	}
	if isNetworkError(err) {
		return true
	}
	if isPostgresError(err) {
		return true
	}
	return isHTTPRetriable(err)
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

// isHTTPRetriable checks for transient HTTP-level errors by status code.
func isHTTPRetriable(err error) bool {
	var statusErr *HTTPStatusError
	if !errors.As(err, &statusErr) {
		return false
	}

	switch statusErr.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}

	return false
}

// HTTPStatusError wraps a non-2xx HTTP response so the retryer can classify it.
type HTTPStatusError struct {
	StatusCode int
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("unexpected http status: %d", e.StatusCode)
}
