package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
	Trace LogLevel = "trace"
	Panic LogLevel = "panic"
	Fatal LogLevel = "fatal"
)

var loggerLevelMap = map[LogLevel]log.Level{
	Debug: log.DebugLevel,
	Info:  log.InfoLevel,
	Error: log.ErrorLevel,
	Warn:  log.WarnLevel,
	Trace: log.TraceLevel,
	Panic: log.PanicLevel,
	Fatal: log.FatalLevel,
}

func NewLogger() {
	log.SetLevel(loggerLevelMap[LogLevel(os.Getenv("LOG_LEVEL"))])
	log.SetReportCaller(true)
	log.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "2006-01-02 15:04:05",
		ShowFullLevel:   true,
		CallerFirst:     true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			_, filename := filepath.Split(f.File)
			return fmt.Sprintf(" (%s:%d)", filename, f.Line)
		},
	})

	// Setup log file
	if err := setupLogFile(os.Getenv("SERVICE_NAME")); err != nil {
		log.Errorf("Failed to setup log file: %v", err)
	}
}

func setupLogFile(service_name string) error {
	// Create logs directory if not exists
	logsDir := "log"
	if os.Getenv("LOG_DIR") != "" {
		logsDir = os.Getenv("LOG_DIR")
	}

	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Setup log file name
	if service_name == "" {
		service_name = "app"
	}
	logFileName := fmt.Sprintf("%s.log", service_name)
	logFilePath := filepath.Join(logsDir, logFileName)

	// Setup lumberjack for log rotation
	logFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100,  // megabytes - rotate when file reaches this size
		MaxBackups: 30,   // keep 30 old log files
		MaxAge:     30,   // days - delete old log files after 30 days
		Compress:   true, // compress rotated files
		LocalTime:  true, // use local time for filenames
	}

	// Write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Infof("Logging to file: %s", logFilePath)
	return nil
}
