// Package handlers provides HTTP server handlers
package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Handler interface {
	CreateTopic(http.ResponseWriter, *http.Request)
	ListTopics(http.ResponseWriter, *http.Request)
	GetTopicByID(http.ResponseWriter, *http.Request)
	DeleteTopicByID(http.ResponseWriter, *http.Request)

	CreateMessage(http.ResponseWriter, *http.Request)
	ListMessagesByTopicID(http.ResponseWriter, *http.Request)
	DeleteMessageByID(http.ResponseWriter, *http.Request)
}

type handler struct {
	topics   repo.RepositoryTopic
	messages repo.RepositoryMessage
}

func New(repo repo.Repository) Handler {
	return &handler{
		topics:   repo.Topic,
		messages: repo.Message,
	}
}

func (h *handler) ListTopics(w http.ResponseWriter, r *http.Request) {
	list(w, r, h.topics)
}

func (h *handler) GetTopicByID(w http.ResponseWriter, r *http.Request) {
	getByID(w, r, h.topics)
}

func (h *handler) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	getByID(w, r, h.messages)
}

func (h *handler) DeleteTopicByID(w http.ResponseWriter, r *http.Request) {
	deleteByID[factcheck.Topic](w, r, h.topics)
}

func (h *handler) DeleteMessageByID(w http.ResponseWriter, r *http.Request) {
	deleteByID[factcheck.Message](w, r, h.messages)
}

func (h *handler) ListMessagesByTopicID(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), h.messages.ListByTopic)
}

func (h *handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	create(
		w, r, h.topics,
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
				Status:       factcheck.StatusTopicPending,
				Result:       topic.Result,
				ResultStatus: topic.ResultStatus,
				CreatedAt:    time.Now(),
				UpdatedAt:    nil,
			}
		}),
	)
}

func (h *handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	create(
		w, r, h.messages,
		createCheck(func(_ context.Context, m factcheck.Message) error {
			if m.Text == "" {
				return errors.New("empty text")
			}
			if !m.Type.IsValid() {
				return fmt.Errorf("invalid type '%s'", m.Type)
			}
			return nil
		}),
		createModify(func(_ context.Context, m factcheck.Message) factcheck.Message {
			return factcheck.Message{
				ID:        uuid.New().String(),
				TopicID:   m.TopicID,
				Text:      m.Text,
				Type:      m.Type,
				CreatedAt: time.Now(),
				UpdatedAt: nil,
			}
		}),
	)
}
