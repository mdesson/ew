package Config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	EchoPort     int    `env:"ECHO_PORT,default=8000"`
	FrontendPort int    `env:"FRONTEND_PORT,default=8080"`
	LogLevel     string `env:"LOG_LEVEL,default=info"`
}

func New() (Config, error) {
	ctx := context.Background()

	var c Config
	if err := envconfig.Process(ctx, &c); err != nil {
		return c, err
	}

	return c, nil
}
