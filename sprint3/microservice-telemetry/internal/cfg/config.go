package cfg

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database
	Server
}

type Server struct {
	Port int64 `env:"SERVER_PORT" env-default:"8080"`
}

type Database struct {
	URL    string `env:"DATABASE_URL"`
	Token  string `env:"DATABASE_TOKEN"`
	Org    string `env:"DATABASE_ORG"`
	Bucket string `env:"DATABASE_BUCKET"`
}

func Read() (Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
