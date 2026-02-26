package startup

import (
	"log"

	"github.com/HugManh/cronjob/configs"
	"github.com/HugManh/cronjob/internal/routing"
	"github.com/HugManh/cronjob/pkg/db/postgres"
	"github.com/HugManh/cronjob/pkg/taskmanager"
)

type Shutdown = func()

var migrate = false

func Server() {
	config := configs.LoadConfig()
	router, shutdown := create(config)
	defer shutdown()
	router.Start(config.ServerHost, config.ServerPort)

}

func create(cfg *configs.Config) (routing.Router, Shutdown) {

	// Init Database
	config := postgres.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
		DBName:   cfg.DBName,
		SSLMode:  "disable", // or "require" based on your needs
	}

	db, err := postgres.NewDatabase(config)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if migrate {
		if err := db.Migrate(); err != nil {
			log.Fatalf("Database migration failed: %v", err)
		}
	}

	// Init TaskManager
	tm := taskmanager.NewTaskManager()
	tm.Cron.Start()

	// Register tasks from DB
	if err := tm.LoadTasksFromDB(db.DB()); err != nil {
		log.Fatalf("Registering tasks from DB failed: %v", err)
	}

	router := routing.NewRouter("debug")
	router.LoadControllers(db.DB(), tm)

	shutdown := func() {
		db.Disconnect()
		tm.Cron.Stop()
	}

	return router, shutdown
}
