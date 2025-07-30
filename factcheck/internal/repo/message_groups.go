package repo

import (
	"context"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
)

type MessageGroups interface {
	Create(context.Context, factcheck.MessageGroup, ...Option) (factcheck.MessageGroup, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.MessageGroup, error)
	GetBySHA1(ctx context.Context, sha1 string, opts ...Option) (factcheck.MessageGroup, error)
	ListByTopic(ctx context.Context, topicID string, opts ...Option) ([]factcheck.MessageGroup, error)
	AssignTopic(ctx context.Context, id string, topicID string, opts ...Option) (factcheck.MessageGroup, error)
	UnassignTopic(ctx context.Context, id string, opts ...Option) (factcheck.MessageGroup, error)
	Delete(ctx context.Context, id string, opts ...Option) error
}

func NewMessageGroups(queries *postgres.Queries) MessageGroups {
	return &messageGroups{queries: queries}
}

type messageGroups struct {
	queries *postgres.Queries
}

func (m *messageGroups) Create(ctx context.Context, group factcheck.MessageGroup, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	if group.Text == "" {
		slog.WarnContext(ctx, "empty group.text", "group_id", group.ID)
	}
	if group.TextSHA1 == "" {
		slog.WarnContext(ctx, "empty group.text_sha1", "group_id", group.ID)
	}
	params, err := postgres.MessageGroupCreator(group)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	created, err := queries.CreateMessageGroup(ctx, params)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	return postgres.ToMessageGroup(created)
}

func (m *messageGroups) GetByID(ctx context.Context, id string, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	result, err := queries.GetMessageGroup(ctx, uuid)
	if err != nil {
		return factcheck.MessageGroup{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToMessageGroup(result)
}

func (m *messageGroups) ListByTopic(ctx context.Context, topicID string, opts ...Option) ([]factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return nil, err
	}
	result, err := queries.ListMessageGroupsByTopic(ctx, topicUUID)
	if err != nil {
		return nil, err
	}
	return postgres.ToMessageGroups(result)
}

func (m *messageGroups) AssignTopic(ctx context.Context, id string, topicID string, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	result, err := queries.AssignMessageGroupToTopic(ctx, postgres.AssignMessageGroupToTopicParams{
		ID:      uuid,
		TopicID: topicUUID,
	})
	if err != nil {
		return factcheck.MessageGroup{}, handleNotFound(err, filter{
			"id":       id,
			"topic_id": topicID,
		})
	}
	return postgres.ToMessageGroup(result)
}

func (m *messageGroups) UnassignTopic(ctx context.Context, id string, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	result, err := queries.UnassignMessageGroupFromTopic(ctx, uuid)
	if err != nil {
		return factcheck.MessageGroup{}, handleNotFound(err, filter{"id": id})
	}
	return postgres.ToMessageGroup(result)
}

func (m *messageGroups) GetBySHA1(ctx context.Context, sha1 string, opts ...Option) (factcheck.MessageGroup, error) {
	queries := queries(m.queries, options(opts...))
	result, err := queries.GetMessageGroupBySHA1(ctx, sha1)
	if err != nil {
		return factcheck.MessageGroup{}, handleNotFound(err, filter{"sha1": sha1})
	}
	return postgres.ToMessageGroup(result)
}

func (m *messageGroups) Delete(ctx context.Context, id string, opts ...Option) error {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return err
	}
	err = queries.DeleteMessageGroup(ctx, uuid)
	return handleNotFound(err, filter{"id": id})
}
