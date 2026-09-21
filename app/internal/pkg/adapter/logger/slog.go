package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// SlogLogger implements application logging using log/slog.
type SlogLogger struct {
	sl *slog.Logger
}

// NewSlogLogger builds a slog logger from config.
func NewSlogLogger(cfg Config) (*SlogLogger, error) {
	lvl, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	var handler slog.Handler
	switch cfg.Format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	default:
		return nil, fmt.Errorf("format: unknown value %q", cfg.Format)
	}

	return &SlogLogger{sl: slog.New(handler)}, nil
}

// Debug implements the Debug method of the Logger port.
func (s *SlogLogger) Debug(ctx context.Context, msg string, args ...any) {
	s.sl.Log(ctx, slog.LevelDebug, msg, args...)
}

// Info implements the Info method of the Logger port.
func (s *SlogLogger) Info(ctx context.Context, msg string, args ...any) {
	s.sl.Log(ctx, slog.LevelInfo, msg, args...)
}

// Warn implements the Warn method of the Logger port.
func (s *SlogLogger) Warn(ctx context.Context, msg string, args ...any) {
	s.sl.Log(ctx, slog.LevelWarn, msg, args...)
}

// Error implements the Error method of the Logger port.
func (s *SlogLogger) Error(ctx context.Context, msg string, args ...any) {
	s.sl.Log(ctx, slog.LevelError, msg, args...)
}

// parseLevel parses a slog level from a string.
func parseLevel(level string) (slog.Leveler, error) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return nil, fmt.Errorf("parse log level: unknown level %q", level)
	}
}
