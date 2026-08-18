package config

import (
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	DeviceName        string
	APIBaseURL        string
	MQTTBrokerURL     string
	StateFile         string
	HeartbeatInterval time.Duration
	LogLevel          slog.Level
}

func GetConfig() *Config {
	once.Do(func() {
		instance = load()
	})

	return instance
}

func load() *Config {
	return &Config{
		DeviceName:        env("DEVICE_NAME", defaultDeviceName()),
		APIBaseURL:        env("API_BASE_URL", "http://localhost:8000/api/v1"),
		MQTTBrokerURL:     env("MQTT_BROKER_URL", "tcp://localhost:1883"),
		StateFile:         env("STATE_FILE", "avahub-state.json"),
		HeartbeatInterval: duration(env("HEARTBEAT_INTERVAL", "60s"), time.Minute),
		LogLevel:          level(env("LOG_LEVEL", "info")),
	}
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func defaultDeviceName() string {
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "Ava Hub"
	}

	return host
}

func duration(value string, fallback time.Duration) time.Duration {
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}

	return parsed
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
