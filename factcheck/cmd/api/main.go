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
	const (
		name = "factcheck-api"
		addr = ":8080"
	)
	srv, cleanup, err := di.InitializeServer()
	if err != nil {
		panic(err)
	}
	defer func() {
		slog.Info("server cleaning up")
		cleanup()
		slog.Info("server cleanup completed, exiting...")
	}()

	go func() {
		slog.Info("server starting", "addr", addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("server shutting down...")

	timeout := time.Second * 30
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "timeout", timeout.String(), "error", err)
	}
}
