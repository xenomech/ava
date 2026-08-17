package db

import (
	"errors"
	"time"

	"ava/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func Connect(pgconfig *PostgresConfig) (*gorm.DB, error) {
	if pgconfig.Host == "" || pgconfig.User == "" ||
		pgconfig.Password == "" || pgconfig.Database == "" ||
		pgconfig.Port == "" {
		logger.Error("DB_CONNECTION_ERROR", logger.String("message", "Postgres configuration is incomplete"))

		return nil, errors.New("postgres configuration is incomplete")
	}

	if pgconfig.SSLMode == "" {
		pgconfig.SSLMode = "disable"
	}

	dsn := "host=" + pgconfig.Host +
		" user=" + pgconfig.User +
		" password=" + pgconfig.Password +
		" dbname=" + pgconfig.Database +
		" port=" + pgconfig.Port +
		" sslmode=" + pgconfig.SSLMode

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		logger.Error("DB_CONNECTION_ERROR", logger.Err(err))

		return nil, err
	}

	sqlDB, err := database.DB()
	if err != nil {
		logger.Error("DB_POOL_CONFIG_ERROR", logger.Err(err))

		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	logger.Info("DB_CONNECTED")

	return database, nil
}

func Disconnect(database *gorm.DB) error {
	if database == nil {
		return nil
	}

	sqlDB, err := database.DB()
	if err != nil {
		logger.Error("DB_DISCONNECT_ERROR", logger.Err(err))

		return err
	}

	if err := sqlDB.Close(); err != nil {
		logger.Error("DB_DISCONNECT_ERROR", logger.Err(err))

		return err
	}

	logger.Info("DB_DISCONNECTED")

	return nil
}
