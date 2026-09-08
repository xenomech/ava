package db

import (
	"errors"
	"strings"
	"time"

	"ava/api/internal/model"
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

func Migrate(database *gorm.DB) error {
	if database == nil {
		return nil
	}

	if err := database.AutoMigrate(
		&model.User{},
		&model.Tenant{},
		&model.TenantMembership{},
		&model.Session{},
		&model.APIToken{},
		&model.Token{},
		&model.Room{},
		&model.Hub{},
		&model.HubAuthorization{},
		&model.Device{},
		&model.Flow{},
		&model.FlowStep{},
		&model.Scene{},
		&model.SceneTarget{},
	); err != nil {
		logger.Error("DB_MIGRATION_ERROR", logger.Err(err))

		return err
	}

	// Uniqueness has to mean "among the rows that still exist", so soft deleted rows must not hold a name.
	for _, index := range []partialIndex{
		{table: "rooms", name: "idx_room_tenant_name", columns: "tenant_id, name"},
		{table: "tenant_memberships", name: "idx_membership_tenant_user", columns: "tenant_id, user_id"},
	} {
		if err := onlyLiveRows(database, index); err != nil {
			logger.Error("DB_MIGRATION_ERROR", logger.Err(err))

			return err
		}
	}

	if database.Migrator().HasColumn(&model.Device{}, "room") {
		if err := database.Migrator().DropColumn(&model.Device{}, "room"); err != nil {
			logger.Error("DB_MIGRATION_ERROR", logger.Err(err))

			return err
		}
	}

	logger.Info("DB_MIGRATED")

	return nil
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

type partialIndex struct {
	table   string
	name    string
	columns string
}

// onlyLiveRows rewrites a unique index to ignore soft deleted rows, checking first to avoid a needless lock.
func onlyLiveRows(database *gorm.DB, index partialIndex) error {
	var definition string

	err := database.Raw(
		"SELECT indexdef FROM pg_indexes WHERE tablename = ? AND indexname = ?",
		index.table, index.name,
	).Scan(&definition).Error
	if err != nil {
		return err
	}

	if strings.Contains(definition, "WHERE") {
		return nil
	}

	if definition != "" {
		if err := database.Exec("DROP INDEX IF EXISTS " + index.name).Error; err != nil {
			return err
		}
	}

	statement := "CREATE UNIQUE INDEX IF NOT EXISTS " + index.name +
		" ON " + index.table + " (" + index.columns + ") WHERE deleted_at IS NULL"

	if err := database.Exec(statement).Error; err != nil {
		return err
	}

	logger.Info("DB_INDEX_NARROWED", logger.String("index", index.name))

	return nil
}
