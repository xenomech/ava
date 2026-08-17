package logger

import (
	"os"
	"sync"
	"time"

	"ava/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	log *zap.Logger
}

type Field = zap.Field

var (
	instance *logger
	one      sync.Once
)

func String(key, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Any(key string, value any) Field {
	return zap.Any(key, value)
}

func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Err(err error) Field {
	return zap.Error(err)
}

func Duration(key string, d time.Duration) Field {
	return zap.Duration(key, d)
}

func getInstance() *logger {
	one.Do(func() {
		env := config.GetConfig().ServerEnv

		var zapLogger *zap.Logger
		var err error

		if env == "production" {
			zapLogger, err = zap.NewProduction()
		} else {
			zapLogger, err = newDevelopmentLogger()
		}

		if err != nil {
			panic("Failed to init logger: " + err.Error())
		}

		instance = &logger{log: zapLogger.WithOptions(zap.AddCallerSkip(1))}
	})

	return instance
}

func newDevelopmentLogger() (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	return zap.New(core, zap.AddCaller()), nil
}

func Info(msg string, fields ...Field) {
	getInstance().log.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	getInstance().log.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	getInstance().log.Error(msg, fields...)
}

func Debug(msg string, fields ...Field) {
	getInstance().log.Debug(msg, fields...)
}

func Sync() {
	_ = getInstance().log.Sync()
}
