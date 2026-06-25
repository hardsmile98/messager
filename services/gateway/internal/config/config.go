package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	AuthGRPCURL string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:        os.Getenv("PORT"),
		AuthGRPCURL: os.Getenv("AUTH_GRPC_URL"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Port == "" {
		return errors.New("PORT is required")
	}

	if c.AuthGRPCURL == "" {
		return errors.New("AUTH_GRPC_URL is required")
	}

	return nil
}
