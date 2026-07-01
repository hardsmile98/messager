package httptransport

import (
	"net/http"

	"gateway/internal/config"
	"gateway/internal/transport/http/handler"
	"gateway/internal/transport/http/middleware"

	"github.com/go-chi/chi/v5"
	pbAuth "github.com/hardsmile98/messager/sdk/auth/v1"
	pbChat "github.com/hardsmile98/messager/sdk/chat/v1"
)

func NewRouter(
	authClient pbAuth.AuthServiceClient,
	chatClient pbChat.ChatServiceClient,
	cfg *config.Config,
) http.Handler {
	authHandler := handler.NewAuthHandler(authClient, cfg)
	chatHandler := handler.NewChatHandler(chatClient)

	r := chi.NewRouter()
	r.Use(middleware.CORS(cfg.CORSAllowedOrigins))

	r.Post("/api/v1/auth/register", authHandler.Register)
	r.Post("/api/v1/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authClient))
		r.Post("/api/v1/auth/logout", authHandler.Logout)
		r.Post("/api/v1/auth/refresh-token", authHandler.RefreshToken)

		r.Post("/api/v1/chats/private", chatHandler.CreatePrivateChat)
		r.Get("/api/v1/chats", chatHandler.GetUserChats)
		r.Get("/api/v1/chats/{chat_id}", chatHandler.GetChatInfo)
	})

	return r
}
