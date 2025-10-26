package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

func Parse() (Config, error) {
	
	godotenv.Load()

	var cfg Config
	err := env.Parse(&cfg)

	return cfg, err
}

type Config struct {
	DB DB

	APIKey string `env:"EXCHANGERATE_API_KEY" envDefault:""`
	Port   string `env:"APP_PORT" envDefault:"4000"`
}

type DB struct {
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD, unset" `
	Name     string `env:"DB_NAME" envDefault:"balance"`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}
