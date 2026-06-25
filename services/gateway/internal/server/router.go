package server

import (
	"gateway/internal/config"
	"gateway/internal/handler"
	"gateway/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	pb "github.com/hardsmile98/messager/sdk/auth/v1"
)

func newRouter(authClient pb.AuthServiceClient, conf *config.Config) http.Handler {
	authHandler := handler.NewAuthHandler(authClient)

	r := chi.NewRouter()
	r.Use(middleware.CorsMiddleware)

	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login", authHandler.Login)
	r.Post("/api/v1/auth/logout", authHandler.Logout)
	r.Post("/api/v1/auth/refresh-token", authHandler.RefreshToken)

	return r
}
