package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/openapi"
)

func New(r repo.Repository) openapi.StrictServerInterface {
	return &handler{
		topics:   r.Topic,
		messages: r.Message,
	}
}

type handler struct {
	topics   repo.RepositoryTopic
	messages repo.RepositoryMessage
}

func (h *handler) CreateTopic(ctx context.Context, request openapi.CreateTopicRequestObject) (openapi.CreateTopicResponseObject, error) {
	id := uuid.New()
	topic := factcheck.Topic{
		ID:        id.String(),
		Name:      request.Body.Name,
		Status:    factcheck.StatusTopicPending,
		CreatedAt: time.Now(),
	}
	created, err := h.topics.Create(ctx, topic)
	if err != nil {
		return openapi.CreateTopic500TextResponse(err.Error()), nil
	}
	if created.ID != id.String() {
		return openapi.CreateTopic500TextResponse(fmt.Sprintf("topic id mismatch: %s != %s", created.ID, id.String())), nil
	}
	return openapi.CreateTopic201JSONResponse(openapi.Topic{
		CreatedAt: created.CreatedAt,
		Id:        id,
		Name:      created.Name,
		Status:    openapi.TopicStatus(created.Status),
		UpdatedAt: created.UpdatedAt,
	}), nil
}

func (h *handler) ListTopics(ctx context.Context, request openapi.ListTopicsRequestObject) (openapi.ListTopicsResponseObject, error) {
	list, err := h.topics.List(ctx)
	if err != nil {
		return openapi.ListTopics500TextResponse(err.Error()), nil
	}
	topics := make([]openapi.Topic, len(list))
	for i := range list {
		t := &list[i]
		id, err := uuid.Parse(t.ID)
		if err != nil {
			return openapi.ListTopics500TextResponse(err.Error()), nil
		}
		topics[i] = openapi.Topic{
			CreatedAt: t.CreatedAt,
			Id:        id,
			Name:      t.Name,
			Status:    openapi.TopicStatus(t.Status),
			UpdatedAt: t.UpdatedAt,
		}
	}
	return openapi.ListTopics200JSONResponse(topics), nil
}

func (h *handler) GetTopicByID(ctx context.Context, request openapi.GetTopicByIDRequestObject) (openapi.GetTopicByIDResponseObject, error) {
	topic, err := h.topics.GetByID(ctx, request.Id.String())
	if err != nil {
		return openapi.GetTopicByID500TextResponse(err.Error()), nil
	}
	id, err := uuid.Parse(topic.ID)
	if err != nil {
		return openapi.GetTopicByID500TextResponse(err.Error()), nil
	}
	return openapi.GetTopicByID200JSONResponse(openapi.Topic{
		Id:        id,
		CreatedAt: topic.CreatedAt,
		Name:      topic.Name,
		Status:    openapi.TopicStatus(topic.Status),
		UpdatedAt: topic.UpdatedAt,
	}), nil
}

func (h *handler) DeleteTopicByID(ctx context.Context, request openapi.DeleteTopicByIDRequestObject) (openapi.DeleteTopicByIDResponseObject, error) {
	err := h.topics.Delete(ctx, request.Id.String())
	if err != nil {
		return openapi.DeleteTopicByID500TextResponse(err.Error()), nil
	}
	return openapi.DeleteTopicByID200TextResponse("ok"), nil
}
