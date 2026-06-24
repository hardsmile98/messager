package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type VerifyToken struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func GenerateToken(userID string, secret string, ttlMinutes int) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttlMinutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateToken(token string, secret string) (VerifyToken, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return "", errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return VerifyToken{}, err
	}

	claims, ok := parsedToken.Claims.(*Claims)

	if !ok || !parsedToken.Valid {
		return VerifyToken{}, errors.New("invalid token")
	}

	return VerifyToken{
		UserID:    claims.UserID,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
