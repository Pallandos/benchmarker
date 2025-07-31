package config

import (
	"os"

	"github.com/joho/godotenv"
)

// AppConfig holds the configuration for the application
type AppConfig struct {
	StackName string `env:"STACK_NAME"`
	LogPath   string `env:"LOG_PATH"`
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

	return &config, nil
}
