package handlers

import (
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type handler struct {
	topics repo.RepositoryTopic
}

func (h *handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	create(w, r, h.topics)
}

func (h *handler) ListTopics(w http.ResponseWriter, r *http.Request) {
	list(w, r, h.topics)
}
