package config

import (
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/mikestefanello/hooks"
	"github.com/mikestefanello/hooks-example/pkg/app"
	"github.com/samber/do"
)

type (
	// Config provides system configuration
	Config interface {
		// GetHTTP returns HTTP configuration
		GetHTTP() HTTPConfig

		// GetApp returns App configuration
		GetApp() AppConfig
	}

	// Base stores complete configuration
	Base struct {
		HTTP HTTPConfig
		App  AppConfig
	}

	// HTTPConfig stores HTTP configuration
	HTTPConfig struct {
		Hostname     string        `env:"HTTP_HOSTNAME"`
		Port         uint16        `env:"HTTP_PORT,default=8000"`
		ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT,default=5s"`
		WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT,default=10s"`
		IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT,default=2m"`
	}

	// AppConfig stores application configuration
	AppConfig struct {
		Name    string        `env:"APP_NAME,default=Hooks"`
		Timeout time.Duration `env:"APP_TIMEOUT,default=20s"`
	}
)

func init() {
	// Provide dependencies during app boot process
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewConfig)
	})
}

// NewConfig creates a new Config instance
func NewConfig(i *do.Injector) (Config, error) {
	var cfg Base
	err := envdecode.StrictDecode(&cfg)
	return &cfg, err
}

// GetHTTP returns HTTP configuration
func (c *Base) GetHTTP() HTTPConfig {
	return c.HTTP
}

// GetApp returns app configuration
func (c *Base) GetApp() AppConfig {
	return c.App
}
