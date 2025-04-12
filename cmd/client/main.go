package client

import (
	"context"
	"log"
)

func main() {
	cfg, err := config.LoadClientConfig("configs/client.yaml")
	if err != nil {
		log.Fatal("Failed to load config: %v", err)
	}

	logger := logger.New(cfg.LogLevel)

	pcClient, err := client.NewPcClubClient(logger, cfg)
	if err != nil {
		logger.Fatal("Failed to create client", "error", err)
	}

	ctx := context.Background()
	if err := pcClient.Run(ctx); err != nil {
		logger.Fatal("Client error", "error", err)
	}
}
