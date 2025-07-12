package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func (h *handler) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), h.messages.GetByID)
}

func (h *handler) DeleteMessageByID(w http.ResponseWriter, r *http.Request) {
	deleteByID[factcheck.Message](w, r, h.messages.Delete)
}

func (h *handler) ListMessagesByTopicID(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), h.messages.ListByTopic)
}

func (h *handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	create(
		w, r, h.messages.Create,
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
				CreatedAt: utils.TimeNow(),
				UpdatedAt: nil,
			}
		}),
	)
}
