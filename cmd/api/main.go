package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/mell62/journeyveil/internal/api"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	address := os.Getenv("HTTP_ADDRESS")
	if address == "" {
		address = ":8080"
	}

	server, err := api.NewServer(address, api.NewHandler(logger))
	if err != nil {
		logger.Error("create API server", "error", err)
		os.Exit(1)
	}

	logger.Info("API server listening", "address", server.Addr())
	if err := server.Run(ctx); err != nil {
		logger.Error("run API server", "error", err)
		os.Exit(1)
	}
}
