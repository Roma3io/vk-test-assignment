package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	GRPCServer GRPCServer `yaml:"grpc_server"`
}

type GRPCServer struct {
	Port    int           `yaml:"port" env-default:"50051"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func Load() (*Config, error) {
	var cfg Config
	path := fetchConfig()
	if path == "" {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("failed to load config from env: %w", err)
		}
		fmt.Printf("No config file provided, using defaults or environment variables: %+v", cfg)
		return &cfg, nil
	}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}
	return &cfg, nil
}

func fetchConfig() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
