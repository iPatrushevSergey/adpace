package postgres

import (
	"fmt"
	"time"
)

type RetryConfig struct {
	Attempts  int           `mapstructure:"max_attempts"`
	BaseDelay time.Duration `mapstructure:"base_delay"`
	MaxDelay  time.Duration `mapstructure:"max_delay"`
}

func (c *RetryConfig) Validate() error {
	if c.Attempts < 0 {
		return fmt.Errorf("max_attempts must be non-negative")
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
