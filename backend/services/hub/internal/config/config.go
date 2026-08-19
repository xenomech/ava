package config

import (
	"os"
	"sync"
	"time"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	HubName           string
	APIBaseURL        string
	MQTTBrokerURL     string
	StateFile         string
	HeartbeatInterval time.Duration
	SyncInterval      time.Duration
	DiscoveryTimeout  time.Duration
	LogLevel          string
	Env               string
}

func GetConfig() *Config {
	once.Do(func() {
		instance = load()
	})

	return instance
}

func load() *Config {
	return &Config{
		HubName:           env("HUB_NAME", defaultHubName()),
		APIBaseURL:        env("API_BASE_URL", "http://localhost:8000/api/v1"),
		MQTTBrokerURL:     env("MQTT_BROKER_URL", "tcp://localhost:1883"),
		StateFile:         env("STATE_FILE", "avahub-state.json"),
		HeartbeatInterval: duration(env("HEARTBEAT_INTERVAL", "60s"), time.Minute),
		SyncInterval:      duration(env("SYNC_INTERVAL", "60s"), time.Minute),
		DiscoveryTimeout:  duration(env("DISCOVERY_TIMEOUT", "4s"), 4*time.Second),
		LogLevel:          env("LOG_LEVEL", "info"),
		Env:               env("HUB_ENV", "local"),
	}
}

func env(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}

	return fallback
}

func defaultHubName() string {
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
