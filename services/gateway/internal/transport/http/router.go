package httptransport

import (
	"net/http"

	"gateway/internal/transport/http/handler"
	"gateway/internal/transport/http/middleware"

	"github.com/go-chi/chi/v5"
	pb "github.com/hardsmile98/messager/sdk/auth/v1"
)

func NewRouter(authClient pb.AuthServiceClient) http.Handler {
	authHandler := handler.NewAuth(authClient)

	r := chi.NewRouter()
	r.Use(middleware.CORS)

	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login", authHandler.Login)
	r.Post("/api/v1/auth/logout", authHandler.Logout)
	r.Post("/api/v1/auth/refresh-token", authHandler.RefreshToken)

	return r
}
