package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/models/postgres"
)

// RepositoryMessage defines the interface for message data operations
type RepositoryMessage interface {
	Create(ctx context.Context, message factcheck.Message) (factcheck.Message, error)
	GetByID(ctx context.Context, id string) (factcheck.Message, error)
	ListByTopic(ctx context.Context, topicID string) ([]factcheck.Message, error)
	Update(ctx context.Context, message factcheck.Message) (factcheck.Message, error)
	Delete(ctx context.Context, id string) error
}

// repositoryMessage implements RepositoryMessage
type repositoryMessage struct {
	queries *postgres.Queries
}

// NewRepositoryMessage creates a new message repository
func NewRepositoryMessage(queries *postgres.Queries) RepositoryMessage {
	return &repositoryMessage{
		queries: queries,
	}
}

// Create creates a new message using the message adapter
func (r *repositoryMessage) Create(ctx context.Context, msg factcheck.Message) (factcheck.Message, error) {
	params, err := message(msg)
	if err != nil {
		return factcheck.Message{}, err
	}

	dbMessage, err := r.queries.CreateMessage(ctx, params)
	if err != nil {
		return factcheck.Message{}, err
	}

	return messageDomain(dbMessage), nil
}

// GetByID retrieves a message by ID using the messageDomain adapter
func (r *repositoryMessage) GetByID(ctx context.Context, id string) (factcheck.Message, error) {
	messageID, err := stringToUUID(id)
	if err != nil {
		return factcheck.Message{}, err
	}

	dbMessage, err := r.queries.GetMessage(ctx, messageID)
	if err != nil {
		return factcheck.Message{}, err
	}

	return messageDomain(dbMessage), nil
}

// ListByTopic retrieves messages by topic ID using the messageDomain adapter
func (r *repositoryMessage) ListByTopic(ctx context.Context, topicID string) ([]factcheck.Message, error) {
	topicUUID, err := stringToUUID(topicID)
	if err != nil {
		return nil, err
	}

	dbMessages, err := r.queries.ListMessagesByTopic(ctx, topicUUID)
	if err != nil {
		return nil, err
	}

	messages := make([]factcheck.Message, len(dbMessages))
	for i, dbMessage := range dbMessages {
		messages[i] = messageDomain(dbMessage)
	}

	return messages, nil
}

// Update updates a message using the messageUpdate adapter
func (r *repositoryMessage) Update(ctx context.Context, msg factcheck.Message) (factcheck.Message, error) {
	params, err := messageUpdate(msg)
	if err != nil {
		return factcheck.Message{}, err
	}

	dbMessage, err := r.queries.UpdateMessage(ctx, params)
	if err != nil {
		return factcheck.Message{}, err
	}

	return messageDomain(dbMessage), nil
}

// Delete deletes a message by ID using the stringToUUID adapter
func (r *repositoryMessage) Delete(ctx context.Context, id string) error {
	messageID, err := stringToUUID(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteMessage(ctx, messageID)
}
