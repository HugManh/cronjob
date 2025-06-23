package startup

import (
	"context"
	"log"

	"github.com/HugManh/cronjob/internal/network"
	"github.com/HugManh/cronjob/pkg/config"
	"github.com/HugManh/cronjob/pkg/db"
	"github.com/HugManh/cronjob/pkg/taskmanager"
)

type Shutdown = func()

func Server() {
	config := config.LoadConfig()
	router, shutdown := create(config)
	defer shutdown()
	router.Start(config.ServerHost, config.ServerPort)

}

func create(cfg *config.Config) (network.Router, Shutdown) {
	context := context.Background()

	// Init Database
	dbConfig := db.DbConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
		DBName:   cfg.DBName,
		SSLMode:  "disable", // or "require" based on your needs
	}

	db := db.NewDatabase(context, dbConfig)
	db.Connect()

	// Init TaskManager
	tm := taskmanager.NewTaskManager()
	tm.Cron.Start()

	// Register tasks from DB
	if err := tm.LoadTasksFromDB(db.GetInstance().Database); err != nil {
		log.Fatalf("Registering tasks from DB failed: %v", err)
	}

	router := network.NewRouter("debug")
	router.LoadControllers(db.GetInstance().Database, tm)

	shutdown := func() {
		db.Disconnect()
		tm.Cron.Stop()
	}

	return router, shutdown
}
