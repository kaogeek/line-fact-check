package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
)

func main() {
	container, cleanup, err := di.InitializeContainer()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	defer func() {
		slog.InfoContext(ctx, "server cleaning up")
		cleanup()
		slog.InfoContext(ctx, "server cleanup completed, exiting...")
	}()
	go func() {
		slog.InfoContext(ctx, "server starting", "addr", container.Conf.HTTP.ListenAddr)
		err := container.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.InfoContext(ctx, "server shutting down...")

	timeout := time.Second * 30
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := container.Server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "server forced to shutdown", "timeout", timeout.String(), "error", err)
	}
}
