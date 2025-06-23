package db

import (
	"context"
	"fmt"

	mtask "github.com/HugManh/cronjob/internal/tasks/model"
	mslack "github.com/HugManh/cronjob/internal/slack/model"
	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetInstance() *database
	Connect()
	Disconnect()
}

type DbConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DBName   string
	SSLMode  string // "disable", "require", etc.
}

type database struct {
	Database *gorm.DB
	context  context.Context
	config   DbConfig
}

func NewDatabase(ctx context.Context, config DbConfig) Database {
	return &database{
		context: ctx,
		config:  config,
	}
}

func (db *database) GetInstance() *database {
	return db
}

func (db *database) Connect() {
	cfg := db.config
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	log.Printf("Connecting to database with DSN: %s", dsn)

	log.Println("connecting postgres...", dsn)
	client, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal("connection to postgres failed!: ", err)
	}

	client.AutoMigrate(&mtask.Task{}, &mslack.Slack{})
	db.Database = client
	log.Println("connected to Postgres!")
}

func (db *database) Disconnect() {
	fmt.Println("disconnecting postgres...")
	sqlDB, err := db.Database.DB()
	if err != nil {
		log.Panic(err)
	}

	if err := sqlDB.Close(); err != nil {
		log.Panic(err)
	}
	fmt.Println("disconnected postgres!")
}
