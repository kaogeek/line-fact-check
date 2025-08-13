package handler

import (
	"net/http"
	"strings"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
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

func (h *handler) CountMessageGroupsDynamic(w http.ResponseWriter, r *http.Request) {
	opts := toMessageGroupOptions(r)
	counts, err := h.groups.CountDynamic(r.Context(), opts...)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	result := make(map[string]int64)
	for k, v := range counts {
		result[string(k)] = v
		result["total"] += v
	}
	sendJSON(r.Context(), w, http.StatusOK, result)
}

func toMessageGroupOptions(r *http.Request) []repo.OptionMessageGroup {
	query := r.URL.Query().Get
	text, idIn, idNotIn, statusesIn := query("like_message_text"), query("in_id"), query("not_in_id"), query("statuses_in")

	var opts []repo.OptionMessageGroup

	if text != "" {
		opts = append(opts, repo.MessageGroupLikeMessageText(text))
	}

	if idIn != "" {
		parts := strings.Split(idIn, ",")
		opts = append(opts, repo.MessageGroupIDIn(parts))
	}

	if idNotIn != "" {
		parts := strings.Split(idNotIn, ",")
		opts = append(opts, repo.MessageGroupIDNotIn(parts))
	}

	if statusesIn != "" {
		parts := strings.Split(statusesIn, ",")
		if len(parts) != 0 {
			statuses := utils.MapNoError(parts, func(s string) factcheck.StatusMGroup {
				return factcheck.StatusMGroup(s)
			})
			opts = append(opts, repo.MessageGroupStatusesIn(statuses))
		}
	}

	return opts
}
