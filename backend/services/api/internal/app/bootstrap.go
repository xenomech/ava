package app

import (
	"context"
	"time"

	"ava/api/config"
	"ava/api/internal/broker"
	"ava/api/internal/controller"
	"ava/api/internal/db"
	"ava/api/internal/middleware"
	"ava/api/internal/repository"
	"ava/api/internal/routes"
	"ava/api/internal/services"
	"ava/pkg/logger"

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
	Config  *config.Config
	DB      *gorm.DB
	Service *services.Service
	Broker  *broker.Broker
	Fiber   *fiber.App
}

func Bootstrap(ctx context.Context) (*App, error) {
	cfg := config.GetConfig()

	logger.Init(cfg.ServerEnv, cfg.LogLevel)

	database, err := db.Connect(&db.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Database: cfg.DBDatabase,
		SSLMode:  cfg.DBSSLMode(),
	})
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(database); err != nil {
		disconnect(database)

		return nil, err
	}

	messages, err := broker.Connect(ctx, broker.Config{
		URL:      cfg.MQTTBrokerURL,
		Username: cfg.MQTTUsername,
		Password: cfg.MQTTPassword,
	})
	if err != nil {
		logger.Warn("MQTT_UNAVAILABLE", logger.Err(err))
	}

	var (
		commander   services.Commander
		provisioner services.Provisioner
	)

	if messages != nil {
		if err := messages.EnsureControlPlane(ctx, cfg.MQTTUsername); err != nil {
			logger.Warn("MQTT_CONTROL_PLANE_FAILED", logger.Err(err))
		}

		commander = messages
		provisioner = messages
	}

	service := services.NewService(repository.NewRepository(database), commander, provisioner)

	return &App{
		Config:  cfg,
		DB:      database,
		Service: service,
		Broker:  messages,
		Fiber:   buildFiber(ctx, cfg, service),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.Broker.Listen(ctx, a.Service.Device.ApplyReportedState, a.Service.Hub.ApplyPresence)

	listenErr := make(chan error, 1)

	go func() {
		logger.Info("SERVER_STARTED", logger.String("port", a.Config.Port))

		if err := a.Fiber.Listen(":" + a.Config.Port); err != nil {
			listenErr <- err
		}
	}()

	select {
	case err := <-listenErr:
		logger.Error("SERVER_START_ERROR", logger.Err(err))

		return err
	case <-ctx.Done():
	}

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

	a.Broker.Close()
	disconnect(a.DB)

	logger.Sync()
}

func disconnect(database *gorm.DB) {
	if err := db.Disconnect(database); err != nil {
		logger.Error("DB_DISCONNECT_ERROR", logger.Err(err))
	}
}

func buildFiber(life context.Context, cfg *config.Config, service *services.Service) *fiber.App {
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

	routes.AddRoutes(server, controller.NewController(life, service), mw)

	return server
}
