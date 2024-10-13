package cfg

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database
	Server
	Kafka
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

type Kafka struct {
	Host               string `env:"KAFKA_HOST" env-default:"localhost"`
	Port               int64  `env:"KAFKA_PORT" env-default:"9092"`
	StatusChangedTopic string `env:"KAFKA_TOPIC_STATUS_CHANGED" env-default:"devices.status"`
}

func Read() (Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	return cfg, nil
}
