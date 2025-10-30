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
	DB

	APIKey string `env:"EXCHANGERATE_API_KEY"`
	Port   string `env:"APP_PORT"`
}

type DB struct {
	User     string `env:"DB_USER" `
	Password string `env:"DB_PASSWORD"  `
	Name     string `env:"DB_NAME" `
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT" `
	SSLMode  string `env:"DB_SSLMODE" `
}
