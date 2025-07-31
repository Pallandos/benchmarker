package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	StackName string `env:"STACK_NAME"`
}

func LoadConfig(path string) (*AppConfig, error) {

	var config AppConfig
	err := godotenv.Load(path)

	if err != nil {
		return nil, err
	}

	config.StackName = os.Getenv("STACK_NAME")
	return &config, nil
}
