package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	GRPCAddr string         `yaml:"grpc_addr"`
	HTTPAddr string         `yaml:"http_addr"`
	LogLevel string         `yaml:"log_level"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

type ClientConfig struct {
	ServerAddr string           `yaml:"server_addr"`
	PcId       string           `yaml:"pc_id"`
	LogLevel   string           `yaml:"log_level"`
	Lock       LockConfig       `yaml:"lock"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type AuthConfig struct {
	SecretKey     string        `yaml:"secret_key"`
	TokenDuration time.Duration `yaml:"token_duration"`
}

type LockConfig struct {
	CheckInterval time.Duration `yaml:"check_interval"`
}

type MonitoringConfig struct {
	ReportInterval time.Duration `yaml:"report_interval"`
	MaxCPUUsage    float64       `yaml:"max_cpu_usage"`
	MaxGPUUsage    float64       `yaml:"max_gpu_usage"`
}

func LoadServerConfig(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadClientConfig(path string) (*ClientConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ClientConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
