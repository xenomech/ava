package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ava/api/internal/config"
	"ava/api/internal/controller"
	"ava/api/internal/db"
	"ava/api/internal/middleware"
	"ava/api/internal/repository"
	"ava/api/internal/routes"
	"ava/api/internal/services"
	"ava/api/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSAllowedOrigins,
		AllowMethods:     cfg.CORSAllowedMethods,
		AllowHeaders:     cfg.CORSAllowedHeaders,
		AllowCredentials: true,
		MaxAge:           cfg.CORSMaxAge,
	}))

	repo := repository.NewRepository(database)
	service := services.NewService(repo)
	mw := middleware.NewMiddleware(service)
	ctrl := controller.NewController(service)

	app.Use(mw.RequestTrace)
	app.Use(middleware.SecurityHeaders(cfg.ServerEnv))

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
