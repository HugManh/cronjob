package main

import (
	"github.com/HugManh/cronjob/configs"
	"github.com/HugManh/cronjob/internal/startup"
	"github.com/HugManh/cronjob/pkg/logger"
)

func init() {
	cfg := configs.Load()
	if err := logger.Configure(logger.Config{
		ServiceName: cfg.ServiceName,
		Level:       cfg.LogLevel,
		Directory:   cfg.LogDir,
	}); err != nil {
		panic(err)
	}
}

func main() {
	startup.Server()
}
