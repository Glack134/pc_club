package client

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/Glack134/pc_club/internal/pkg/locker"
	"github.com/Glack134/pc_club/internal/pkg/monitoring"
	"github.com/Glack134/pc_club/pkg/config"
	"github.com/Glack134/pc_club/pkg/logger"
	"google.golang.org/grpc"
)

type PcClubClient struct {
	logger  logger.Logger
	config  *config.ClientConfig
	locker  locker.Locker
	monitor monitoring.HardwareMonitor
	conn    *grpc.ClientConn
	pcID    string
}

func NewPcClubClient(log logger.Logger, cfg *config.ClientConfig) (*PcClubClient, error) {
	return &PcClubClient{
		logger: log,
		config: cfg,
	}, nil
}

func NewClient(logger logger.Logger) *PcClubClient {
	return &PcClubClient{
		logger:  logger,
		locker:  locker.NewWindowsLocker(),
		monitor: monitoring.NewHardwareMonitor(),
		pcID:    generatePCID(),
	}
}

func (c *PcClubClient) AutoDiscover() (string, error) {
	// Сканируем локальную сеть
	ips, _ := net.LookupIP("pcclub._tcp.local")
	for _, ip := range ips {
		if c.testConnection(ip.String()) {
			return ip.String(), nil
		}
	}
	return "", errors.New("server not found")
}

func (c *PcClubClient) testConnection(ip string) bool {
	conn, err := net.DialTimeout("tcp", ip+":50051", 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (c *PcClubClient) Run(ctx context.Context) error {
	// реализация
	return nil
}

func (c *PcClubClient) Connect(serverIP string) error {
	conn, err := grpc.Dial(serverIP+":50051", grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.conn = conn
	c.logger.Info("Connected to server", logger.Field{Key: "server", Value: serverIP})
	return nil
}

func generatePCID() string {
	interfaces, _ := net.Interfaces()
	for _, i := range interfaces {
		if i.HardwareAddr != nil {
			return i.HardwareAddr.String()
		}
	}
	return "unknown_" + time.Now().Format("20060102150405")
}
