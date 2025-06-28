package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kaogeek/line-fact-check/pillars"

	"github.com/kaogeek/line-fact-check/factcheck"
)

func main() {
	const (
		name = "factcheck-api"
		addr = ":8080"
	)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/", pillars.HandlerEcho(name))
	r.Handle("/health", pillars.HandlerOk(name))

	r.Get("/topic", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(map[string]any{
			"topic": factcheck.Topic{
				ID:        "some-topic-id",
				Name:      "TOPIC-FOOBAR",
				Status:    factcheck.StatusTopicResolved,
				Result:    "this is a fake nees",
				CreatedAt: time.Now(),
				UpdatedAt: nil,
			},
			"client": map[string]any{
				"uri":        r.RequestURI,
				"url":        r.URL.String(),
				"user-agent": r.UserAgent(),
				"host":       r.Host,
			},
		})
		if err != nil {
			contentTypeText(w.Header())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		contentTypeJSON(w.Header())
		w.WriteHeader(200)
		w.Write(j)
	}))

	srv := http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 2,
	}
	slog.Info("listening",
		"addr", addr,
	)
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func contentTypeJSON(h http.Header) {
	h.Add("Content-Type", "application/json; charset=utf-8")
}

func contentTypeText(h http.Header) {
	h.Add("Content-Type", "text/plain")
}
