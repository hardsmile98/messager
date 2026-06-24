package config

import (
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
}

var cfg Config

func parseInt(duration string) int {
	d, err := strconv.Atoi(duration)

	if err != nil {
		return 0
	}

	return d
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg = Config{
		DBUri:           os.Getenv("DB_URI"),
		RedisAddr:       os.Getenv("REDIS_ADDR"),
		GRPCPort:        os.Getenv("GRPC_PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AccessTokenTTL:  parseInt(os.Getenv("ACCESS_TOKEN_TTL_MINUTES")),
		RefreshTokenTTL: parseInt(os.Getenv("REFRESH_TOKEN_TTL_MINUTES")),
	}

	return &cfg, nil
}

func GetConfig() *Config {
	return &cfg
}
