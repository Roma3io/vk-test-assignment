package server

import (
	"go.uber.org/zap"
	"vk-test-assignment/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.Load("config.yaml")
	log := setupLogger(cfg.Env)
	defer log.Sync()

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
