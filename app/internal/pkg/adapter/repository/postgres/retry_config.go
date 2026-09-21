package postgres

import (
	"fmt"
	"time"
)

type RetryConfig struct {
	MaxRetries int           `mapstructure:"max_retries"`
	BaseDelay  time.Duration `mapstructure:"base_delay"`
	MaxDelay   time.Duration `mapstructure:"max_delay"`
}

func (c *RetryConfig) Validate() error {
	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries must be non-negative")
	}
	if c.BaseDelay <= 0 {
		return fmt.Errorf("base_delay must be positive")
	}
	if c.MaxDelay <= 0 {
		return fmt.Errorf("max_delay must be positive")
	}
	if c.MaxDelay < c.BaseDelay {
		return fmt.Errorf("max_delay must not be less than base_delay")
	}
	return nil
}
