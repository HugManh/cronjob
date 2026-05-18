package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLevel       = log.DebugLevel
	defaultLogDir      = "log"
	defaultServiceName = "app"
)

type Config struct {
	ServiceName string
	Level       string
	Directory   string
}

func Configure(config Config) error {
	log.SetLevel(parseLevel(config.Level))
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if err := configureOutput(config.ServiceName, config.Directory); err != nil {
		return err
	}

	return nil
}

func parseLevel(value string) log.Level {
	level, err := log.ParseLevel(value)
	if err != nil {
		return defaultLevel
	}
	return level
}

func configureOutput(serviceName, logDir string) error {
	if logDir == "" {
		logDir = defaultLogDir
	}
	if serviceName == "" {
		serviceName = defaultServiceName
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create log directory: %w", err)
	}

	logFilePath := filepath.Join(logDir, fmt.Sprintf("%s.log", serviceName))
	logFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Infof("logging to file: %s", logFilePath)
	return nil
}
