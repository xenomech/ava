package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ava/config"
	"ava/internal/controller"
	"ava/internal/db"
	"ava/internal/middleware"
	"ava/internal/repository"
	"ava/internal/routes"
	"ava/internal/services"
	"ava/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.GetConfig()

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
		logger.Error("DB_CONNECTION_ERROR", logger.Err(err))
		panic(err)
	}

	if err := db.Migrate(database); err != nil {
		logger.Error("DB_MIGRATION_ERROR", logger.Err(err))
		panic(err)
	}

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		BodyLimit:    4 * 1024 * 1024,
	})

	repo := repository.NewRepository(database)
	service := services.NewService(repo)
	mw := middleware.NewMiddleware(service)
	ctrl := controller.NewController(service)

	routes.AddRoutes(app, ctrl, mw)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("SERVER_STARTED", logger.String("port", cfg.Port))

		if err := app.Listen(":" + cfg.Port); err != nil {
			logger.Error("SERVER_START_ERROR", logger.Err(err))
		}
	}()

	<-ctx.Done()
	logger.Info("SHUTDOWN_SIGNAL_RECEIVED")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		logger.Error("SERVER_SHUTDOWN_ERROR", logger.Err(err))
	}

	if err := db.Disconnect(database); err != nil {
		logger.Error("DB_DISCONNECT_ERROR", logger.Err(err))
	}

	logger.Info("SERVER_STOPPED")
	logger.Sync()
}
