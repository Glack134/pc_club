package main

import (
	"log"
	"net"

	"github.com/Glack134/pc_club/internal/app/server"
	"github.com/Glack134/pc_club/internal/storage/sqlite"
	"github.com/Glack134/pc_club/pkg/config"
	"github.com/Glack134/pc_club/pkg/logger"
	"google.golang.org/grpc"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadServerConfig("configs/server.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	log := logger.New(cfg.LogLevel) // Переименовали переменную в log

	// Инициализация хранилища
	storage, err := sqlite.NewStorage(cfg.Database.DSN)
	if err != nil {
		log.Fatal("Failed to initialize storage", logger.Field{Key: "error", Value: err})
	}

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	srv := server.NewPcClubServer(log, storage)

	// Регистрация сервисов
	srv.RegisterServices(grpcServer)

	// Запуск сервера
	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatal("Failed to listen",
			logger.Field{Key: "address", Value: cfg.GRPCAddr},
			logger.Field{Key: "error", Value: err},
		)
	}

	log.Info("Server started", logger.Field{Key: "address", Value: cfg.GRPCAddr})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve", logger.Field{Key: "error", Value: err})
	}
}
