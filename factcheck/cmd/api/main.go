package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
)

func main() {
	container, cleanup, err := di.InitializeContainer()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer func() {
		slog.InfoContext(ctx, "[main] server cleaning up")
		cleanup()
		slog.InfoContext(ctx, "[main] server cleanup completed, exiting...")
	}()

	quit := make(chan os.Signal, 1) // Buffered so it won't block on 2x Ctrl-C
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		slog.InfoContext(ctx, "[main] server starting", "config_http", container.Config.HTTP)
		err := container.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "[main] server error", "error", err)
		}
	}()

	<-quit
	slog.InfoContext(ctx, "[main] server shutting down...")
}
