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
	Host         string `env:"DATABASE_HOST" env-default:"localhost"`
	Port         int64  `env:"DATABASE_PORT" env-default:"5432"`
	User         string `env:"DATABASE_USER"`
	Password     string `env:"DATABASE_PASSWORD"`
	DatabaseName string `env:"DATABASE_NAME"`
}

func Read() (Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
