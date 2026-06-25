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
	UserID    string
	ExpiresAt time.Time
}

func GenerateToken(userID, secret string, ttlMinutes int) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(ttlMinutes) * time.Minute)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func ValidateToken(token, secret string) (VerifyToken, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
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

	if claims.UserID == "" {
		return VerifyToken{}, errors.New("invalid token claims")
	}

	return VerifyToken{
		UserID:    claims.UserID,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}
