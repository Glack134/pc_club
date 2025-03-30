package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/ini.v1"
)

type Config struct {
	ServerAddress string `ini:"server_address"`
	PcID          string `ini:"pc_id"`
	AuthToken     string `ini:"auth_token"`
}

func loadConfig() (*Config, error) {
	// Получаем абсолютный путь к config.ini
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}
	configPath := filepath.Join(filepath.Dir(exePath), "config.ini")

	// Проверяем существование файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config.ini not found in: %s", filepath.Dir(exePath))
	}

	// Проверяем, что это текстовый файл
	if isBinary(configPath) {
		return nil, fmt.Errorf("config.ini is a binary file")
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	config := new(Config)
	if err := cfg.MapTo(config); err != nil {
		return nil, fmt.Errorf("config mapping error: %v", err)
	}

	return config, nil
}

func isBinary(filepath string) bool {
	file, err := os.Open(filepath)
	if err != nil {
		return true
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return true
	}

	return strings.Contains(http.DetectContentType(buffer), "text/plain") == false
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Принудительная проверка адреса
	if config.ServerAddress == "" || config.ServerAddress == "localhost:50051" {
		log.Fatal("Server address not configured properly in config.ini")
	}

	// Настройка gRPC клиента с явным указанием IPv4
	conn, err := grpc.Dial(
		config.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "tcp4", addr) // Явно указываем IPv4
		}),
	)
	if err != nil {
		log.Fatalf("Connection failed to %s: %v", config.ServerAddress, err)
	}
	defer conn.Close()

	client := rpc.NewAdminServiceClient(conn)

	// Тестовый запрос
	resp, err := client.GrantAccess(context.Background(), &rpc.GrantRequest{
		UserId:    "test",
		PcId:      config.PcID,
		Minutes:   60,
		AuthToken: config.AuthToken,
	})
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	log.Printf("Server response: %v", resp.Message)
}
