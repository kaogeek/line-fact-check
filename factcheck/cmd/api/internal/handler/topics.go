// Package handlers provides HTTP server handlers
package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func (h *handler) ListTopics(w http.ResponseWriter, r *http.Request) {
	list(w, r, h.topics.List)
}

func (h *handler) GetTopicByID(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), h.topics.GetByID)
}

func (h *handler) DeleteTopicByID(w http.ResponseWriter, r *http.Request) {
	deleteByID[factcheck.Topic](w, r, h.topics.Delete)
}

func (h *handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	create(
		w, r, h.topics.Create,
		createCheck(func(_ context.Context, topic factcheck.Topic) error {
			if topic.ID != "" {
				return fmt.Errorf("unexpected topic id (expecting empty topic id): '%s'", topic.Status)
			}
			if topic.Status != "" {
				return fmt.Errorf("unexpected topic status (expecting empty topic status): '%s'", topic.Status)
			}
			return nil
		}),
		createModify(func(_ context.Context, topic factcheck.Topic) factcheck.Topic {
			return factcheck.Topic{
				ID:           uuid.New().String(),
				Name:         topic.Name,
				Description:  topic.Description,
				Status:       factcheck.StatusTopicPending,
				Result:       topic.Result,
				ResultStatus: factcheck.StatusTopicResultNone,
				CreatedAt:    utils.TimeNow(),
				UpdatedAt:    nil,
			}
		}),
	)
}

func (h *handler) UpdateTopicStatus(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Status string `json:"status"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	topic, err := h.topics.UpdateStatus(r.Context(), paramID(r), factcheck.StatusTopic(body.Status))
	if err != nil {
		if repo.IsNotFound(err) {
			var notFoundErr *repo.ErrNotFound
			if errors.As(err, &notFoundErr) {
				errNotFound(w, notFoundErr.Error())
			} else {
				errNotFound(w, fmt.Sprintf("topic not found: %s", paramID(r)))
			}
			return
		}
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, topic, http.StatusOK)
}

func (h *handler) UpdateTopicDescription(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Description string `json:"description"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	topic, err := h.topics.UpdateDescription(r.Context(), paramID(r), body.Description)
	if err != nil {
		if repo.IsNotFound(err) {
			var notFoundErr *repo.ErrNotFound
			if errors.As(err, &notFoundErr) {
				errNotFound(w, notFoundErr.Error())
			} else {
				errNotFound(w, fmt.Sprintf("topic not found: %s", paramID(r)))
			}
			return
		}
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, topic, http.StatusOK)
}

func (h *handler) UpdateTopicName(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Name string `json:"name"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	topic, err := h.topics.UpdateName(r.Context(), paramID(r), body.Name)
	if err != nil {
		if repo.IsNotFound(err) {
			var notFoundErr *repo.ErrNotFound
			if errors.As(err, &notFoundErr) {
				errNotFound(w, notFoundErr.Error())
			} else {
				errNotFound(w, fmt.Sprintf("topic not found: %s", paramID(r)))
			}
			return
		}
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, topic, http.StatusOK)
}
