package app

import (
	"context"
	"time"

	"ava/api/config"
	"ava/api/internal/controller"
	"ava/api/internal/db"
	"ava/api/internal/middleware"
	"ava/api/internal/repository"
	"ava/api/internal/routes"
	"ava/api/internal/services"
	"ava/pkg/logger"
	"ava/pkg/mqtt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

const (
	readTimeout     = 10 * time.Second
	idleTimeout     = 120 * time.Second
	shutdownTimeout = 10 * time.Second
	bodyLimit       = 4 * 1024 * 1024
)

type App struct {
	Config    *config.Config
	DB        *gorm.DB
	Service   *services.Service
	Publisher *mqtt.Client
	Fiber     *fiber.App
}

func Bootstrap(ctx context.Context) (*App, error) {
	cfg := config.GetConfig()

	logger.Init(cfg.ServerEnv, cfg.LogLevel)

	database, err := connectDatabase(cfg)
	if err != nil {
		return nil, err
	}

	publisher := connectBroker(ctx, cfg)

	var commander services.Commander
	if publisher != nil {
		commander = publisher
	}

	service := services.NewService(repository.NewRepository(database), commander)

	return &App{
		Config:    cfg,
		DB:        database,
		Service:   service,
		Publisher: publisher,
		Fiber:     buildFiber(cfg, service),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		logger.Info("SERVER_STARTED", logger.String("port", a.Config.Port))

		if err := a.Fiber.Listen(":" + a.Config.Port); err != nil {
			logger.Error("SERVER_START_ERROR", logger.Err(err))
		}
	}()

	<-ctx.Done()
	logger.Info("SHUTDOWN_SIGNAL_RECEIVED")

	if err := a.Fiber.ShutdownWithTimeout(shutdownTimeout); err != nil {
		logger.Error("SERVER_SHUTDOWN_ERROR", logger.Err(err))

		return err
	}

	return nil
}

func (a *App) Close() {
	if a == nil {
		return
	}

	a.Publisher.Close()

	if err := db.Disconnect(a.DB); err != nil {
		logger.Error("DB_DISCONNECT_ERROR", logger.Err(err))
	}

	logger.Sync()
}

func buildFiber(cfg *config.Config, service *services.Service) *fiber.App {
	server := fiber.New(fiber.Config{
		ReadTimeout: readTimeout,
		IdleTimeout: idleTimeout,
		BodyLimit:   bodyLimit,
	})

	server.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     cfg.CORSAllowedMethods,
		AllowHeaders:     cfg.CORSAllowedHeaders,
		AllowCredentials: true,
		MaxAge:           cfg.CORSMaxAge,
	}))

	mw := middleware.NewMiddleware(service)

	server.Use(mw.RequestTrace)
	server.Use(middleware.SecurityHeaders(cfg.ServerEnv))
	server.Use(middleware.VerifyOrigin(cfg.CORSAllowedOrigins))

	routes.AddRoutes(server, controller.NewController(service), mw)

	return server
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	sslMode := "require"
	if cfg.ServerEnv == "local" {
		sslMode = "disable"
	}

	database, err := db.Connect(&db.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBDatabase,
		SSLMode:  sslMode,
	})
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(database); err != nil {
		return nil, err
	}

	return database, nil
}

func connectBroker(ctx context.Context, cfg *config.Config) *mqtt.Client {
	publisher, err := mqtt.Connect(ctx, &mqtt.Options{
		BrokerURL: cfg.MQTTBrokerURL,
		ClientID:  "ava-api",
	})
	if err != nil {
		logger.Warn("MQTT_UNAVAILABLE", logger.Err(err))

		return nil
	}

	return publisher
}
