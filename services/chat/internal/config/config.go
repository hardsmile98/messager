package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUri     string
	RedisAddr string
	GRPCPort  string
}

func parseInt(value string) (int, error) {
	if value == "" {
		return 0, errors.New("value is required")
	}

	d, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	return d, nil
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBUri:     os.Getenv("DB_URI"),
		RedisAddr: os.Getenv("REDIS_ADDR"),
		GRPCPort:  os.Getenv("GRPC_PORT"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DBUri == "" {
		return errors.New("DB_URI is required")
	}

	if c.RedisAddr == "" {
		return errors.New("REDIS_ADDR is required")
	}

	if c.GRPCPort == "" {
		return errors.New("GRPC_PORT is required")
	}

	return nil
}
