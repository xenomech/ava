package logger

import (
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

type logger struct {
	log *zap.Logger
}

var (
	instance *logger
	one      sync.Once

	settings = struct {
		env   string
		level string
	}{env: "local", level: "info"}
)

func Init(env, level string) {
	settings.env = env
	settings.level = level
}

func String(key, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
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

func getInstance() *logger {
	one.Do(func() {
		instance = &logger{log: build().WithOptions(zap.AddCallerSkip(1))}
	})

	return instance
}

func build() *zap.Logger {
	level := parseLevel(settings.level)

	if settings.env == "production" {
		config := zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(level)

		built, err := config.Build()
		if err != nil {
			return zap.NewNop()
		}

		return built
	}

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
		level,
	)

	return zap.New(core, zap.AddCaller())
}

func parseLevel(name string) zapcore.Level {
	switch strings.ToLower(name) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
