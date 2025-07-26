// Package handler provides HTTP server handlers
package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func (h *handler) ListAllTopics(w http.ResponseWriter, r *http.Request) {
	list(w, r, func(ctx context.Context) ([]factcheck.Topic, error) {
		return h.topics.List(r.Context(), 0, 0)
	})
}

// TODO: register route
func (h *handler) ListTopics(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := limitOffSet(r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	topics, err := h.topics.List(r.Context(), limit, offset)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, topics, http.StatusOK)
}

func (h *handler) GetTopicByID(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), func(ctx context.Context, id string) (factcheck.Topic, error) {
		return h.topics.GetByID(ctx, id)
	})
}

func (h *handler) DeleteTopicByID(w http.ResponseWriter, r *http.Request) {
	deleteByID[factcheck.Topic](w, r, func(ctx context.Context, s string) error {
		return h.topics.Delete(ctx, s)
	})
}

func (h *handler) ListTopicsHome(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := limitOffSet(r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	opts := toTopicOptions(r)
	topics, err := h.topics.ListDynamicV2(r.Context(), limit, offset, opts...)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, topics, http.StatusOK)
}

func (h *handler) CountTopicsHome(w http.ResponseWriter, r *http.Request) {
	opts := toTopicOptions(r)
	counts, err := h.topics.CountByStatusDynamic(r.Context(), opts...)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	result := make(map[string]int64)
	for k, v := range counts {
		result[string(k)] = v
		result["total"] += v
	}
	sendJSON(w, result, http.StatusOK)
}

func (h *handler) CreateTopic(w http.ResponseWriter, r *http.Request) {
	create(
		w, r,
		func(ctx context.Context, topic factcheck.Topic) (factcheck.Topic, error) {
			return h.topics.Create(ctx, topic)
		},
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
				ID:          utils.NewID().String(),
				Name:        topic.Name,
				Description: topic.Description,
				Status:      factcheck.StatusTopicPending,
				Result:      topic.Result,
				CreatedAt:   utils.TimeNow(),
				UpdatedAt:   nil,
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
		handleNotFound(w, err, "topic", paramID(r))
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
		handleNotFound(w, err, "topic", paramID(r))
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
		handleNotFound(w, err, "topic", paramID(r))
		return
	}
	sendJSON(w, topic, http.StatusOK)
}

func (h *handler) ListTopicMessages(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), func(ctx context.Context, id string) ([]factcheck.MessageV2, error) {
		return h.messagesv2.ListByTopic(ctx, id)
	})
}

func (h *handler) ListTopicMessageGroups(w http.ResponseWriter, r *http.Request) {
	getBy(w, r, paramID(r), func(ctx context.Context, s string) ([]factcheck.MessageGroup, error) {
		return h.groups.ListByTopic(ctx, s)
	})
}

func (h *handler) PostAnswer(w http.ResponseWriter, r *http.Request) {
	data, err := decode[struct {
		Text string `json:"text"`
	}](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	h.answers.Create(r.Context(), factcheck.Answer{
		ID:        utils.NewID().String(),
		TopicID:   paramID(r),
		Text:      data.Text,
		CreatedAt: utils.TimeNow(),
	})
}

func toTopicOptions(r *http.Request) []repo.OptionTopic {
	query := r.URL.Query().Get
	id, text, statuses := query("like_id"), query("like_message_text"), query("in_statuses")
	var opts []repo.OptionTopic
	if statuses != "" {
		parts := strings.Split(statuses, ",")
		if len(parts) != 0 {
			statuses := utils.MapNoError(parts, func(s string) factcheck.StatusTopic {
				return factcheck.StatusTopic(s)
			})
			opts = append(opts, repo.TopicInStatuses(statuses))
		}
	}
	if id != "" {
		opts = append(opts, repo.TopicLikeID(id))
	}
	if text != "" {
		opts = append(opts, repo.TopicLikeMessageText(text))
	}
	return opts
}
