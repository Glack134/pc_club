package server

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadServerConfig("config/server.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.New(cfg.LogLevel)

	storage, err := server.NewStorage(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to initialize storage", "error", err)
	}

	grpcServer := grpc.NewServer()
	server := server.NewPcClubServer(logger, storage, cfg)

	server.RegisterServices(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		logger.Fatal("Failed to listen", "address", cfg.GRPCAddr, "error", err)
	}

	logger.Info("Server started", "address", cfg.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", "error", err)
	}
}
