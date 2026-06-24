package service

import (
	"auth/internal/config"
	"auth/internal/jwt"
	"auth/internal/model"
	"auth/internal/repository"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *AuthService) generateRefreshToken() (string, string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)

	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	hash := sha256.Sum256([]byte(token))

	tokenHash := hex.EncodeToString(hash[:])

	return token, tokenHash, nil
}

func (s *AuthService) generateTokens(ctx context.Context, userID string) (string, string, error) {
	accessToken, err := jwt.GenerateToken(userID, s.config.JWTSecret, s.config.AccessTokenTTL)

	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refresh, refreshHash, err := s.generateRefreshToken()

	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	expiresAt := time.Now().Add(time.Duration(s.config.RefreshTokenTTL) * time.Minute)

	newRefreshToken := &model.RefreshToken{
		UserID:    userID,
		TokenHash: refreshHash,
		Device:    "DEVICE_ID",
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	err = s.refreshTokenRepo.SaveRefreshToken(ctx, newRefreshToken)

	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to save refresh token: %v", err)
	}

	return accessToken, refresh, nil
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email and password are required")
	}

	exists, err := s.userRepo.GetUserByUsername(ctx, req.Username)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check if username exists: %v", err)
	}

	if exists != nil {
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	id, err := s.userRepo.CreateUser(ctx, user)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	return &pb.RegisterResponse{
		UserId:       id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
