package repo

import (
	"context"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

type MessagesV2Groups interface {
	Create(context.Context, factcheck.MessageGroup, ...Option) (factcheck.MessageGroup, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.MessageGroup, error)
	GetBySHA1(ctx context.Context, sha1 string, opts ...Option) (factcheck.MessageGroup, error)
	AssignTopic(ctx context.Context, id string, topicID string, opts ...Option) (factcheck.MessageGroup, error)
	UnassignTopic(ctx context.Context, id string, opts ...Option) (factcheck.MessageGroup, error)
	Delete(ctx context.Context, id string, opts ...Option) error
}

func NewMessagesV2Groups(queries *postgres.Queries) MessagesV2Groups {
	return &messagesV2Groups{queries: queries}
}

type messagesV2Groups struct {
	queries *postgres.Queries
}

func (m *messagesV2Groups) Create(ctx context.Context, group factcheck.MessageGroup, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	if group.Text == "" {
		slog.WarnContext(ctx, "empty group.text", "group_id", group.ID)
	}
	if group.TextSHA1 == "" {
		slog.WarnContext(ctx, "empty group.text_sha1", "group_id", group.ID)
	}
	params, err := postgres.MessageV2GroupCreator(group)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	created, err := queries.CreateMessageV2Group(ctx, params)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	return postgres.ToMessageV2Group(created)
}

func (m *messagesV2Groups) GetByID(ctx context.Context, id string, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	result, err := queries.GetMessageV2Group(ctx, uuid)
	if err != nil {
		return factcheck.MessageGroup{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToMessageV2Group(result)
}

func (m *messagesV2Groups) ListByTopic(ctx context.Context, topicID string, opts ...Option) ([]factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return nil, err
	}
	result, err := queries.ListMessageV2GroupsByTopic(ctx, topicUUID)
	if err != nil {
		return nil, err
	}
	return postgres.ToMessageV2Groups(result)
}

func (m *messagesV2Groups) AssignTopic(ctx context.Context, id string, topicID string, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	result, err := queries.AssignMessageV2GroupToTopic(ctx, postgres.AssignMessageV2GroupToTopicParams{
		ID:      uuid,
		TopicID: topicUUID,
	})
	if err != nil {
		return factcheck.MessageGroup{}, handleNotFound(err, filter{
			"id":       id,
			"topic_id": topicID,
		})
	}
	return postgres.ToMessageV2Group(result)
}

func (m *messagesV2Groups) UnassignTopic(ctx context.Context, id string, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	result, err := queries.UnassignMessageV2GroupFromTopic(ctx, uuid)
	if err != nil {
		return factcheck.MessageGroup{}, handleNotFound(err, filter{"id": id})
	}
	return postgres.ToMessageV2Group(result)
}

func (m *messagesV2Groups) GetBySHA1(ctx context.Context, sha1 string, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	sha1Text, err := postgres.Text(sha1)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	result, err := queries.GetMessageV2GroupBySHA1(ctx, sha1Text)
	if err != nil {
		return factcheck.MessageGroup{}, handleNotFound(err, filter{"sha1": sha1})
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
