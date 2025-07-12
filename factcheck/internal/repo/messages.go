package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// Messages defines the interface for message data operations
type Messages interface {
	Create(ctx context.Context, message factcheck.Message) (factcheck.Message, error)
	GetByID(ctx context.Context, id string) (factcheck.Message, error)
	ListByTopic(ctx context.Context, topicID string) ([]factcheck.Message, error)
	Update(ctx context.Context, message factcheck.Message) (factcheck.Message, error)
	Delete(ctx context.Context, id string) error
}

// messages implements RepositoryMessage
type messages struct {
	queries *postgres.Queries
}

// NewRepositoryMessage creates a new message repository
func NewRepositoryMessage(queries *postgres.Queries) Messages {
	return &messages{
		queries: queries,
	}
}

// Create creates a new message using the message adapter
func (m *messages) Create(ctx context.Context, msg factcheck.Message) (factcheck.Message, error) {
	params, err := message(msg)
	if err != nil {
		return factcheck.Message{}, err
	}
	dbMessage, err := m.queries.CreateMessage(ctx, params)
	if err != nil {
		return factcheck.Message{}, err
	}
	return messageDomain(dbMessage), nil
}

// GetByID retrieves a message by ID using the messageDomain adapter
func (m *messages) GetByID(ctx context.Context, id string) (factcheck.Message, error) {
	messageID, err := uuid(id)
	if err != nil {
		return factcheck.Message{}, err
	}
	dbMessage, err := m.queries.GetMessage(ctx, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return factcheck.Message{}, &ErrNotFound{
				Err:    err,
				Filter: map[string]string{"id": id},
			}
		}
		return factcheck.Message{}, err
	}
	return messageDomain(dbMessage), nil
}

// ListByTopic retrieves messages by topic ID using the messageDomain adapter
func (m *messages) ListByTopic(ctx context.Context, topicID string) ([]factcheck.Message, error) {
	topicUUID, err := uuid(topicID)
	if err != nil {
		return nil, err
	}
	dbMessages, err := m.queries.ListMessagesByTopic(ctx, topicUUID)
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
func (m *messages) Update(ctx context.Context, msg factcheck.Message) (factcheck.Message, error) {
	params, err := messageUpdate(msg)
	if err != nil {
		return factcheck.Message{}, err
	}
	dbMessage, err := m.queries.UpdateMessage(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return factcheck.Message{}, &ErrNotFound{
				Err:    err,
				Filter: map[string]string{"id": msg.ID},
			}
		}
		return factcheck.Message{}, err
	}
	return messageDomain(dbMessage), nil
}

// Delete deletes a message by ID using the stringToUUID adapter
func (m *messages) Delete(ctx context.Context, id string) error {
	messageID, err := uuid(id)
	if err != nil {
		return err
	}
	err = m.queries.DeleteMessage(ctx, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &ErrNotFound{
				Err:    err,
				Filter: map[string]string{"id": id},
			}
		}
		return err
	}
	return nil
}
