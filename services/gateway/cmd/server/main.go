package main

import (
	"gateway/internal/config"
	"gateway/internal/server"
	"log/slog"
	"os"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	conf, err := config.LoadConfig()

	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := server.RunHTTPServer(conf); err != nil {
		slog.Error("failed to run http server", "error", err)
		os.Exit(1)
	}
}
