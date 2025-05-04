package project

import (
	"context"
	"errors"
)

// Config defined all necessary parameters for the project wizard to execute.
type Config struct {
	ProjectName string
	ProjectDir  string
	ModName     string
	MuxoVersion string
}

// ErrConfigNotFound is passed to panic if config cannot be found inside the context.
var ErrConfigNotFound = errors.New("project wizard config was not found")

type key int

var configWizardKey key

// ConfigFromContext returns the User value stored in ctx, if any.
func ConfigFromContext(ctx context.Context) (*Config, bool) {
	cfg, ok := ctx.Value(configWizardKey).(*Config)
	return cfg, ok
}

// NewConfigContext returns a new Context that carries value cfg.
func NewConfigContext(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, configWizardKey, cfg)
}
