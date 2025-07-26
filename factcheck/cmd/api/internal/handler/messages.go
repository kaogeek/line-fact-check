package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func (h *handler) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), func(ctx context.Context, id string) (factcheck.Message, error) {
		return h.messages.GetByID(ctx, id)
	})
}

func (h *handler) DeleteMessageByID(w http.ResponseWriter, r *http.Request) {
	deleteByID[factcheck.Message](w, r, func(ctx context.Context, s string) error {
		return h.messages.Delete(ctx, s)
	})
}

func (h *handler) ListMessagesInGroup(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, chi.URLParam(r, "group_id"), func(ctx context.Context, s string) ([]factcheck.MessageV2, error) {
		return h.messagesv2.ListByGroup(ctx, s)
	})
}

func (h *handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	create(
		w, r, func(ctx context.Context, m factcheck.Message) (factcheck.Message, error) {
			return h.messages.Create(ctx, m)
		},
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
				ID:        utils.NewID().String(),
				TopicID:   m.TopicID,
				Text:      m.Text,
				Type:      m.Type,
				CreatedAt: utils.TimeNow(),
				UpdatedAt: nil,
			}
		}),
	)
}
