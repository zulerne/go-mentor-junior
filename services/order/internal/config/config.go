package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

const (
	EnvLocal = "local"
	EnvProd  = "production"

	envPrefix = ""

	defaultTimeout         = 5 * time.Second
	defaultIdleTimeout     = 60 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

type Config struct {
	Env  string     `koanf:"env"  validate:"required,oneof=local production"`
	HTTP HTTPConfig `koanf:"http" validate:"omitempty"`
}

type HTTPConfig struct {
	Address         string        `koanf:"address"`
	Timeout         time.Duration `koanf:"timeout"`
	IdleTimeout     time.Duration `koanf:"idletimeout"`
	ShutdownTimeout time.Duration `koanf:"shutdowntimeout"`
}

func Load(v *validator.Validate) (*Config, error) {
	k := koanf.New(".")

	defaults := map[string]any{
		"env":                  EnvLocal,
		"http.address":         ":8080",
		"http.timeout":         defaultTimeout,
		"http.idletimeout":     defaultIdleTimeout,
		"http.shutdowntimeout": defaultShutdownTimeout,
	}

	if err := k.Load(confmap.Provider(defaults, "."), nil); err != nil {
		return nil, fmt.Errorf("load defaults: %w", err)
	}

	// Env vars override defaults.
	// ENV=production            → env
	// HTTP_ADDRESS=:9090        → http.address
	// HTTP_IDLETIMEOUT=30s      → http.idletimeout
	// HTTP_SHUTDOWNTIMEOUT=15s  → http.shutdowntimeout
	if err := k.Load(env.Provider(envPrefix, ".", func(s string) string {
		s = strings.TrimPrefix(s, envPrefix)
		s = strings.ToLower(s)
		return strings.ReplaceAll(s, "_", ".")
	}), nil); err != nil {
		return nil, fmt.Errorf("load env: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := v.Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}
