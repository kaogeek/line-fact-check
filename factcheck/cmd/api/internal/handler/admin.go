package handler

import (
	"context"
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck"
)

func (h *handler) adminMiddleware(w http.ResponseWriter, r *http.Request) {
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
		errBadRequest(w, "missing group_id")
		return
	}
	msg, err := h.messagesv2.AssignGroup(r.Context(), id, body.GroupID)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(r.Context(), w, http.StatusOK, msg)
}

func (h *handler) AssignGroupTopic(w http.ResponseWriter, r *http.Request) {
	id := paramID(r)
	if id == "" {
		errBadRequest(w, "missing message_group_id")
		return
	}
	body, err := decode[struct {
		TopicID string `json:"topic_id"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	if body.TopicID == "" {
		errBadRequest(w, "missing topic_id")
		return
	}
	group, err := h.groups.AssignTopic(r.Context(), id, body.TopicID)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(r.Context(), w, http.StatusOK, group)
}

func (h *handler) PostAnswer(w http.ResponseWriter, r *http.Request) {
	data, err := decode[struct {
		Text string `json:"text"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	user, err := h.getUserInfo(r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	answer, _, _, err := h.service.ResolveTopic(r.Context(), user, paramID(r), data.Text)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(r.Context(), w, http.StatusOK, answer)
}

func (h *handler) DeleteTopicByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserInfo(r)
	if err != nil {
		errBadRequest(w, "error getting user info from request")
		return
	}
	if user.UserID == "" { //nolint
		errAuth(w)
		return
	}
	if user.UserType != factcheck.TypeUserMessageAdmin {
		errAuth(w)
		return
	}
	deleteByID[factcheck.Topic](w, r, func(ctx context.Context, s string) error {
		return h.topics.Delete(ctx, s)
	})
}

func (h *handler) DeleteAnswerByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserInfo(r)
	if err != nil {
		errBadRequest(w, "error getting user info from request")
		return
	}
	if user.UserID == "" { //nolint
		errAuth(w)
		return
	}
	if user.UserType != factcheck.TypeUserMessageAdmin {
		errAuth(w)
		return
	}
	deleteByID[factcheck.Answer](w, r, func(ctx context.Context, id string) error {
		return h.answers.Delete(ctx, id)
	})
}

func (h *handler) DeleteGroupByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserInfo(r)
	if err != nil {
		errBadRequest(w, "error getting user info from request")
		return
	}
	if user.UserID == "" { //nolint
		errAuth(w)
		return
	}
	if user.UserType != factcheck.TypeUserMessageAdmin {
		errAuth(w)
		return
	}
	deleteByID[factcheck.MessageGroup](w, r, func(ctx context.Context, id string) error {
		return h.groups.Delete(ctx, id)
	})
}
