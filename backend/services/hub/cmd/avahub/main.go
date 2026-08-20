package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"ava/hub/internal/app"
	"ava/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		logger.Error("HUB_FAILED", logger.Err(err))
		logger.Sync()
		os.Exit(1)
	}

	logger.Info("HUB_STOPPED")
	logger.Sync()
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	instance, err := app.Bootstrap(ctx)
	if err != nil {
		return err
	}

	defer instance.Close()

	if err := instance.Run(ctx); err != nil && ctx.Err() == nil {
		return err
	}

	return nil
}
