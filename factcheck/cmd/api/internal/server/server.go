// Package server defines top-level http server for factcheck-api
package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kaogeek/line-fact-check/pillars"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func New(conf config.Config, h handler.Handler) *http.Server {
	messages, messageGroups, topics := chi.NewMux(), chi.NewMux(), chi.NewMux()

	messages.Post("/", h.SubmitMessage)
	messages.Put("/{id}/assign-message-group", h.AssignMessageGroup)
	messages.Delete("/", h.DeleteMessageByID)

	messageGroups.Put("/{id}/assign-topic", h.AssignMessageGroup)

	topics.Post("/", h.CreateTopic)
	topics.Get("/all", h.ListAllTopics)
	topics.Get("/", h.ListTopicsHome)
	topics.Get("/count", h.CountTopicsHome)
	topics.Get("/{id}", h.GetTopicByID)
	topics.Get("/{id}/answer", h.GetAnswer)
	topics.Get("/{id}/answers", h.ListAnswers)
	topics.Get("/{id}/messages", h.ListTopicMessages)
	topics.Get("/{id}/message-group", h.ListTopicMessageGroups)
	topics.Put("/{id}/status", h.UpdateTopicStatus)
	topics.Put("/{id}/description", h.UpdateTopicDescription)
	topics.Put("/{id}/name", h.UpdateTopicName)
	topics.Delete("/{id}", h.DeleteTopicByID)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Handle("/", pillars.HandlerEcho(conf.AppName))
	r.Handle("/health", pillars.HandlerOk(conf.AppName))
	r.Mount("/topics", topics)
	r.Mount("/messages", messages)
	r.Mount("/message-groups", messageGroups)

	return &http.Server{
		Addr:         utils.DefaultIfZero(conf.HTTP.ListenAddr, ":8080"),
		Handler:      r,
		ReadTimeout:  utils.DefaultIfZero(time.Duration(conf.HTTP.TimeoutMsRead)*time.Millisecond, time.Second),
		WriteTimeout: utils.DefaultIfZero(time.Duration(conf.HTTP.TimeoutMsWrite)*time.Millisecond, time.Second),
	}
}
