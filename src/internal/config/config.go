package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Host string `env:"SERVER_HOST" envDefault:"localhost"`
		Port string `env:"SERVER_PORT" envDefault:"8080"`
	}
	PostgreSQL struct {
		Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
		Port     string `env:"POSTGRES_PORT" envDefault:"5432"`
		User     string `env:"POSTGRES_USER" envDefault:"postgres"`
		Password string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
		DBName   string `env:"POSTGRES_DB" envDefault:"survey"`
	}
	MongoDB struct {
		URI string `env:"MONGODB_URI" envDefault:"mongodb://localhost:27017"`
	}
	JWT JWTConfig
}

type JWTConfig struct {
	Secret          string        `env:"JWT_SECRET" envDefault:"your-secret-key"`
	AccessDuration  time.Duration `env:"ACCESS_TOKEN_DURATION" envDefault:"15m"`
	RefreshDuration time.Duration `env:"REFRESH_TOKEN_DURATION" envDefault:"24h"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
