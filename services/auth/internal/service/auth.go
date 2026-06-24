package service

import (
	"auth/internal/config"
	"auth/internal/errors"
	"auth/internal/jwt"
	"auth/internal/model"
	"auth/internal/repository"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	minPasswordLength = 8
	defaultDeviceID   = "unknown"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
	config           *config.Config
}

func NewAuthService(
	userRepo repository.UserRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	config *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		config:           config,
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}

func generateRefreshToken() (string, string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	return token, hashToken(token), nil
}

func (s *AuthService) generateTokens(ctx context.Context, userID, device string) (string, string, error) {
	accessToken, err := jwt.GenerateToken(userID, s.config.JWTSecret, s.config.AccessTokenTTL)

	if err != nil {
		return "", "", errors.InternalError(ctx, "failed to generate access token", err)
	}

	refresh, refreshHash, err := generateRefreshToken()

	if err != nil {
		return "", "", errors.InternalError(ctx, "failed to generate refresh token", err)
	}

	now := time.Now()

	newRefreshToken := &model.RefreshToken{
		UserID:    userID,
		TokenHash: refreshHash,
		Device:    device,
		ExpiresAt: now.Add(time.Duration(s.config.RefreshTokenTTL) * time.Minute),
		CreatedAt: now,
	}

	if err := s.refreshTokenRepo.SaveRefreshToken(ctx, newRefreshToken); err != nil {
		return "", "", errors.InternalError(ctx, "failed to save refresh token", err)
	}

	return accessToken, refresh, nil
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.InvalidArgument("username, email and password are required")
	}

	exists, err := s.userRepo.GetUserByUsername(ctx, req.Username)

	if err != nil {
		return nil, errors.InternalError(ctx, "failed to check if username exists", err)
	}

	if exists != nil {
		return nil, errors.AlreadyExists("username already exists")
	}

	emailExists, err := s.userRepo.GetUserByEmail(ctx, req.Email)

	if err != nil {
		return nil, errors.InternalError(ctx, "failed to check if email exists", err)
	}

	if emailExists != nil {
		return nil, errors.AlreadyExists("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.InternalError(ctx, "failed to hash password", err)
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	id, err := s.userRepo.CreateUser(ctx, user)

	if err != nil {
		return nil, errors.InternalError(ctx, "failed to create user", err)
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, id, req.Device)

	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		UserId:       id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.InvalidArgument("username and password are required")
	}

	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)

	if err != nil {
		return nil, errors.InternalError(ctx, "failed to get user by username", err)
	}

	if user == nil {
		return nil, errors.Unauthenticated("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.Unauthenticated("invalid credentials")
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user.ID, req.Device)

	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.RefreshToken == "" {
		return nil, errors.InvalidArgument("refresh token is required")
	}

	tokenHash := hashToken(req.RefreshToken)

	refreshToken, err := s.refreshTokenRepo.GetRefreshToken(ctx, tokenHash)

	if err != nil {
		return nil, errors.InternalError(ctx, "failed to get refresh token", err)
	}

	if refreshToken == nil {
		return nil, errors.Unauthenticated("invalid refresh token")
	}

	if err := s.refreshTokenRepo.RevokeRefreshToken(ctx, tokenHash, refreshToken.UserID, refreshToken.ExpiresAt); err != nil {
		return nil, errors.InternalError(ctx, "failed to revoke refresh token", err)
	}

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, errors.InvalidArgument("refresh token is required")
	}

	tokenHash := hashToken(req.RefreshToken)

	isBlacklisted, err := s.refreshTokenRepo.IsBlacklisted(ctx, tokenHash)

	if err != nil {
		return nil, errors.InternalError(ctx, "failed to check if refresh token is blacklisted", err)
	}

	if isBlacklisted {
		return nil, errors.Unauthenticated("refresh token has been revoked")
	}

	refreshToken, err := s.refreshTokenRepo.GetRefreshToken(ctx, tokenHash)

	if err != nil {
		return nil, errors.InternalError(ctx, "failed to get refresh token", err)
	}

	if refreshToken == nil {
		return nil, errors.Unauthenticated("invalid refresh token")
	}

	accessToken, newRefreshToken, err := s.generateTokens(ctx, refreshToken.UserID, refreshToken.Device)

	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.RevokeRefreshToken(ctx, tokenHash, refreshToken.UserID, refreshToken.ExpiresAt); err != nil {
		return nil, errors.InternalError(ctx, "failed to revoke refresh token", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if req.AccessToken == "" {
		return nil, errors.InvalidArgument("access token is required")
	}

	claims, err := jwt.ValidateToken(req.AccessToken, s.config.JWTSecret)

	if err != nil {
		return nil, errors.Unauthenticated("invalid access token")
	}

	return &pb.VerifyTokenResponse{
		UserId:    claims.UserID,
		ExpiresAt: timestamppb.New(claims.ExpiresAt),
	}, nil
}
