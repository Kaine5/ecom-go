package main

import (
	"context"
	"ecom-go/internal/repository"
	"os"
	"os/signal"
	"syscall"

	"ecom-go/internal/config"
	"ecom-go/pkg/logger"
)

func main() {
	logger.Info("Starting worker service")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", "error", err)
	}

	// Set up repository
	repoFactory, err := repository.NewFactory(cfg)
	if err != nil {
		logger.Fatal("Failed to create repository factory", "error", err)
	}
	defer func() {
		if err := repoFactory.Close(); err != nil {
			logger.Error("Error closing database connection", "error", err)
		}
	}()

	// Set up services
	// TODO: Add other services here

	// Create a context that is canceled when a signal is received
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start workers (to be implemented)
	// go startOrderProcessor(ctx, log, cfg)
	// go startNotificationProcessor(ctx, log, cfg)

	logger.Info("Worker service started")

	// Wait for signal
	sig := <-sigCh
	logger.Info("Received signal, shutting down", "signal", sig)

	// Give workers a chance to finish gracefully
	cancel()
	logger.Info("Worker service stopped")
}
