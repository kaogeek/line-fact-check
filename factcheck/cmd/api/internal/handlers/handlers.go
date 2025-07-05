// Package handlers provides HTTP server handlers
package handlers

import (
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Handler interface {
	CreateTopic(http.ResponseWriter, *http.Request)
	ListTopics(http.ResponseWriter, *http.Request)
	GetTopicByID(http.ResponseWriter, *http.Request)
}

type handler struct {
	topics repo.RepositoryTopic
}

func New(repo *repo.Repository) Handler {
	return &handler{
		topics: repo.Topic,
	}
}

func (h *handler) CreateTopic(w http.ResponseWriter, r *http.Request)  { create(w, r, h.topics) }
func (h *handler) ListTopics(w http.ResponseWriter, r *http.Request)   { list(w, r, h.topics) }
func (h *handler) GetTopicByID(w http.ResponseWriter, r *http.Request) { getByID(w, r, h.topics) }
