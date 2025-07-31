package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds the configuration for the application
type AppConfig struct {
	StackName       string        `env:"STACK_NAME"`
	LogPath         string        `env:"LOG_PATH"`
	MonitorInterval time.Duration `env:"MONITOR_INTERVAL"`
	MonitorDuration time.Duration `env:"MONITOR_DURATION"`
}

// LoadConfig loads the application configuration from a .env file
func LoadConfig(path string) (*AppConfig, error) {

	var config AppConfig
	err := godotenv.Load(path)

	if err != nil {
		return nil, err
	}

	config.StackName = os.Getenv("STACK_NAME")
	config.LogPath = os.Getenv("LOG_PATH")

	config.MonitorInterval = 5 * time.Second
	if interval := os.Getenv("MONITOR_INTERVAL"); interval != "" {
		if d, err := time.ParseDuration(interval); err == nil {
			config.MonitorInterval = d
		}
	}

	config.MonitorDuration = 60 * time.Second
	if duration := os.Getenv("MONITOR_DURATION"); duration != "" {
		if d, err := time.ParseDuration(duration); err == nil {
			config.MonitorDuration = d
		}
	}

	return &config, nil
}
