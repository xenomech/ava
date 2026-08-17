package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ava/config"
	"ava/internal/controller"
	"ava/internal/middleware"
	"ava/internal/repository"
	"ava/internal/routes"
	"ava/internal/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.GetConfig()

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		BodyLimit:    4 * 1024 * 1024,
	})

	repo := repository.NewRepository(nil)
	service := services.NewService(repo)
	mw := middleware.NewMiddleware(service)
	ctrl := controller.NewController(service)

	routes.AddRoutes(app, ctrl, mw)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("SERVER_STARTED port=%s", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Printf("SERVER_START_ERROR error=%v", err)
		}
	}()

	<-ctx.Done()
	log.Print("SHUTDOWN_SIGNAL_RECEIVED")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Printf("SERVER_SHUTDOWN_ERROR error=%v", err)
	}

	log.Print("SERVER_STOPPED")
}
