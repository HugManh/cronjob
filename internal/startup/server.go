package startup

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/HugManh/cronjob/configs"
	"github.com/HugManh/cronjob/internal/routing"
	"github.com/HugManh/cronjob/internal/service"
	"github.com/HugManh/cronjob/pkg/database/postgres"
	"github.com/HugManh/cronjob/pkg/logger"
)

type Shutdown = func()

type databaseClient interface {
	DB() *gorm.DB
	Close() error
}

func Server() {
	cfg := configs.Load()
	router, shutdown := create(cfg)
	defer shutdown()
	router.Start(cfg.ServerHost, cfg.ServerPort)

}

func create(cfg *configs.Config) (routing.Router, Shutdown) {

	// Init Database
	db, err := createDatabase(cfg)
	if err != nil {
		logger.Fatalf("database connection failed: %v", err)
	}

	// Init TaskManager
	tm := service.NewTaskManager()
	tm.Cron.Start()

	// Register tasks from DB
	logger.Info("loading tasks from database...")
	if err := tm.Startup(db.DB()); err != nil {
		logger.Fatalf("registering tasks from DB failed: %v", err)
	}

	router := routing.NewRouter(ginMode(cfg.Env))
	router.LoadControllers(db.DB(), tm)

	shutdown := func() {
		db.Close()
		tm.Cron.Stop()
	}

	return router, shutdown
}

func createDatabase(cfg *configs.Config) (databaseClient, error) {
	switch cfg.DBEngine {
	case "postgres", "supabase":
		return postgres.Open(postgres.Config{
			Host:     cfg.DBHost,
			Port:     cfg.DBPort,
			User:     cfg.DBUser,
			Password: cfg.DBPass,
			DBName:   cfg.DBName,
			SSLMode:  cfg.DBSSLMode,
		})
	case "mysql":
		return nil, fmt.Errorf("database engine %q is configured but no MySQL driver is implemented", cfg.DBEngine)
	default:
		return nil, fmt.Errorf("unsupported database engine %q", cfg.DBEngine)
	}
}

func ginMode(env string) string {
	if env == "production" {
		return "release"
	}
	return "debug"
}
