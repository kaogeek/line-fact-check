package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/kaogeek/line-fact-check/pillars"
)

func main() {
	const (
		name = "foo-api"
		addr = ":8080"
	)

	mux := http.NewServeMux()
	mux.Handle("/", pillars.HandlerEcho(name))
	mux.Handle("/health", pillars.HandlerOk(name))

	srv := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 2,
	}

	slog.InfoContext(context.Background(), "listening",
		"addr", addr,
	)
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
