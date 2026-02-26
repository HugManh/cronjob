package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogMode string

const (
	Production  LogMode = "production"
	Development LogMode = "development"
)

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
	Fatal LogLevel = "fatal"
)

type LogEncoding string

const (
	Console LogEncoding = "console"
)

// type RotateLogger struct {
// 	MaxSize    int  `mapstructure:"max_size" json:"max_size"` // megabytes
// 	MaxBackups int  `mapstructure:"max_backups" json:"max_backups"`
// 	MaxAge     int  `mapstructure:"max_age" json:"max_age"`   // days
// 	Compress   bool `mapstructure:"compress" json:"compress"` // disabled by default
// }

// Config định cấu hình logger với các tùy chọn log level, encoding, đường dẫn file, và rotation
type Config struct {
	Level    LogLevel    `mapstructure:"level" json:"level"`       // Mức log: debug, info, warn, error, fatal
	Encoding LogEncoding `mapstructure:"encoding" json:"encoding"` // Định dạng: console hoặc json
	LogDir   string      `mapstructure:"log_dir" json:"log_dir"`   // Đường dẫn thư mục chứa log
	LogFile  string      `mapstructure:"log_file" json:"log_file"` // Tên file log
	ZapType  string      `mapstructure:"zap_type" json:"zap_type"` // Loại zap: sugar hoặc không

	// Cấu hình rotation cho log files
	MaxSize    int  `mapstructure:"max_size" json:"max_size"` // megabytes
	MaxBackups int  `mapstructure:"max_backups" json:"max_backups"`
	MaxAge     int  `mapstructure:"max_age" json:"max_age"`   // days
	Compress   bool `mapstructure:"compress" json:"compress"` // Nén file log cũ (disabled by default)
}

var loggerLevelMap = map[LogLevel]zapcore.Level{
	Debug: zapcore.DebugLevel,
	Info:  zapcore.InfoLevel,
}

type ILogger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Logger struct {
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
	zapSugar    bool
}

var Logging *Logger = &Logger{}

// setDefaults thiết lập giá trị mặc định cho các trường config nếu chưa được khởi tạo
func setDefaults(cfg *Config) {
	if cfg.Level == "" {
		cfg.Level = Info
	}
	if cfg.Encoding == "" {
		cfg.Encoding = Console
	}
	if cfg.LogDir == "" {
		cfg.LogDir = "./logs"
	}
	if cfg.LogFile == "" {
		cfg.LogFile = "app.log"
	}
	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100 // 100 MB
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 5
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30 // 30 days
	}
	if cfg.ZapType == "" {
		cfg.ZapType = "sugar"
	}
}

// getPath tạo đường dẫn đầy đủ cho file log và tạo thư mục nếu chưa tồn tại
func getPath(logDir, logFile string) string {
	fullDir := filepath.Join(logDir)
	if _, err := os.Stat(fullDir); os.IsNotExist(err) {
		if err := os.MkdirAll(fullDir, 0755); err != nil {
			panic(err)
		}
	}

	return filepath.Join(fullDir, logFile)
}

// configure khởi tạo WriteSyncer với lumberjack rotation logger, ghi log cả stdout và file
func configure(cfg *Config) zapcore.WriteSyncer {
	path := getPath(cfg.LogDir, cfg.LogFile)
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    cfg.MaxSize, // megabytes
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,   // days
		Compress:   cfg.Compress, // disabled by default
	})
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(w),
	)
}

func GetLogger() ILogger {
	return Logging
}

// NewLogger khởi tạo logger từ config, thiết lập các core cho các log level khác nhau
func NewLogger(cfg *Config) ILogger {
	setDefaults(cfg)

	var coreArr []zapcore.Core
	logLevel, exist := loggerLevelMap[cfg.Level]
	if !exist {
		logLevel = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.LevelKey = "level"
	encoderCfg.CallerKey = "caller"
	encoderCfg.TimeKey = "time"
	encoderCfg.NameKey = "name"
	encoderCfg.MessageKey = "msg"
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeDuration = zapcore.MillisDurationEncoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	var encoder zapcore.Encoder
	if cfg.Encoding == Console {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	// Thiết lập priority filter: Error level ghi vào stderr, Info/Debug ghi vào file
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= logLevel
	})

	infoFileCore := zapcore.NewCore(encoder, configure(cfg), lowPriority)
	// infoFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)), lowPriority)
	errorFileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr)), highPriority)

	coreArr = append(coreArr, infoFileCore)
	coreArr = append(coreArr, errorFileCore)

	// Kết hợp các cores để xử lý log từ cả hai priority level
	core := zapcore.NewTee(coreArr...)
	loggerzap := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugarLogger := loggerzap.Sugar()

	Logging = &Logger{
		sugarLogger: sugarLogger,
		logger:      loggerzap,
		zapSugar:    strings.Contains(cfg.ZapType, "sugar"),
	}
	return Logging
}

func (l *Logger) Debug(args ...interface{}) {
	if l.zapSugar {
		l.sugarLogger.Debug(args...)
		return
	}
	l.logger.Debug(fmt.Sprint(args...))
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	if l.zapSugar {
		l.sugarLogger.Debugf(template, args...)
		return
	}
	l.logger.Debug(fmt.Sprintf(template, args...))
}

func (l *Logger) Info(args ...interface{}) {
	if l.zapSugar {
		l.sugarLogger.Info(args...)
		return
	}
	l.logger.Info(fmt.Sprint(args...))
}

func (l *Logger) Infof(template string, args ...interface{}) {
	if l.zapSugar {
		l.sugarLogger.Infof(template, args...)
		return
	}
	l.logger.Info(fmt.Sprintf(template, args...))
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *Logger) Error(args ...interface{}) {
	if l.zapSugar {
		l.sugarLogger.Error(args...)
		return
	}
	l.logger.Error(fmt.Sprint(args...))
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	if l.zapSugar {
		l.sugarLogger.Errorf(template, args...)
		return
	}
	l.logger.Error(fmt.Sprintf(template, args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
