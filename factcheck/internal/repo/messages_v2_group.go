package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

type MessagesV2Groups interface {
	Create(context.Context, factcheck.MessageV2Group, ...Option) (factcheck.MessageV2Group, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.MessageV2Group, error)
	AssignTopic(ctx context.Context, id string, topicID string, opts ...Option) (factcheck.MessageV2Group, error)
	UnassignTopic(ctx context.Context, id string, opts ...Option) (factcheck.MessageV2Group, error)
	Delete(ctx context.Context, id string, opts ...Option) error
}

func NewMessagesV2Groups(queries *postgres.Queries) MessagesV2Groups {
	return &messagesV2Groups{queries: queries}
}

type messagesV2Groups struct {
	queries *postgres.Queries
}

func (m *messagesV2Groups) Create(ctx context.Context, group factcheck.MessageV2Group, opts ...Option) (factcheck.MessageV2Group, error) {
	queries := queries(m.queries, options(opts...))
	params, err := postgres.MessageV2GroupCreator(group)
	if err != nil {
		return factcheck.MessageV2Group{}, err
	}
	created, err := queries.CreateMessageV2Group(ctx, params)
	if err != nil {
		return factcheck.MessageV2Group{}, err
	}
	return postgres.ToMessageV2Group(created)
}

func (m *messagesV2Groups) GetByID(ctx context.Context, id string, opts ...Option) (factcheck.MessageV2Group, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageV2Group{}, err
	}
	result, err := queries.GetMessageV2Group(ctx, uuid)
	if err != nil {
		return factcheck.MessageV2Group{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToMessageV2Group(result)
}

func (m *messagesV2Groups) AssignTopic(ctx context.Context, id string, topicID string, opts ...Option) (factcheck.MessageV2Group, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageV2Group{}, err
	}
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return factcheck.MessageV2Group{}, err
	}
	result, err := queries.AssignMessageV2GroupToTopic(ctx, postgres.AssignMessageV2GroupToTopicParams{
		ID:      uuid,
		TopicID: topicUUID,
	})
	if err != nil {
		return factcheck.MessageV2Group{}, handleNotFound(err, filter{
			"id":       id,
			"topic_id": topicID,
		})
	}
	return postgres.ToMessageV2Group(result)
}

func (m *messagesV2Groups) UnassignTopic(ctx context.Context, id string, opts ...Option) (factcheck.MessageV2Group, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageV2Group{}, err
	}
	result, err := queries.UnassignMessageV2GroupFromTopic(ctx, uuid)
	if err != nil {
		return factcheck.MessageV2Group{}, handleNotFound(err, filter{"id": id})
	}
	return postgres.ToMessageV2Group(result)
}

func (m *messagesV2Groups) Delete(ctx context.Context, id string, opts ...Option) error {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return err
	}
	err = queries.DeleteMessageV2Group(ctx, uuid)
	return handleNotFound(err, filter{"id": id})
}
