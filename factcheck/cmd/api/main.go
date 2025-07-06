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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kaogeek/line-fact-check/pillars"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
)

func main() {
	const (
		name = "factcheck-api"
		addr = ":8080"
	)

	container, cleanup, err := di.InitializeContainer()
	if err != nil {
		panic(err)
	}
	defer func() {
		slog.Info("server cleaning up")
		cleanup()
		slog.Info("server cleanup completed, exiting...")
	}()

	topics, messages := chi.NewMux(), chi.NewMux()
	topics.Get("/", container.Handler.ListTopics)
	topics.Get("/{id}", container.Handler.GetTopicByID)
	topics.Post("/", container.Handler.CreateTopic)
	topics.Delete("/{id}", container.Handler.DeleteTopicByID)
	messages.Get("/{id}", container.Handler.ListMessagesByTopicID)
	messages.Post("/", container.Handler.CreateMessage)
	messages.Delete("/", container.Handler.DeleteMessageByID)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Handle("/", pillars.HandlerEcho(name))
	r.Handle("/health", pillars.HandlerOk(name))
	r.Mount("/topics", topics)
	r.Mount("/messages", messages)

	srv := http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 2,
	}

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
