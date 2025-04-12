package client

import (
	"context"
	"log"

	"github.com/Glack134/pc_club/internal/app/client"
	"github.com/Glack134/pc_club/pkg/config"
	"github.com/Glack134/pc_club/pkg/logger"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadClientConfig("configs/client.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	log := logger.New(cfg.LogLevel)

	// Создание клиента
	pcClient, err := client.NewPcClubClient(log, cfg)
	if err != nil {
		log.Fatal("Failed to create client")
	}

	// Запуск основных компонентов
	ctx := context.Background()
	if err := pcClient.Run(ctx); err != nil {
		log.Fatal("Client error")
	}
}
