package main

import (
	"context"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"syscall"
	"vk-test-assignment/internal/config"
	"vk-test-assignment/internal/pubsubservice"
	"vk-test-assignment/pkg/subpub"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error %v", err)
	}
	logger := setupLogger(cfg.Env)
	defer logger.Sync()
	bus := subpub.NewSubPub()
	server := pubsubservice.NewPubSubServer(bus, cfg, logger)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	if err := server.Start(ctx); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
	<-ctx.Done()
	logger.Info("Server stopped")
}

func setupLogger(env string) *zap.Logger {
	var logger *zap.Logger
	switch env {
	case envProd:
		logger, _ = zap.NewProduction()
	case envLocal:
		logger, _ = zap.NewDevelopment()
	case envDev:
		logger, _ = zap.NewDevelopment()
	}
	return logger
}
