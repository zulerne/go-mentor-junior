package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	EnvLocal = "local"
	EnvProd  = "production"
)

const (
	defaultTimeout         = 5 * time.Second
	defaultIdleTimeout     = 60 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

// TODO (review): Is it okay to use validatorv10 for config validation? or is it better to use other libs or manually?

type Config struct {
	Env        string     `validate:"required,oneof=local production"`
	HTTPConfig HTTPConfig `validate:"omitempty"`
}

type HTTPConfig struct {
	Address         string        `validate:"omitempty"`
	Timeout         time.Duration `validate:"omitempty"`
	IdleTimeout     time.Duration `validate:"omitempty"`
	ShutdownTimeout time.Duration `validate:"omitempty"`
}

func Load() (*Config, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	var errs []error

	timeout, err := parseDuration(os.Getenv("HTTP_TIMEOUT"), defaultTimeout)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to parse HTTP_TIMEOUT: %w", err))
	}
	idleTimeout, err := parseDuration(os.Getenv("HTTP_IDLE_TIMEOUT"), defaultIdleTimeout)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to parse HTTP_IDLE_TIMEOUT: %w", err))
	}
	shutdownTimeout, err := parseDuration(os.Getenv("HTTP_SHUTDOWN_TIMEOUT"), defaultShutdownTimeout)
	if err != nil {
		errs = append(errs, fmt.Errorf("failed to parse HTTP_SHUTDOWN_TIMEOUT: %w", err))
	}

	cfg := &Config{
		Env: parseString(os.Getenv("ENV"), EnvLocal),
		HTTPConfig: HTTPConfig{
			Address:         parseString(os.Getenv("HTTP_ADDR"), ":8080"),
			Timeout:         timeout,
			IdleTimeout:     idleTimeout,
			ShutdownTimeout: shutdownTimeout,
		},
	}

	if err = validate.Struct(cfg); err != nil {
		errs = append(errs, fmt.Errorf("failed to validate config: %w", err))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return cfg, nil
}

func parseString(val string, defVal string) string {
	if val == "" {
		return defVal
	}
	return val
}

func parseDuration(val string, def time.Duration) (time.Duration, error) {
	if val == "" {
		return def, nil
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration from string: %w", err)
	}
	return dur, nil
}
