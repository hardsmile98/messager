package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUri           string
	RedisAddr       string
	JWTPrivateKey   string
	JWTPublicKey    string
	AccessTokenTTL  int
	RefreshTokenTTL int
}

var cfg Config

func parseInt(duration string) int {
	d, err := strconv.Atoi(duration)

	if err != nil {
		return 0
	}

	if d < 0 {
		return 0
	}

	return d
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg = Config{
		DBUri:           os.Getenv("DB_URI"),
		RedisAddr:       os.Getenv("REDIS_ADDR"),
		JWTPrivateKey:   os.Getenv("JWT_PRIVATE_KEY"),
		JWTPublicKey:    os.Getenv("JWT_PUBLIC_KEY"),
		AccessTokenTTL:  parseInt(os.Getenv("ACCESS_TOKEN_TTL_MINUTES")),
		RefreshTokenTTL: parseInt(os.Getenv("REFRESH_TOKEN_TTL_MINUTES")),
	}

	return &cfg, nil
}

func GetConfig() *Config {
	return &cfg
}
