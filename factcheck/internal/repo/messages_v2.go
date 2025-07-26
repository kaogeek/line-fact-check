package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

type MessagesV2 interface {
	Create(ctx context.Context, message factcheck.MessageV2, opts ...Option) (factcheck.MessageV2, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.MessageV2, error)
	ListByTopic(ctx context.Context, topicID string, opts ...Option) ([]factcheck.MessageV2, error)
	AssignTopic(ctx context.Context, messageID string, topicID string, opts ...Option) (factcheck.MessageV2, error)
	UnassignTopic(ctx context.Context, messageID string, opts ...Option) (factcheck.MessageV2, error)
	ListByGroup(ctx context.Context, groupID string, opts ...Option) ([]factcheck.MessageV2, error)
	AssignGroup(ctx context.Context, messageID string, groupID string, opts ...Option) (factcheck.MessageV2, error)
	Delete(ctx context.Context, id string, opts ...Option) error
}

func NewMessagesV2(queries *postgres.Queries) MessagesV2 {
	return &messagesV2{queries: queries}
}

type messagesV2 struct {
	queries *postgres.Queries
}

func (m *messagesV2) Create(ctx context.Context, msg factcheck.MessageV2, opts ...Option) (factcheck.MessageV2, error) {
	queries := queries(m.queries, options(opts...))
	params, err := postgres.MessageV2Creator(msg)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	created, err := queries.CreateMessageV2(ctx, params)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	return postgres.ToMessageV2(created)
}

func (m *messagesV2) GetByID(ctx context.Context, id string, opts ...Option) (factcheck.MessageV2, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	result, err := queries.GetMessageV2(ctx, uuid)
	if err != nil {
		return factcheck.MessageV2{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToMessageV2(result)
}

func (m *messagesV2) ListByTopic(ctx context.Context, topicID string, opts ...Option) ([]factcheck.MessageV2, error) {
	queries := queries(m.queries, options(opts...))
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return nil, err
	}
	list, err := queries.ListMessagesV2ByTopic(ctx, topicUUID)
	if err != nil {
		return nil, err
	}
	return utils.Map(list, postgres.ToMessageV2)
}

func (m *messagesV2) ListByGroup(ctx context.Context, groupID string, opts ...Option) ([]factcheck.MessageV2, error) {
	queries := queries(m.queries, options(opts...))
	groupUUID, err := postgres.UUID(groupID)
	if err != nil {
		return nil, err
	}
	list, err := queries.ListMessagesV2ByGroup(ctx, groupUUID)
	if err != nil {
		return nil, err
	}
	return utils.Map(list, postgres.ToMessageV2)
}

func (m *messagesV2) Delete(ctx context.Context, id string, opts ...Option) error {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return err
	}
	err = queries.DeleteMessageV2(ctx, uuid)
	if err != nil {
		return handleNotFound(err, map[string]string{"id": id})
	}
	return nil
}

func (m *messagesV2) AssignTopic(ctx context.Context, messageID string, topicID string, opts ...Option) (factcheck.MessageV2, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(messageID)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	msg, err := queries.AssignMessageV2ToTopic(ctx, postgres.AssignMessageV2ToTopicParams{
		ID:      uuid,
		TopicID: topicUUID,
	})
	if err != nil {
		return factcheck.MessageV2{}, handleNotFound(err, map[string]string{"message_id": messageID, "topic_id": topicID})
	}
	return postgres.ToMessageV2(msg)
}

func (m *messagesV2) UnassignTopic(ctx context.Context, messageID string, opts ...Option) (factcheck.MessageV2, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(messageID)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	msg, err := queries.UnassignMessageV2FromTopic(ctx, uuid)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	return postgres.ToMessageV2(msg)
}

func (m *messagesV2) AssignGroup(ctx context.Context, messageID string, groupID string, opts ...Option) (factcheck.MessageV2, error) {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(messageID)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	groupUUID, err := postgres.UUID(groupID)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	msg, err := queries.AssignMessageV2ToMessageGroup(ctx, postgres.AssignMessageV2ToMessageGroupParams{
		ID:      uuid,
		GroupID: groupUUID,
	})
	if err != nil {
		return factcheck.MessageV2{}, handleNotFound(err, map[string]string{"message_id": messageID, "topic_id": groupID})
	}
	return postgres.ToMessageV2(msg)
}
