// Package handlers provides HTTP server handlers
package handlers

import (
	"fmt"
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck"
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

func (h *handler) ListTopics(w http.ResponseWriter, r *http.Request)   { list(w, r, h.topics) }
func (h *handler) GetTopicByID(w http.ResponseWriter, r *http.Request) { getByID(w, r, h.topics) }

func (h *handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	topic, err := decode[factcheck.Topic](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	if topic.Status != "" {
		errBadRequest(w, fmt.Sprintf("unexpected topic status: '%s'", topic.Status))
		return
	}
	topic.Status = factcheck.StatusTopicPending
	created, err := h.topics.Create(r.Context(), topic)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, created, http.StatusCreated)
}
