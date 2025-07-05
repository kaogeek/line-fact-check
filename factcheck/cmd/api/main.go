package main

import (
	"log/slog"
	"net/http"
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

	container, err := di.InitializeContainer()
	if err != nil {
		panic(err)
	}

	topics := chi.NewMux()
	topics.Get("/", container.Handler.ListTopics)
	topics.Post("/", container.Handler.CreateTopic)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/", pillars.HandlerEcho(name))
	r.Handle("/health", pillars.HandlerOk(name))
	r.Mount("/topics", topics) // TODO: get by id

	srv := http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 2,
	}

	slog.Info("listening", "addr", addr)
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
