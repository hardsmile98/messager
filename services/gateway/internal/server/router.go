package server

import (
	"gateway/internal/config"
	"gateway/internal/handler"
	"gateway/internal/middleware"
	"net/http"

	"github.com/go-chi/chi"
	pb "github.com/hardsmile98/messager/sdk/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRouter(conf *config.Config) (http.Handler, error) {
	authConnection, err := grpc.NewClient(
		conf.AuthGRPCURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}

	authClient := pb.NewAuthServiceClient(authConnection)

	authHandler := handler.NewAuthHandler(authClient)

	r := chi.NewRouter()
	r.Use(middleware.CorsMiddleware)

	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login", authHandler.Login)
	r.Post("/api/v1/auth/logout", authHandler.Logout)
	r.Post("/api/v1/auth/refresh-token", authHandler.RefreshToken)

	return r, nil
}
