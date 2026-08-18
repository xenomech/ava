package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	Port      string
	ServerEnv string

	JwtSecretKey     string
	JwtAccessExpiry  time.Duration
	JwtRefreshExpiry time.Duration

	DeviceCodeExpiry time.Duration
	HubPollInterval  time.Duration
	HubTokenExpiry   time.Duration

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBDatabase string

	CORSAllowedOrigins string
	CORSAllowedMethods string
	CORSAllowedHeaders string
	CORSMaxAge         int

	ResendAPIKey    string
	ResendFromEmail string
	ResendFromName  string
	AppURL          string
}

func GetConfig() *Config {
	once.Do(func() {
		instance = load()
	})

	return instance
}

func load() *Config {
	v := viper.New()

	v.AutomaticEnv()

	setDefaults(v)

	v.SetConfigFile(envFile())
	v.SetConfigType("env")

	if err := v.ReadInConfig(); err != nil && !isMissingConfigFile(err) {
		panic(fmt.Sprintf("CONFIG_READ_ERROR: %v", err))
	}

	return &Config{
		Port:      v.GetString("PORT"),
		ServerEnv: v.GetString("SERVER_ENV"),

		JwtSecretKey:     v.GetString("JWT_SECRET"),
		JwtAccessExpiry:  v.GetDuration("JWT_ACCESS_EXPIRY"),
		JwtRefreshExpiry: v.GetDuration("JWT_REFRESH_EXPIRY"),

		DeviceCodeExpiry: v.GetDuration("HUB_CODE_EXPIRY"),
		HubPollInterval:  v.GetDuration("HUB_POLL_INTERVAL"),
		HubTokenExpiry:   v.GetDuration("HUB_TOKEN_EXPIRY"),

		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetString("DB_PORT"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBDatabase: v.GetString("DB_DATABASE"),

		CORSAllowedOrigins: v.GetString("CORS_ALLOWED_ORIGINS"),
		CORSAllowedMethods: v.GetString("CORS_ALLOWED_METHODS"),
		CORSAllowedHeaders: v.GetString("CORS_ALLOWED_HEADERS"),
		CORSMaxAge:         v.GetInt("CORS_MAX_AGE"),

		ResendAPIKey:    v.GetString("RESEND_API_KEY"),
		ResendFromEmail: v.GetString("RESEND_FROM_EMAIL"),
		ResendFromName:  v.GetString("RESEND_FROM_NAME"),
		AppURL:          v.GetString("APP_URL"),
	}
}

func envFile() string {
	for _, candidate := range []string{".env", filepath.Join("..", ".env")} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return ".env"
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("PORT", "8000")
	v.SetDefault("SERVER_ENV", "local")

	v.SetDefault("JWT_ACCESS_EXPIRY", 15*time.Minute)
	v.SetDefault("JWT_REFRESH_EXPIRY", 7*24*time.Hour)

	v.SetDefault("HUB_CODE_EXPIRY", 10*time.Minute)
	v.SetDefault("HUB_POLL_INTERVAL", 5*time.Second)
	v.SetDefault("HUB_TOKEN_EXPIRY", time.Hour)

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")

	v.SetDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	v.SetDefault("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS")
	v.SetDefault("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,X-Requested-With,X-User-ID,X-Tenant-ID,X-REQUEST-ID,X-TIME,X-Device-Name")
	v.SetDefault("CORS_MAX_AGE", 86400)

	v.SetDefault("RESEND_FROM_EMAIL", "noreply@example.com")
	v.SetDefault("RESEND_FROM_NAME", "Ava")
	v.SetDefault("APP_URL", "http://localhost:3000")
}

func isMissingConfigFile(err error) bool {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return errors.Is(pathErr, fs.ErrNotExist)
	}

	var notFoundErr viper.ConfigFileNotFoundError

	return errors.As(err, &notFoundErr)
}
