package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// Messages defines the interface for message data operations
type Messages interface {
	Create(ctx context.Context, message factcheck.Message, opts ...Option) (factcheck.Message, error)
	GetByID(ctx context.Context, id string) (factcheck.Message, error)
	ListByTopic(ctx context.Context, topicID string) ([]factcheck.Message, error)
	Update(ctx context.Context, message factcheck.Message) (factcheck.Message, error)
	AssignTopic(ctx context.Context, messageID string, topicID string) (factcheck.Message, error)
	Delete(ctx context.Context, id string) error
}

// messages implements RepositoryMessage
type messages struct {
	queries *postgres.Queries
}

// NewMessages creates a new message repository
func NewMessages(queries *postgres.Queries) Messages {
	return &messages{
		queries: queries,
	}
}

// Create creates a new message using the message adapter
func (m *messages) Create(ctx context.Context, msg factcheck.Message, opts ...Option) (factcheck.Message, error) {
	options := Options{}
	for i := range opts {
		options = opts[i](options)
	}
	params, err := postgres.MessageCreator(msg)
	if err != nil {
		return factcheck.Message{}, err
	}
	query := m.queries.CreateMessage
	if options.tx != nil {
		query = m.queries.WithTx(options.tx).CreateMessage
	}
	dbMessage, err := query(ctx, params)
	if err != nil {
		return factcheck.Message{}, err
	}
	return postgres.ToMessage(dbMessage), nil
}

// GetByID retrieves a message by ID using the messageDomain adapter
func (m *messages) GetByID(ctx context.Context, id string) (factcheck.Message, error) {
	messageID, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Message{}, err
	}
	dbMessage, err := m.queries.GetMessage(ctx, messageID)
	if err != nil {
		return factcheck.Message{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToMessage(dbMessage), nil
}

// ListByTopic retrieves messages by topic ID using the messageDomain adapter
func (m *messages) ListByTopic(ctx context.Context, topicID string) ([]factcheck.Message, error) {
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return nil, err
	}
	dbMessages, err := m.queries.ListMessagesByTopic(ctx, topicUUID)
	if err != nil {
		return nil, err
	}
	messages := make([]factcheck.Message, len(dbMessages))
	for i, dbMessage := range dbMessages {
		messages[i] = postgres.ToMessage(dbMessage)
	}
	return messages, nil
}

// Update updates a message using the messageUpdate adapter
func (m *messages) Update(ctx context.Context, msg factcheck.Message) (factcheck.Message, error) {
	params, err := postgres.MessageUpdater(msg)
	if err != nil {
		return factcheck.Message{}, err
	}
	dbMessage, err := m.queries.UpdateMessage(ctx, params)
	if err != nil {
		return factcheck.Message{}, handleNotFound(err, map[string]string{"id": msg.ID})
	}
	return postgres.ToMessage(dbMessage), nil
}

// AssignTopic assigns a message to a different topic
func (m *messages) AssignTopic(ctx context.Context, messageID string, topicID string) (factcheck.Message, error) {
	// Convert string IDs to pgtype.UUID
	msgUUID, err := postgres.UUID(messageID)
	if err != nil {
		return factcheck.Message{}, err
	}
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return factcheck.Message{}, err
	}

	// Update the message's topic_id
	dbMessage, err := m.queries.AssignMessageToTopic(ctx, postgres.AssignMessageToTopicParams{
		ID:      msgUUID,
		TopicID: topicUUID,
	})
	if err != nil {
		return factcheck.Message{}, handleNotFound(err, map[string]string{"message_id": messageID, "topic_id": topicID})
	}

	return postgres.ToMessage(dbMessage), nil
}

// Delete deletes a message by ID using the stringToUUID adapter
func (m *messages) Delete(ctx context.Context, id string) error {
	messageID, err := postgres.UUID(id)
	if err != nil {
		return err
	}
	err = m.queries.DeleteMessage(ctx, messageID)
	if err != nil {
		return handleNotFound(err, map[string]string{"id": id})
	}
	return nil
}
