package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// Messages defines the interface for message data operations
type Messages interface {
	Create(ctx context.Context, message factcheck.Message, opts ...Option) (factcheck.Message, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.Message, error)
	ListByTopic(ctx context.Context, topicID string, opts ...Option) ([]factcheck.Message, error)
	AssignTopic(ctx context.Context, messageID string, topicID string, opts ...Option) (factcheck.Message, error)
	Delete(ctx context.Context, id string, opts ...Option) error
}

// messages implements RepositoryMessage
type messages struct {
	queries *postgres.Queries
}

// NewMessages creates a new message repository
func NewMessages(queries *postgres.Queries) Messages {
	return &messages{queries: queries}
}

// Create creates a new message using the message adapter
func (m *messages) Create(ctx context.Context, msg factcheck.Message, opts ...Option) (factcheck.Message, error) {
	queries := queries(m.queries, options(opts...))
	params, err := postgres.MessageCreator(msg)
	if err != nil {
		return factcheck.Message{}, err
	}
	created, err := queries.CreateMessage(ctx, params)
	if err != nil {
		return factcheck.Message{}, err
	}
	return postgres.ToMessage(created), nil
}

// GetByID retrieves a message by ID using the messageDomain adapter
func (m *messages) GetByID(ctx context.Context, id string, opts ...Option) (factcheck.Message, error) {
	queries := queries(m.queries, options(opts...))
	messageID, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Message{}, err
	}
	result, err := queries.GetMessage(ctx, messageID)
	if err != nil {
		return factcheck.Message{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToMessage(result), nil
}

// ListByTopic retrieves messages by topic ID using the messageDomain adapter
func (m *messages) ListByTopic(ctx context.Context, topicID string, opts ...Option) ([]factcheck.Message, error) {
	queries := queries(m.queries, options(opts...))
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return nil, err
	}
	list, err := queries.ListMessagesByTopic(ctx, topicUUID)
	if err != nil {
		return nil, err
	}
	messages := make([]factcheck.Message, len(list))
	for i, dbMessage := range list {
		messages[i] = postgres.ToMessage(dbMessage)
	}
	return messages, nil
}

// AssignTopic assigns a message to a different topic
func (m *messages) AssignTopic(ctx context.Context, messageID string, topicID string, opts ...Option) (factcheck.Message, error) {
	queries := queries(m.queries, options(opts...))
	msgUUID, err := postgres.UUID(messageID)
	if err != nil {
		return factcheck.Message{}, err
	}
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return factcheck.Message{}, err
	}
	msg, err := queries.AssignMessageToTopic(ctx, postgres.AssignMessageToTopicParams{
		ID:      msgUUID,
		TopicID: topicUUID,
	})
	if err != nil {
		return factcheck.Message{}, handleNotFound(err, map[string]string{"message_id": messageID, "topic_id": topicID})
	}
	return postgres.ToMessage(msg), nil
}

// Delete deletes a message by ID using the stringToUUID adapter
func (m *messages) Delete(ctx context.Context, id string, opts ...Option) error {
	queries := queries(m.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return err
	}
	err = queries.DeleteMessage(ctx, uuid)
	if err != nil {
		return handleNotFound(err, map[string]string{"id": id})
	}
	return nil
}
