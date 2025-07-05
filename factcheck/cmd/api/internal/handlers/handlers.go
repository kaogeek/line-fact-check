package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type handler struct {
	topics repo.RepositoryTopic
}

func (h *handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	b, err := body(r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	var topic factcheck.Topic
	err = json.Unmarshal(b, &topic)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	created, err := h.topics.Create(r.Context(), topic)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	err = replyJson(w, created, http.StatusCreated)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
}

func (h *handler) ListTopics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	list, err := h.topics.List(ctx)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	err = replyJson(w, list, http.StatusOK)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
}
