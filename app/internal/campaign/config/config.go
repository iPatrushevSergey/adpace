// Package config loads campaign-manager service settings.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	mapstructure "github.com/go-viper/mapstructure/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/logger"
	"github.com/iPatrushevSergey/adpace/app/internal/pkg/adapter/repository/postgres"
)

// Config holds grouped campaign-manager configuration.
type Config struct {
	Server  Server               `mapstructure:"server"`
	Logger  logger.Config        `mapstructure:"logger"`
	DBPool  postgres.PoolConfig  `mapstructure:"db_pool"`
	DBRetry postgres.RetryConfig `mapstructure:"db_retry"`
}

// LoadConfig loads campaign-manager config.
// Field priority: flags > env > yaml > viper defaults.
// Config file path: flag -c > env CONFIG > default path.
func LoadConfig() (Config, error) {
	fs, err := newFlagSet()
	if err != nil {
		return Config{}, err
	}

	dotenvPath, dotenvLoaded, err := loadDotEnv()
	if err != nil {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	if dotenvLoaded {
		_, _ = fmt.Fprintf(os.Stderr, "config: loaded dotenv %s\n", dotenvPath)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "config: dotenv file not found, continuing with env/yaml/defaults")
	}

	configPath := resolveConfigPath(fs)

	v := viper.New()
	setDefaults(v)

	if err := readConfigFile(v, configPath); err != nil {
		return Config{}, err
	}

	bindEnv(v)

	if err := bindFlags(v, fs); err != nil {
		return Config{}, fmt.Errorf("bind flags: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg, viper.DecodeHook(durationDecodeHook())); err != nil {
		return Config{}, fmt.Errorf("unmarshal: %w", err)
	}

	if err := finalizeConfig(&cfg, configPath); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func newFlagSet() (*pflag.FlagSet, error) {
	fs := pflag.NewFlagSet("campaign", pflag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringP("address", "a", "", "listen address (host:port)")
	fs.StringP("config", "c", "", "path to config file")
	fs.StringP("log-level", "l", "", "logging level")
	fs.String("shutdown-timeout", "", "graceful shutdown timeout")
	fs.StringP("database-uri", "d", "", "PostgreSQL DSN")
	fs.String("tls-cert", "", "path to TLS certificate PEM")
	fs.String("tls-key", "", "path to TLS private key PEM")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("flag parsing: %w", err)
	}
	return fs, nil
}

func loadDotEnv() (string, bool, error) {
	for _, file := range []string{"app/.env", ".env"} {
		if _, err := os.Stat(file); err == nil {
			if err := godotenv.Load(file); err != nil {
				return "", false, fmt.Errorf("load %s: %w", file, err)
			}
			return file, true, nil
		}
	}
	return "", false, nil
}

func resolveConfigPath(fs *pflag.FlagSet) string {
	path := "app/configs/campaign.yaml"
	source := "default"
	if env, ok := os.LookupEnv("CONFIG"); ok {
		if trimmedPath := strings.TrimSpace(env); trimmedPath != "" {
			path = trimmedPath
			source = "env"
		}
	}
	if fs.Changed("config") {
		flagPath, err := fs.GetString("config")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "config: read -c flag: %v\n", err)
		} else {
			path = flagPath
			source = "flag"
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "config: config file path source=%s path=%s\n", source, path)
	return path
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.address", "127.0.0.1:8080")
	v.SetDefault("server.shutdown_timeout", "10s")
	v.SetDefault("server.cert_file", "")
	v.SetDefault("server.key_file", "")
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("db_pool.uri", "")
	v.SetDefault("db_pool.max_conns", 25)
	v.SetDefault("db_pool.min_conns", 5)
	v.SetDefault("db_pool.max_conn_life", "1h")
	v.SetDefault("db_pool.max_conn_idle", "30m")
	v.SetDefault("db_pool.health_check", "1m")
	v.SetDefault("db_retry.max_retries", 3)
	v.SetDefault("db_retry.base_delay", "100ms")
	v.SetDefault("db_retry.max_delay", "2s")
}

func readConfigFile(v *viper.Viper, configPath string) error {
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			_, _ = fmt.Fprintf(os.Stderr, "config: file not found path=%s\n", v.ConfigFileUsed())
			return nil
		}
		return fmt.Errorf("read config: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stderr, "config: loaded %s\n", v.ConfigFileUsed())
	return nil
}

func bindEnv(v *viper.Viper) {
	v.SetEnvPrefix("CAMPAIGN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}

func bindFlags(v *viper.Viper, fs *pflag.FlagSet) error {
	bindings := map[string]string{
		"server.address":          "address",
		"server.shutdown_timeout": "shutdown-timeout",
		"logger.level":            "log-level",
		"db_pool.uri":             "database-uri",
		"server.cert_file":        "tls-cert",
		"server.key_file":         "tls-key",
	}
	for key, flagName := range bindings {
		f := fs.Lookup(flagName)
		if f == nil {
			return fmt.Errorf("flag not found: %s", flagName)
		}
		if err := v.BindPFlag(key, f); err != nil {
			return fmt.Errorf("bind %s: %w", flagName, err)
		}
	}
	return nil
}

// finalizeConfig validates loaded config and applies cross-cutting rules.
func finalizeConfig(cfg *Config, configPath string) error {
	if err := cfg.Logger.Validate(); err != nil {
		return fmt.Errorf("logger: %w", err)
	}
	if err := cfg.DBPool.Validate(); err != nil {
		return fmt.Errorf("db_pool: %w", err)
	}
	if err := cfg.DBRetry.Validate(); err != nil {
		return fmt.Errorf("db_retry: %w", err)
	}
	if err := cfg.Server.Validate(); err != nil {
		return fmt.Errorf("server: %w", err)
	}

	if !strings.HasSuffix(strings.ToLower(filepath.Base(configPath)), ".prod.yaml") {
		return nil
	}
	if !cfg.Server.TLSEnabled() {
		return fmt.Errorf("production config requires TLS (server.cert_file and server.key_file)")
	}
	return nil
}

// durationDecodeHook maps incoming scalars to time.Duration.
func durationDecodeHook() mapstructure.DecodeHookFunc {
	durType := reflect.TypeOf(time.Duration(0))
	return func(from reflect.Type, to reflect.Type, data any) (any, error) {
		if to != durType {
			return data, nil
		}

		var d Duration
		switch v := data.(type) {
		case string:
			if err := d.Set(v); err != nil {
				return nil, err
			}
			return d.Duration, nil
		case int:
			return time.Duration(v) * time.Second, nil
		case int64:
			return time.Duration(v) * time.Second, nil
		case float64:
			return time.Duration(v) * time.Second, nil
		default:
			return nil, fmt.Errorf("unsupported duration type %T", data)
		}
	}
}
