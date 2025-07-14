package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

type userInfo struct {
	UserType factcheck.TypeUserMessage `json:"user_type"`
	UserID   string                    `json:"user_id"`
}

// TODO: use middleware to parse the metadata and save it to req context
// when doing auth
func (h *handler) getUserInfo(_ *http.Request) (userInfo, error) {
	return userInfo{
		UserType: factcheck.TypeUserMessageAdmin,
		UserID:   "user-mock-getuserinfo",
	}, nil
}

func (h *handler) NewUserMessage(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Text    string `json:"text"`
		TopicID string `json:"topic_id"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	meta, err := h.getUserInfo(r)
	if err != nil {
		slog.Error("error getting user info",
			"err", err,
			"message", body.Text,
		)
		errInternalError(w, "error getting user info")
		return
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		errInternalError(w, fmt.Sprintf("error encoding metadata: %s", err.Error()))
		return
	}

	statusMessage := factcheck.StatusMessageSubmitted
	if body.TopicID != "" {
		statusMessage = factcheck.StatusMessageTopicSubmitted
	}
	now := utils.TimeNow()
	userMessage := factcheck.UserMessage{
		ID:        uuid.NewString(),
		Type:      meta.UserType,
		RepliedAt: nil,
		Metadata:  metaJSON,
		CreatedAt: now,
		UpdatedAt: nil,
	}
	message := factcheck.Message{
		ID:            uuid.NewString(),
		UserMessageID: userMessage.ID,
		Type:          factcheck.TypeMessageText,
		Status:        statusMessage,
		TopicID:       body.TopicID,
		Text:          body.Text,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	tx, err := h.repository.Begin(r.Context())
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	// As per doc, they defer rollback
	// https://docs.sqlc.dev/en/stable/howto/transactions.html
	defer tx.Rollback(r.Context())

	createdUserMessage, err := h.userMessages.Create(r.Context(), userMessage)
	if err != nil {
		errInternalError(w, fmt.Sprintf("error creating user_messages: %s", err.Error()))
		return
	}
	createdMessage, err := h.messages.Create(r.Context(), message)
	if err != nil {
		errInternalError(w, fmt.Sprintf("error creating messages: %s", err.Error()))
		return
	}

	tx.Commit(r.Context())

	sendJSON(w, map[string]any{
		"user_message": createdUserMessage,
		"message":      createdMessage,
	}, http.StatusCreated)
}
