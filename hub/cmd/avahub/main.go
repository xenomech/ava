package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ava/hub/internal/config"
	"ava/hub/internal/hub"
)

func main() {
	if err := run(); err != nil {
		slog.Error("HUB_FAILED", slog.String("error", err.Error()))
		slog.Info("HUB_STOPPED")
		os.Exit(1)
	}

	slog.Info("HUB_STOPPED")
}

func run() error {
	cfg := config.GetConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("HUB_STARTED",
		slog.String("device_name", cfg.DeviceName),
		slog.String("api_base_url", cfg.APIBaseURL),
		slog.String("state_file", cfg.StateFile),
	)

	if err := hub.New(cfg).Run(ctx); err != nil && ctx.Err() == nil {
		return err
	}

	return nil
}
