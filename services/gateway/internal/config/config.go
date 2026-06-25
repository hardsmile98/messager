package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  string
	AuthGRPCURL           string
	HTTPReadTimeout       time.Duration
	HTTPWriteTimeout      time.Duration
	HTTPIdleTimeout       time.Duration
	HTTPReadHeaderTimeout time.Duration
	CookieSecure          bool
	CookieDomain          string
	CORSAllowedOrigins    []string
}

func parseDurationSeconds(envKey string, defaultSeconds int) (time.Duration, error) {
	value := os.Getenv(envKey)
	if value == "" {
		return time.Duration(defaultSeconds) * time.Second, nil
	}

	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a positive integer", envKey)
	}

	if seconds <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", envKey)
	}

	return time.Duration(seconds) * time.Second, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	httpReadTimeout, err := parseDurationSeconds("HTTP_READ_TIMEOUT_SECONDS", 10)
	if err != nil {
		return nil, err
	}

	httpWriteTimeout, err := parseDurationSeconds("HTTP_WRITE_TIMEOUT_SECONDS", 10)
	if err != nil {
		return nil, err
	}

	httpIdleTimeout, err := parseDurationSeconds("HTTP_IDLE_TIMEOUT_SECONDS", 60)
	if err != nil {
		return nil, err
	}

	httpReadHeaderTimeout, err := parseDurationSeconds("HTTP_READ_HEADER_TIMEOUT_SECONDS", 5)
	if err != nil {
		return nil, err
	}

	cookieSecure := os.Getenv("COOKIE_SECURE") == "true"

	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")

	cfg := &Config{
		Port:                  os.Getenv("PORT"),
		AuthGRPCURL:           os.Getenv("AUTH_GRPC_URL"),
		HTTPReadTimeout:       httpReadTimeout,
		HTTPWriteTimeout:      httpWriteTimeout,
		HTTPIdleTimeout:       httpIdleTimeout,
		HTTPReadHeaderTimeout: httpReadHeaderTimeout,
		CookieSecure:          cookieSecure,
		CookieDomain:          os.Getenv("COOKIE_DOMAIN"),
		CORSAllowedOrigins:    splitCSV(corsOrigins),
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
