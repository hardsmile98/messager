package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUri           string
	RedisAddr       string
	JWTSecret       string
	AccessTokenTTL  int
	RefreshTokenTTL int
	GRPCPort        string
	TLSCertFile     string
	TLSKeyFile      string
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

	accessTokenTTL, err := parseInt(os.Getenv("ACCESS_TOKEN_TTL_MINUTES"))

	if err != nil {
		return nil, fmt.Errorf("ACCESS_TOKEN_TTL_MINUTES: %w", err)
	}

	refreshTokenTTL, err := parseInt(os.Getenv("REFRESH_TOKEN_TTL_MINUTES"))

	if err != nil {
		return nil, fmt.Errorf("REFRESH_TOKEN_TTL_MINUTES: %w", err)
	}

	cfg := &Config{
		DBUri:           os.Getenv("DB_URI"),
		RedisAddr:       os.Getenv("REDIS_ADDR"),
		GRPCPort:        os.Getenv("GRPC_PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		TLSCertFile:     os.Getenv("TLS_CERT_FILE"),
		TLSKeyFile:      os.Getenv("TLS_KEY_FILE"),
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

	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}

	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters")
	}

	if c.AccessTokenTTL <= 0 {
		return errors.New("ACCESS_TOKEN_TTL_MINUTES must be greater than 0")
	}

	if c.RefreshTokenTTL <= 0 {
		return errors.New("REFRESH_TOKEN_TTL_MINUTES must be greater than 0")
	}

	if (c.TLSCertFile == "") != (c.TLSKeyFile == "") {
		return errors.New("TLS_CERT_FILE and TLS_KEY_FILE must both be set or both be empty")
	}

	return nil
}
