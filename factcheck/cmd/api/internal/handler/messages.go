package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kaogeek/line-fact-check/factcheck"
)

// TODO: use middleware to parse the metadata and save it to req context
// when doing auth
func (h *handler) getUserInfo(_ *http.Request) (factcheck.UserInfo, error) {
	return factcheck.UserInfo{
		UserType: factcheck.TypeUserMessageAdmin,
		UserID:   "user-mock-getuserinfo",
	}, nil
}

func (h *handler) GetMessageByID(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), func(ctx context.Context, id string) (factcheck.MessageV2, error) {
		return h.messagesv2.GetByID(ctx, id)
	})
}

func (h *handler) DeleteMessageByID(w http.ResponseWriter, r *http.Request) {
	deleteByID[factcheck.MessageV2](w, r, func(ctx context.Context, s string) error {
		return h.messagesv2.Delete(ctx, s)
	})
}

func (h *handler) ListMessagesInGroup(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, chi.URLParam(r, "group_id"), func(ctx context.Context, s string) ([]factcheck.MessageV2, error) {
		return h.messagesv2.ListByGroup(ctx, s)
	})
}

func (h *handler) AssignMessageGroup(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)
	if id == "" {
		errBadRequest(w, "missing message_id")
		return
	}
	body, err := decode[struct {
		GroupID string `json:"group_id"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	if body.GroupID == "" {
		errBadRequest(w, err.Error())
		return
	}
	msg, err := h.messagesv2.AssignGroup(r.Context(), id, body.GroupID)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(r.Context(), w, http.StatusOK, msg)
}

func (h *handler) SubmitMessage(w http.ResponseWriter, r *http.Request) {
	body, err := decode[struct {
		Text    string `json:"text"`
		TopicID string `json:"topic_id"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	userInfo, err := h.getUserInfo(r)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	topic, msg, group, err := h.service.Submit(
		r.Context(),
		userInfo,
		body.Text,
		body.TopicID,
	)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(r.Context(), w, http.StatusCreated, map[string]any{
		"topic":   topic,
		"group":   group,
		"message": msg,
	})
}
