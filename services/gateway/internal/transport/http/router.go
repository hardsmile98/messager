package httptransport

import (
	"net/http"

	"gateway/internal/config"
	"gateway/internal/transport/http/handler"
	"gateway/internal/transport/http/middleware"

	"github.com/go-chi/chi/v5"
	pb "github.com/hardsmile98/messager/sdk/auth/v1"
)

func NewRouter(authClient pb.AuthServiceClient, cfg *config.Config) http.Handler {
	authHandler := handler.NewAuth(authClient, cfg)

	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))

	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authClient))
		r.Post("/api/v1/auth/logout", authHandler.Logout)
		r.Post("/api/v1/auth/refresh-token", authHandler.RefreshToken)
	})

	return r
}
