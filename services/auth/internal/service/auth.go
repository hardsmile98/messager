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
	"google.golang.org/protobuf/types/known/timestamppb"
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

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateRefreshToken() (string, string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)

	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to generate refresh token: %v", err)
	}

	token := base64.RawURLEncoding.EncodeToString(b)

	tokenHash := hashToken(token)

	return token, tokenHash, nil
}

func (s *AuthService) generateTokens(ctx context.Context, userID string) (string, string, error) {
	accessToken, err := jwt.GenerateToken(userID, s.config.JWTSecret, s.config.AccessTokenTTL)

	if err != nil {
		return "", "", status.Errorf(codes.Internal, "failed to generate access token: %v", err)
	}

	refresh, refreshHash, err := generateRefreshToken()

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

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	user, err := s.userRepo.GetUserByUsername(ctx, req.Username)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user by username: %v", err)
	}

	if user == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	accessToken, refreshToken, err := s.generateTokens(ctx, user.ID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	return &pb.LoginResponse{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user ID is required")
	}

	err := s.refreshTokenRepo.RevokeRefreshTokenByUserID(ctx, req.UserId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete refresh token: %v", err)
	}

	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	tokenHash := hashToken(req.RefreshToken)

	isBlacklisted, err := s.refreshTokenRepo.IsBlacklisted(ctx, tokenHash)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check if refresh token is blacklisted: %v", err)
	}

	if isBlacklisted {
		return nil, status.Error(codes.Unauthenticated, "refresh token has been revoked")
	}

	refreshToken, err := s.refreshTokenRepo.GetRefreshToken(ctx, tokenHash)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get refresh token: %v", err)
	}

	if refreshToken == nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	if err := s.refreshTokenRepo.RevokeRefreshToken(ctx, tokenHash, refreshToken.ExpiresAt); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to revoke refresh token: %v", err)
	}

	accessToken, newRefreshToken, err := s.generateTokens(ctx, refreshToken.UserID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate tokens: %v", err)
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if req.AccessToken == "" {
		return nil, status.Error(codes.InvalidArgument, "access token is required")
	}

	claims, err := jwt.ValidateToken(req.AccessToken, s.config.JWTSecret)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate access token: %v", err)
	}

	return &pb.VerifyTokenResponse{
		UserId:    claims.UserID,
		ExpiresAt: timestamppb.New(claims.ExpiresAt),
	}, nil
}
