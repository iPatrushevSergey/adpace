package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Server holds HTTP server settings.
type Server struct {
	Address         string        `mapstructure:"address"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	CertFile        string        `mapstructure:"cert_file"`
	KeyFile         string        `mapstructure:"key_file"`
}

// TLSEnabled reports whether both certificate and key paths are configured.
func (s Server) TLSEnabled() bool {
	return s.CertFile != "" && s.KeyFile != ""
}

// Validate trims server fields, normalizes address and verifies TLS cert/key files when configured.
func (s *Server) Validate() error {
	s.Address = strings.TrimSpace(s.Address)
	s.CertFile = strings.TrimSpace(s.CertFile)
	s.KeyFile = strings.TrimSpace(s.KeyFile)

	var addr Address
	if err := addr.Set(s.Address); err != nil {
		return fmt.Errorf("address: %w", err)
	}
	s.Address = addr.String()

	if s.ShutdownTimeout <= 0 {
		return fmt.Errorf("shutdown_timeout must be positive")
	}

	if !s.TLSEnabled() {
		return nil
	}
	info, err := os.Stat(s.CertFile)
	if err != nil {
		return fmt.Errorf("tls certificate: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("tls certificate: %s is a directory", s.CertFile)
	}
	info, err = os.Stat(s.KeyFile)
	if err != nil {
		return fmt.Errorf("tls private key: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("tls private key: %s is a directory", s.KeyFile)
	}
	return nil
}
