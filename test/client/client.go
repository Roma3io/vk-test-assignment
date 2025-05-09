package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "vk-test-assignment/internal/proto/gen"
)

type ClientConfig struct {
	Server struct {
		Host string `yaml:"host" env-default:"localhost"`
		Port int    `yaml:"port" env-default:"50051"`
	} `yaml:"grpc"`
}

func Load() (*ClientConfig, error) {
	var cfg ClientConfig
	path := fetchConfig()
	if path == "" {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("failed to load config from env: %w", err)
		}
		fmt.Printf("No config file provided, using defaults or environment variables: %+v\n", cfg)
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

func main() {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Connecting to server at %s", serverAddr)

	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	client := pb.NewPubSubClient(conn)

	ctx := context.Background()
	stream, err := client.Subscribe(ctx, &pb.SubscribeRequest{Key: "test"})
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	go func() {
		for {
			event, err := stream.Recv()
			if err != nil {
				log.Printf("Subscription ended: %v", err)
				return
			}
			log.Printf("Received event: %s", event.Data)
		}
	}()

	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("message %d", i)
		_, err := client.Publish(ctx, &pb.PublishRequest{
			Key:  "test",
			Data: msg,
		})
		if err != nil {
			log.Printf("Publish failed: %v", err)
		} else {
			log.Printf("Published: %s", msg)
		}
		time.Sleep(1 * time.Second)
	}

	<-ctx.Done()
}
