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

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/Glack134/pc_club/internal/client"
	"github.com/Glack134/pc_club/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/ini.v1"
)

type Config struct {
	ServerAddress string `ini:"server_address"`
	PcID          string `ini:"pc_id"`
	AuthToken     string `ini:"auth_token,omitempty"`
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
`)
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	config := &Config{}
	if err := cfg.Section("client").MapTo(config); err != nil {
		return nil, fmt.Errorf("config mapping failed: %v", err)
	}

	if config.ServerAddress == "" {
		return nil, fmt.Errorf("server_address is required")
	}
	if config.PcID == "" {
		return nil, fmt.Errorf("pc_id is required")
	}

	return config, nil
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

func checkSession(grpcClient rpc.AdminServiceClient, pcID string) (bool, error) {
	resp, err := grpcClient.GetActiveSessions(context.Background(), &rpc.Empty{})
	if err != nil {
		return false, fmt.Errorf("failed to get sessions: %v", err)
	}

	for _, s := range resp.Sessions {
		if s.PcId == pcID && s.ExpiresAt > time.Now().Unix() {
			return true, nil
		}
	}
	return false, nil
}

func unlockPC() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
	case "linux":
		cmd = exec.Command("xdg-screensaver", "reset")
	}
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to unlock PC: %v", err)
	}
}

func lockPC() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32.exe", "user32.dll,LockWorkStation")
	case "linux":
		cmd = exec.Command("xdg-screensaver", "lock")
	}
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to lock PC: %v", err)
	}
}

func runSessionUI(grpcClient rpc.AdminServiceClient, pcID string) {
	a := app.New()
	w := a.NewWindow("PC Club - " + pcID)
	w.SetFullScreen(true)

	timeLeft := binding.NewString()
	timeLabel := widget.NewLabelWithData(timeLeft)
	timeLeft.Set("Checking session time...")

	updateTime := func() {
		resp, err := grpcClient.GetActiveSessions(context.Background(), &rpc.Empty{})
		if err != nil {
			timeLeft.Set("Error checking time")
			return
		}

		for _, s := range resp.Sessions {
			if s.PcId == pcID {
				remaining := time.Until(time.Unix(s.ExpiresAt, 0))
				timeLeft.Set(fmt.Sprintf("Time left: %v", remaining.Round(time.Second)))
				return
			}
		}
		timeLeft.Set("No active session")
	}

	lockBtn := widget.NewButton("Lock Now", func() {
		lockPC()
		w.Close()
	})

	w.SetContent(container.NewVBox(
		timeLabel,
		lockBtn,
	))

	// Первое обновление
	updateTime()

	// Периодическое обновление
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updateTime()
		}
	}()

	w.Show()
	a.Run()
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

	// Проверяем активную сессию
	hasSession, err := checkSession(grpcClient, config.PcID)
	if err != nil {
		log.Printf("Session check error: %v", err)
	}

	if hasSession {
		unlockPC()
		runSessionUI(grpcClient, config.PcID)
	} else {
		lockScreen := client.NewLockScreen(config.PcID)
		lockScreen.SetUnlockCallback(func() {
			unlockPC()
			runSessionUI(grpcClient, config.PcID)
		})
		lockScreen.Show()
		app.New().Run()
	}
}
