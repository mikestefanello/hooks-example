package config

import (
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/mikestefanello/hooks"
	"github.com/mikestefanello/hooks-example/pkg/app"
	"github.com/samber/do"
)

type (
	// Config stores complete configuration
	Config struct {
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
	app.HookBoot.Listen(func(e hooks.Event[*do.Injector]) {
		do.Provide(e.Msg, NewConfig)
	})
}

func NewConfig(i *do.Injector) (*Config, error) {
	var cfg Config
	err := envdecode.StrictDecode(&cfg)
	return &cfg, err
}
