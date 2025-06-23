package config

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	// server
	ServerHost string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	ServerPort uint16 `env:"SERVER_PORT" envDefault:"8080"`
	Env        string `env:"ENVIRONMENT" envDefault:"local"`
	// database
	DBHost string `env:"DB_HOST" envDefault:"localhost"`
	DBPort uint16 `env:"DB_PORT" envDefault:"5432"`
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
