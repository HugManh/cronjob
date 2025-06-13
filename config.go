package main

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string `env:"APP_ENV" envDefault:"development"`
	AppPort string `env:"APP_PORT" envDefault:"8080"`

	DBHost string `env:"DB_HOST" envDefault:"localhost"`
	DBPort int    `env:"DB_PORT" envDefault:"5432"`
	DBUser string `env:"DB_USER" envDefault:"postgres"`
	DBPass string `env:"DB_PASSWORD" envDefault:""`
	DBName string `env:"DB_DATABASE" envDefault:"postgres"`
	DBSsl  bool   `env:"DB_SSL" envDefault:"true"`
}

var (
	cfg  *Config
	once sync.Once
)

// LoadConfig loads environment variables into Config struct
func LoadConfig() *Config {
	once.Do(func() {
		// Optional: load .env file if present (dev environment)
		_ = godotenv.Load()

		cfg = &Config{}
		if err := env.Parse(cfg); err != nil {
			log.Fatalf("Failed to parse env: %v", err)
		}
	})
	return cfg
}
