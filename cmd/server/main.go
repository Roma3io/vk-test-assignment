package server

import (
	"go.uber.org/zap"
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
	cfg, err := config.Load("config.yaml")
	if err != nil {
		panic(err) //stub(need to make hierarchy like env < flag < default)
	}
	log := setupLogger(cfg.Env)
	defer log.Sync()
	bus := subpub.NewSubPub()
	server := pubsubservice.NewPubSubServer(bus, cfg, log)
	server.Start()

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
