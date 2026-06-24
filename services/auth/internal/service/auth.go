package service

import (
	"auth/internal/config"
	"auth/internal/repository"
	pb "github.com/hardsmile98/messager/sdk/auth/v1"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userRepo         *repository.UserRepo
	refreshTokenRepo *repository.RefreshTokenRepo
	config           *config.Config
}

func NewAuthService(
	userRepo *repository.UserRepo,
	refreshTokenRepo *repository.RefreshTokenRepo,
	config *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		config:           config,
	}
}
