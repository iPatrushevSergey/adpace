package postgres

import (
	"fmt"
	"strings"
	"time"
)

type PoolConfig struct {
	URI         string        `mapstructure:"uri"`
	MaxConns    int32         `mapstructure:"max_conns"`
	MinConns    int32         `mapstructure:"min_conns"`
	MaxConnLife time.Duration `mapstructure:"max_conn_life"`
	MaxConnIdle time.Duration `mapstructure:"max_conn_idle"`
	HealthCheck time.Duration `mapstructure:"health_check"`
}

func (c *PoolConfig) Validate() error {
	c.URI = strings.TrimSpace(c.URI)
	if c.URI == "" {
		return fmt.Errorf("uri is required")
	}
	if c.MaxConns <= 0 {
		return fmt.Errorf("max_conns must be positive")
	}
	if c.MinConns < 0 {
		return fmt.Errorf("min_conns must be non-negative")
	}
	if c.MinConns > c.MaxConns {
		return fmt.Errorf("min_conns must not exceed max_conns")
	}
	if c.MaxConnLife <= 0 {
		return fmt.Errorf("max_conn_life must be positive")
	}
	if c.MaxConnIdle <= 0 {
		return fmt.Errorf("max_conn_idle must be positive")
	}
	if c.HealthCheck <= 0 {
		return fmt.Errorf("health_check must be positive")
	}
	return nil
}
