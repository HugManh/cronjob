package postgres

import (
	"fmt"

	"github.com/HugManh/cronjob/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IDatabase interface {
	DB() *gorm.DB
	Migrate() error
	Disconnect() error
}

type Config struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DBName   string
	SSLMode  string // "disable", "require", etc.
}

type database struct {
	Client *gorm.DB
}

func NewDatabase(config Config) (IDatabase, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	client, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &database{
		Client: client,
	}, nil
}

func (db *database) DB() *gorm.DB {
	return db.Client
}

func (db *database) Migrate() error {
	return db.Client.AutoMigrate(&model.Task{}, &model.Slack{})
}

func (db *database) Disconnect() error {
	sqlDB, err := db.Client.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	return nil
}
