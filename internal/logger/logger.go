package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init инициализирует глобальный logger
func Init(environment string, debug bool) error {
	var config zap.Config

	// Определяем конфигурацию на основе окружения
	if environment == "production" {
		config = zap.NewProductionConfig()
		// JSON формат для production
		config.Encoding = "json"
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	} else {
		config = zap.NewDevelopmentConfig()
		// Console формат для development
		config.Encoding = "console"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		
		if debug {
			config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		} else {
			config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		}
	}

	// Логирование в файл и консоль
	config.OutputPaths = []string{"stdout"}
	if environment == "production" {
		config.OutputPaths = append(config.OutputPaths, "server.log")
	}
	config.ErrorOutputPaths = []string{"stderr"}

	// Создаем logger
	var err error
	Log, err = config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	return nil
}

// Sync синхронизирует буферы logger (вызывать при shutdown)
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// Wrapper функции для удобного использования

func Debug(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Fatal(msg, fields...)
	} else {
		// Fallback если logger не инициализирован
		os.Exit(1)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if Log != nil {
		Log.Panic(msg, fields...)
	}
}

