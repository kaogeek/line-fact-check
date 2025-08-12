package handler

import (
	"net/http"
	"strings"

	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

// ListMessageGroupDynamic implements Handler.
func (h *handler) ListMessageGroupDynamic(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := limitOffSet(r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	opts := toMessageGroupOptions(r)
	topics, err := h.groups.ListDynamic(r.Context(), limit, offset, opts...)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(r.Context(), w, http.StatusOK, topics)
}

func toMessageGroupOptions(r *http.Request) []repo.OptionMessageGroup {
	query := r.URL.Query().Get
	text, idIn, idNotIn := query("like_message_text"), query("in_id"), query("not_in_id")

	var opts []repo.OptionMessageGroup

	if text != "" {
		opts = append(opts, repo.MessageGroupLikeMessageText(text))
	}

	if idIn != "" {
		parts := strings.Split(idIn, ",")
		opts = append(opts, repo.MessageGroupIDIn(parts))
	}

	if idNotIn != "" {
		parts := strings.Split(idIn, ",")
		opts = append(opts, repo.MessageGroupIDNotIn(parts))
	}

	return opts
}
