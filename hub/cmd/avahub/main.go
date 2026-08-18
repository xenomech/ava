package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"ava/hub/internal/config"
)

func main() {
	cfg := config.GetConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	slog.Info("HUB_STARTED",
		slog.String("device_id", cfg.DeviceID),
		slog.String("api_base_url", cfg.APIBaseURL),
		slog.String("mqtt_broker_url", cfg.MQTTBrokerURL),
	)

	<-ctx.Done()

	slog.Info("SHUTDOWN_SIGNAL_RECEIVED")
	slog.Info("HUB_STOPPED")
}
