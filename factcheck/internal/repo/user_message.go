package repo

import (
	"context"
	"encoding/json"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/models/postgres"
)

// RepositoryUserMessage defines the interface for user message data operations
type RepositoryUserMessage interface {
	Create(ctx context.Context, userMessage factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error)
	GetByID(ctx context.Context, id string) (factcheck.UserMessage[json.RawMessage], error)
	ListByMessage(ctx context.Context, messageID string) ([]factcheck.UserMessage[json.RawMessage], error)
	Update(ctx context.Context, userMessage factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error)
	Delete(ctx context.Context, id string) error
}

// repositoryUserMessage implements RepositoryUserMessage
type repositoryUserMessage struct {
	queries *postgres.Queries
}

// NewRepositoryUserMessage creates a new user message repository
func NewRepositoryUserMessage(queries *postgres.Queries) RepositoryUserMessage {
	return &repositoryUserMessage{
		queries: queries,
	}
}

// Create creates a new user message using the userMessage adapter
func (r *repositoryUserMessage) Create(ctx context.Context, um factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error) {
	params, err := userMessage(um)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	dbUserMessage, err := r.queries.CreateUserMessage(ctx, params)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	return userMessageDomain(dbUserMessage)
}

// GetByID retrieves a user message by ID using the userMessageDomain adapter
func (r *repositoryUserMessage) GetByID(ctx context.Context, id string) (factcheck.UserMessage[json.RawMessage], error) {
	uuid, err := stringToUUID(id)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}
	dbUserMessage, err := r.queries.GetUserMessage(ctx, uuid)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}
	return userMessageDomain(dbUserMessage)
}

func (r *repositoryUserMessage) ListByMessage(ctx context.Context, messageID string) ([]factcheck.UserMessage[json.RawMessage], error) {
	uuid, err := stringToUUID(messageID)
	if err != nil {
		return nil, err
	}

	dbUserMessages, err := r.queries.ListUserMessagesByMessage(ctx, uuid)
	if err != nil {
		return nil, err
	}

	userMessages := make([]factcheck.UserMessage[json.RawMessage], len(dbUserMessages))
	for i, dbUserMessage := range dbUserMessages {
		userMessage, err := userMessageDomain(dbUserMessage)
		if err != nil {
			return nil, err
		}
		userMessages[i] = userMessage
	}
	return userMessages, nil
}

// Update updates a user message using the userMessageUpdate adapter
func (r *repositoryUserMessage) Update(ctx context.Context, um factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error) {
	params, err := userMessageUpdate(um)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	dbUserMessage, err := r.queries.UpdateUserMessage(ctx, params)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	return userMessageDomain(dbUserMessage)
}

// Delete deletes a user message by ID using the stringToUUID adapter
func (r *repositoryUserMessage) Delete(ctx context.Context, id string) error {
	uuid, err := stringToUUID(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteUserMessage(ctx, uuid)
}
