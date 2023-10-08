package logger2

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
)

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel

	_minLevel = DebugLevel
	_maxLevel = FatalLevel

	InvalidLevel = _maxLevel + 1
)

type (
	Level = zapcore.Level

	Config struct {
		LogDir      string
		ServiceName string
		MaxFileSize int
		MaxAge      int
		MaxBackups  int
		LocalTime   bool
		Compress    bool
		Filename    string
	}

	Logger2 struct {
		zapLogger   *zap.Logger
		LogDir      string
		ServiceName string
		Filename    string
	}
)

var (
	mutex   sync.RWMutex
	_global *Logger2
)

func Logger() *Logger2 {
	mutex.RLock()
	defer mutex.RUnlock()

	l := _global
	return l
}

func NewLoggerDefault(dir string) (*Logger2, error) {
	cfg := &Config{
		LogDir:      filepath.Clean(dir),
		MaxFileSize: 10,
		MaxAge:      28,
		MaxBackups:  7,
		LocalTime:   true,
		Compress:    true,
	}

	return NewLogger(cfg)
}

func NewLogger(cfg *Config) (*Logger2, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, err := os.Stat(cfg.LogDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %s, %v", cfg.LogDir, err)
		}
	}

	if cfg.ServiceName == "" {
		path, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %v", err)
		}

		executable := filepath.Base(path)
		ext := filepath.Ext(executable)

		if ext != "" {
			executable = executable[:len(executable)-len(ext)]
		}

		cfg.ServiceName = executable
	}

	cfg.Filename = filepath.Join(cfg.LogDir, fmt.Sprintf("%s.log", cfg.ServiceName))

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxFileSize,
			MaxAge:     cfg.MaxAge,
			MaxBackups: cfg.MaxBackups,
			LocalTime:  cfg.LocalTime,
			Compress:   cfg.Compress,
		}),
		zapcore.InfoLevel,
	)

	_global = &Logger2{
		zapLogger:   zap.New(core),
		LogDir:      cfg.LogDir,
		ServiceName: cfg.ServiceName,
		Filename:    cfg.Filename,
	}

	return _global, nil
}

func (l *Logger2) Log(level Level, format string, fields ...any) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(1)).Check(level, fmt.Sprintf(format, fields...))
}

func (l *Logger2) Error(message string, fields ...zap.Field) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(1)).Error(message, fields...)
}

func (l *Logger2) Info(message string, fields ...zap.Field) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(1)).Info(message, fields...)
}

func (l *Logger2) Debug(message string, fields ...zap.Field) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(1)).Debug(message, fields...)
}

func (l *Logger2) Warn(message string, fields ...zap.Field) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(1)).Warn(message, fields...)
}

func (l *Logger2) Fatal(message string, fields ...zap.Field) {
	l.zapLogger.WithOptions(zap.AddCallerSkip(1)).Fatal(message, fields...)
}

//func Error(format string, fields ...any) {
//	_global.zapLogger.Error(fmt.Sprintf(format, fields...))
//}
//
//func Info(format string, fields ...any) {
//	_global.zapLogger.Info(fmt.Sprintf(format, fields...))
//}
//
//func ErrorAndStdout(format string, fields ...any) {
//	msg := fmt.Sprintf(format, fields...)
//	_global.zapLogger.Error(msg)
//	fmt.Println(msg)
//}
//
//func InfoAndStdout(format string, fields ...any) {
//	msg := fmt.Sprintf(format, fields...)
//	_global.zapLogger.Info(msg)
//	fmt.Println(msg)
//}
//
//func Fatal(format string, fields ...any) {
//	_global.zapLogger.Fatal(fmt.Sprintf(format, fields...))
//}
//
//func Debug(format string, fields ...any) {
//	_global.zapLogger.Debug(fmt.Sprintf(format, fields...))
//}
//
//func Warn(format string, fields ...any) {
//	_global.zapLogger.Warn(fmt.Sprintf(format, fields...))
//}
//
//func ErrorJson(obj any) {
//	data, _ := json.Marshal(obj)
//	_global.zapLogger.Error(string(data))
//}
//
//func InfoJson(obj any) {
//	data, _ := json.Marshal(obj)
//	_global.zapLogger.Info(string(data))
//}
//
//func DebugJson(obj any) {
//	data, _ := json.Marshal(obj)
//	_global.zapLogger.Debug(string(data))
//}
//
//func FatalJson(obj any) {
//	data, _ := json.Marshal(obj)
//	_global.zapLogger.Fatal(string(data))
//}
//
//func ErrorAndStdoutJson(obj any) {
//	data, _ := json.Marshal(obj)
//	msg := string(data)
//	_global.zapLogger.Error(msg)
//	fmt.Println(msg)
//}
//
//func InfoAndStdoutJson(obj any) {
//	data, _ := json.Marshal(obj)
//	msg := string(data)
//	_global.zapLogger.Info(msg)
//	fmt.Println(msg)
//}
//
//func WarnJson(obj any) {
//	data, _ := json.Marshal(obj)
//	_global.zapLogger.Warn(string(data))
//}
//
//func WarnAndStdout(format string, fields ...any) {
//	msg := fmt.Sprintf(format, fields...)
//	_global.zapLogger.Warn(msg)
//	fmt.Println(msg)
//}
