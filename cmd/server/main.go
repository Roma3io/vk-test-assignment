package main

import (
	"go.uber.org/zap"
	"log"
	"vk-test-assignment/internal/config"
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
	//bus := subpub.NewSubPub()
	//server := pubsubservice.NewPubSubServer(bus, cfg, log)
	//server.Start()
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
