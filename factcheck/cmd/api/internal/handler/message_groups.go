package handler

import "net/http"

func (h *handler) AssignTopic(w http.ResponseWriter, r *http.Request) {
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
