package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	DB() *gorm.DB
	Close() error
}

type Config struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type database struct {
	client *gorm.DB
}

func Open(config Config) (Database, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
		config.SSLMode,
	)

	client, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &database{
		client: client,
	}, nil
}

func (db *database) DB() *gorm.DB {
	return db.client
}

func (db *database) Close() error {
	sqlDB, err := db.client.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Close(); err != nil {
		return err
	}

	return nil
}
