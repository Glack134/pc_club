package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/Glack134/pc_club/internal/auth"
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

func main() {
	config := loadConfig()

	// При запуске проверяем токен из конфига
	if config.AuthToken != "" {
		claims, err := auth.ValidateToken(config.AuthToken)
		if err == nil && claims.PCID == config.PcID {
			unlockPC()
			showSessionUI()
			return
		}
	}

	conn, err := grpc.Dial(config.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gRPC connection failed: %v", err)
	}
	defer conn.Close()

	grpcClient := rpc.NewAdminServiceClient(conn)

	// Проверка активной сессии
	if hasActiveSession(grpcClient, config.PcID) {
		unlockPC() // Разблокируем PC если есть активная сессия
		showSessionUI(grpcClient, config)
	} else {
		lockScreen := client.NewLockScreen(config.PcID)
		lockScreen.SetUnlockCallback(func() {
			unlockPC()
			showSessionUI(grpcClient, config)
		})
		lockScreen.Run()
	}
}

func unlockPC() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
	case "linux":
		cmd = exec.Command("xdg-open", "steam://")
	}
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to unlock PC: %v", err)
	}
}

func showSessionUI(client rpc.AdminServiceClient, config *Config) {
	a := app.New()
	w := a.NewWindow("PC Club - " + config.PcID)
	w.SetFullScreen(true)

	timeLeft := binding.NewString()
	timeLabel := widget.NewLabelWithData(timeLeft)

	// Обновление времени
	go func() {
		for {
			resp, err := client.GetActiveSessions(context.Background(), &rpc.Empty{})
			if err == nil {
				for _, s := range resp.Sessions {
					if s.PcId == config.PcID {
						remaining := time.Until(time.Unix(s.ExpiresAt, 0))
						timeLeft.Set(fmt.Sprintf("Time left: %v", remaining.Round(time.Second)))
					}
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

	w.SetContent(container.NewCenter(
		container.NewVBox(
			timeLabel,
			widget.NewButton("Lock Now", func() {
				lockPC()
				w.Close()
			}),
		),
	))
	w.Show()
}

func lockPC() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
	case "linux":
		cmd = exec.Command("loginctl", "lock-session")
	}
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to lock PC: %v", err)
	}
}
