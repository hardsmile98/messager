package middleware

import (
	"context"
	"net/http"
	"strings"

	"gateway/internal/transport/http/cookie"
	"gateway/internal/transport/http/response"

	pb "github.com/hardsmile98/messager/sdk/auth/v1"
)

type contextKey string

const (
	refreshTokenKey contextKey = "refresh_token"
	userIDKey       contextKey = "user_id"
)

func RequireAuth(client pb.AuthServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken := accessTokenFromRequest(r)

			if accessToken == "" {
				response.JSON(w, http.StatusUnauthorized, map[string]string{
					"error": "authentication required",
				})
				return
			}

			resp, err := client.VerifyToken(r.Context(), &pb.VerifyTokenRequest{
				AccessToken: accessToken,
			})

			if err != nil {
				response.JSON(w, http.StatusUnauthorized, map[string]string{
					"error": "authentication required",
				})
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), userIDKey, resp.UserId))

			refreshToken := cookie.RefreshToken(r)

			r = r.WithContext(context.WithValue(r.Context(), refreshTokenKey, refreshToken))

			next.ServeHTTP(w, r)
		})
	}
}

func accessTokenFromRequest(r *http.Request) string {
	if token := cookie.AccessToken(r); token != "" {
		return token
	}

	auth := r.Header.Get("Authorization")

	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	return ""
}

func RefreshTokenFromContext(r *http.Request) (string, bool) {
	token, ok := r.Context().Value(refreshTokenKey).(string)
	return token, ok && token != ""
}

func UserIDFromContext(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(userIDKey).(string)
	return userID, ok && userID != ""
}
