package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/Glack134/pc_club/internal/client"
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
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %v", err)
	}
	configPath := filepath.Join(filepath.Dir(exePath), "config.ini")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := []byte(`[client]
server_address = 192.168.1.14:50051
pc_id = default-pc
auth_token = `)

		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		return &Config{
			ServerAddress: "192.168.1.14:50051",
			PcID:          "default-pc",
		}, nil
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	config := &Config{}
	if err := cfg.Section("client").MapTo(config); err != nil {
		return nil, fmt.Errorf("config mapping failed: %v", err)
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *Config) error {
	if config.ServerAddress == "" {
		return fmt.Errorf("server_address is required")
	}
	if config.PcID == "" {
		return fmt.Errorf("pc_id is required")
	}
	if config.AuthToken == "" {
		return fmt.Errorf("auth_token is required")
	}
	return nil
}

func createGRPCConnection(address string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}

	return conn, nil
}

func hasActiveSession(client rpc.AdminServiceClient, pcID string) bool {
	resp, err := client.GetActiveSessions(context.Background(), &rpc.Empty{})
	if err != nil {
		return false
	}

	for _, s := range resp.Sessions {
		if s.PcId == pcID && s.ExpiresAt > time.Now().Unix() {
			return true
		}
	}
	return false
}

func startMainUI(client rpc.AdminServiceClient, config *Config) {
	a := app.New()
	w := a.NewWindow("PC Club Client - " + config.PcID)

	// Таймер сессии
	sessionLabel := widget.NewLabel("Active session")
	w.SetContent(container.NewCenter(sessionLabel))

	w.Show()
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Config initialization failed: %v", err)
	}

	conn, err := createGRPCConnection(config.ServerAddress)
	if err != nil {
		log.Fatalf("gRPC connection failed: %v", err)
	}
	defer conn.Close()

	grpcClient := rpc.NewAdminServiceClient(conn)

	// Создаем экран блокировки
	lockScreen := client.NewLockScreen()
	lockScreen.SetUnlockCallback(func() {
		startMainUI(grpcClient, config)
	})

	if hasActiveSession(grpcClient, config.PcID) {
		startMainUI(grpcClient, config)
	} else {
		lockScreen.Show()
	}

	// Запускаем главный цикл приложения
	lockScreen.Window.ShowAndRun()
}
