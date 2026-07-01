package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUri       string
	RedisAddr   string
	GRPCPort    string
	TLSCertFile string
	TLSKeyFile  string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBUri:       os.Getenv("DB_URI"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		GRPCPort:    os.Getenv("GRPC_PORT"),
		TLSCertFile: os.Getenv("TLS_CERT_FILE"),
		TLSKeyFile:  os.Getenv("TLS_KEY_FILE"),
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

	if (c.TLSCertFile == "") != (c.TLSKeyFile == "") {
		return errors.New("TLS_CERT_FILE and TLS_KEY_FILE must both be set or both be empty")
	}

	return nil
}
