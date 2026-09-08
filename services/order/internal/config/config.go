package config

import (
	"os"
	"strconv"
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

func MustLoad() *Config {
	validate := validator.New(validator.WithRequiredStructEnabled())

	cfg := &Config{
		Env: parseString(os.Getenv("ENV"), EnvLocal),
		HTTPConfig: HTTPConfig{
			Address:         parseString(os.Getenv("HTTP_ADDR"), ":8080"),
			Timeout:         parseDuration(os.Getenv("HTTP_TIMEOUT"), defaultTimeout),
			IdleTimeout:     parseDuration(os.Getenv("HTTP_IDLE_TIMEOUT"), defaultIdleTimeout),
			ShutdownTimeout: parseDuration(os.Getenv("HTTP_SHUTDOWN_TIMEOUT"), defaultShutdownTimeout),
		},
	}

	if err := validate.Struct(cfg); err != nil {
		panic("failed to validate config: " + err.Error())
	}

	return cfg
}

func parseInt(val string, defVal int) int {
	if val == "" {
		return defVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		panic("failed to parse int from string")
	}
	return i
}

func parseString(val string, defVal string) string {
	if val == "" {
		return defVal
	}
	return val
}

func parseDuration(val string, def time.Duration) time.Duration {
	if val == "" {
		return def
	}
	dur, err := time.ParseDuration(val)
	if err != nil {
		panic("failed to parse duration from string: " + err.Error())
	}
	return dur
}
