package config

import (
	"errors"
	"fmt"
	"io/fs"
	"sync"

	"github.com/spf13/viper"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	Port      string
	ServerEnv string
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
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("PORT", "8000")
	v.SetDefault("SERVER_ENV", "local")
}

func isMissingConfigFile(err error) bool {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return errors.Is(pathErr, fs.ErrNotExist)
	}

	var notFoundErr viper.ConfigFileNotFoundError

	return errors.As(err, &notFoundErr)
}
