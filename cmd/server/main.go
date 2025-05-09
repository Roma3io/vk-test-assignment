package main

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	logger, err := setupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("Failed to setup logger: %v", err)
	}
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

func setupLogger(env string) (*zap.Logger, error) {
	var cfg zap.Config
	switch env {
	case envProd:
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case envLocal, envDev:
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	default:
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
		log.Printf("Unknown environment '%s', using default logger configuration", env)
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
