package main

import (
	"gateway/internal/config"
	"gateway/internal/server"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	conf, err := config.LoadConfig()

	if err != nil {
		slog.Error("failed to load config", "error", err)
	}

	router, err := server.NewRouter(conf)

	if err != nil {
		slog.Error("failed to create router", "error", err)
	}

	if err := http.ListenAndServe(":"+conf.Port, router); err != nil {
		slog.Error("http server", "error", err)
	}
}
