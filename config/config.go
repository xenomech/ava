package config

import (
	"errors"
	"fmt"
	"io/fs"
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

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBDatabase string
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

	v.SetConfigFile(".env")
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

		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetString("DB_PORT"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBDatabase: v.GetString("DB_DATABASE"),
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("PORT", "8000")
	v.SetDefault("SERVER_ENV", "local")

	v.SetDefault("JWT_ACCESS_EXPIRY", 15*time.Minute)
	v.SetDefault("JWT_REFRESH_EXPIRY", 7*24*time.Hour)

	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
}

func isMissingConfigFile(err error) bool {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return errors.Is(pathErr, fs.ErrNotExist)
	}

	var notFoundErr viper.ConfigFileNotFoundError

	return errors.As(err, &notFoundErr)
}
