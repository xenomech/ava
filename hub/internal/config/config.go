package config

import (
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	DeviceID      string
	APIBaseURL    string
	MQTTBrokerURL string
	LogLevel      slog.Level
}

func GetConfig() *Config {
	once.Do(func() {
		instance = load()
	})

	return instance
}

func load() *Config {
	return &Config{
		DeviceID:      env("DEVICE_ID", ""),
		APIBaseURL:    env("API_BASE_URL", "http://localhost:8000/api/v1"),
		MQTTBrokerURL: env("MQTT_BROKER_URL", "tcp://localhost:1883"),
		LogLevel:      level(env("LOG_LEVEL", "info")),
	}
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func level(name string) slog.Level {
	switch strings.ToLower(name) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
