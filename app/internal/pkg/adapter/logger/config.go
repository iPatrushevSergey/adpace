package logger

import (
	"fmt"
	"strings"
)

// Config defines logger initialization parameters.
type Config struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Validate trims string fields and checks logger settings.
func (c *Config) Validate() error {
	c.Level = strings.ToLower(strings.TrimSpace(c.Level))
	switch c.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("level: unknown value %q", c.Level)
	}

	c.Format = strings.ToLower(strings.TrimSpace(c.Format))
	switch c.Format {
	case "json", "text":
	default:
		return fmt.Errorf("format: unknown value %q", c.Format)
	}

	return nil
}
