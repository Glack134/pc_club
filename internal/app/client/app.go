package client

import (
	"context"
	"time"

	"github.com/Glack134/pc_club/internal/pkg/locker"
	"github.com/Glack134/pc_club/internal/pkg/monitoring"
	"github.com/Glack134/pc_club/pkg/config"
	"github.com/Glack134/pc_club/pkg/logger"
)

type PcClubClient struct {
	logger     logger.Logger
	config     *config.ClientConfig
	locker     locker.Locker
	monitor    monitoring.HardwareMonitor
	process    monitoring.ProcessTracker
	pcId       string
	serverAddr string
}

func NewPcClubClient(log logger.Logger, cfg *config.ClientConfig) (*PcClubClient, error) {
	locker, err := createLocker()
	if err != nil {
		return nil, err
	}

	monitor := monitoring.NewHardwareMonitor()
	process := monitoring.NewProcessTracker()

	return &PcClubClient{
		logger:     log,
		config:     cfg,
		locker:     locker,
		monitor:    monitor,
		process:    process,
		pcId:       cfg.PcId,
		serverAddr: cfg.ServerAddr,
	}, nil
}

func (c *PcClubClient) Run(ctx context.Context) error {
	go c.monitoringLoop(ctx)
	go c.commandLoop(ctx)

	<-ctx.Done()
	return nil
}

func (c *PcClubClient) monitoringLoop(ctx context.Context) {
	ticker := time.NewTicker(c.config.Monitoring.ReportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cpu, _ := c.monitor.GetCPUUsage()
			ram, _ := c.monitor.GetRAMUsage()

			c.logger.Debug("System stats",
				logger.Field{Key: "cpu", Value: cpu},
				logger.Field{Key: "ram", Value: ram},
			)

		case <-ctx.Done():
			return
		}
	}
}

func (c *PcClubClient) commandLoop(ctx context.Context) {
	// Реализация обработки команд от сервера
}

func createLocker() (locker.Locker, error) {
	return locker.NewWindowsLocker(), nil
}
