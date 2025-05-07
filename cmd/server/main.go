package server

import (
	"log"
	"vk-test-assignment/internal/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal("Could not load config: ", err)
	}
}
