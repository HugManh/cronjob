package configs

import (
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `mapstructure:"SERVICE_NAME"`

	ServerHost string `mapstructure:"SERVER_HOST"`
	ServerPort uint16 `mapstructure:"SERVER_PORT"`
	Env        string `mapstructure:"ENVIRONMENT"`

	LogLevel string `mapstructure:"LOG_LEVEL"`
	LogDir   string `mapstructure:"LOG_DIR"`

	DBEngine  string `mapstructure:"DB_ENGINE"`
	DBHost    string `mapstructure:"DB_HOST"`
	DBPort    uint16 `mapstructure:"DB_PORT"`
	DBUser    string `mapstructure:"DB_USER"`
	DBPass    string `mapstructure:"DB_PASSWORD"`
	DBName    string `mapstructure:"DB_NAME"`
	DBSSLMode string `mapstructure:"DB_SSLMODE"`
}

var (
	cfg  *Config
	once sync.Once
)

func Load() *Config {
	once.Do(func() {
		v := viper.New()
		v.SetConfigName(".env")
		v.SetConfigType("env")
		v.AddConfigPath(".")
		v.AutomaticEnv()

		setDefaults(v)

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				log.Fatalf("failed to read config file: %v", err)
			}
		}

		config := &Config{}
		if err := v.Unmarshal(config); err != nil {
			log.Fatalf("failed to unmarshal config: %v", err)
		}

		config.normalize()
		cfg = config
	})

	return cfg
}

func LoadConfig() *Config {
	return Load()
}

func (c *Config) normalize() {
	c.ServiceName = strings.TrimSpace(c.ServiceName)
	c.ServerHost = strings.TrimSpace(c.ServerHost)
	c.Env = strings.ToLower(strings.TrimSpace(c.Env))
	c.LogLevel = strings.ToLower(strings.TrimSpace(c.LogLevel))
	c.LogDir = strings.TrimSpace(c.LogDir)
	c.DBEngine = strings.ToLower(strings.TrimSpace(c.DBEngine))
	c.DBHost = strings.TrimSpace(c.DBHost)
	c.DBUser = strings.TrimSpace(c.DBUser)
	c.DBPass = strings.TrimSpace(c.DBPass)
	c.DBName = strings.TrimSpace(c.DBName)
	c.DBSSLMode = strings.ToLower(strings.TrimSpace(c.DBSSLMode))
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("SERVICE_NAME", "cronjob")
	v.SetDefault("SERVER_HOST", "0.0.0.0")
	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("ENVIRONMENT", "local")
	v.SetDefault("LOG_LEVEL", "debug")
	v.SetDefault("LOG_DIR", "log")
	v.SetDefault("DB_ENGINE", "postgres")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "")
	v.SetDefault("DB_NAME", "postgres")
	v.SetDefault("DB_SSLMODE", "disable")
}
