package repo

import (
	"context"
	"encoding/json"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// UserMessages defines the interface for user message data operations
type UserMessages interface {
	Create(ctx context.Context, userMessage factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error)
	GetByID(ctx context.Context, id string) (factcheck.UserMessage[json.RawMessage], error)
	ListByMessage(ctx context.Context, messageID string) ([]factcheck.UserMessage[json.RawMessage], error)
	Update(ctx context.Context, userMessage factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error)
	Delete(ctx context.Context, id string) error
}

// userMessages implements RepositoryUserMessage
type userMessages struct {
	queries *postgres.Queries
}

// NewUserMessages creates a new user message repository
func NewUserMessages(queries *postgres.Queries) UserMessages {
	return &userMessages{
		queries: queries,
	}
}

// Create creates a new user message using the userMessage adapter
func (u *userMessages) Create(ctx context.Context, um factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error) {
	params, err := userMessage(um)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}
	dbUserMessage, err := u.queries.CreateUserMessage(ctx, params)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}
	return userMessageDomain(dbUserMessage)
}

// GetByID retrieves a user message by ID using the userMessageDomain adapter
func (u *userMessages) GetByID(ctx context.Context, id string) (factcheck.UserMessage[json.RawMessage], error) {
	uuid, err := uuid(id)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}
	dbUserMessage, err := u.queries.GetUserMessage(ctx, uuid)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, handleNotFound(err, map[string]string{"id": id})
	}
	return userMessageDomain(dbUserMessage)
}

func (u *userMessages) ListByMessage(ctx context.Context, messageID string) ([]factcheck.UserMessage[json.RawMessage], error) {
	uuid, err := uuid(messageID)
	if err != nil {
		return nil, err
	}
	dbUserMessages, err := u.queries.ListUserMessagesByMessage(ctx, uuid)
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
func (u *userMessages) Update(ctx context.Context, um factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error) {
	params, err := userMessageUpdate(um)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}
	dbUserMessage, err := u.queries.UpdateUserMessage(ctx, params)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, handleNotFound(err, map[string]string{"id": um.ID})
	}
	return userMessageDomain(dbUserMessage)
}

// Delete deletes a user message by ID using the stringToUUID adapter
func (u *userMessages) Delete(ctx context.Context, id string) error {
	uuid, err := uuid(id)
	if err != nil {
		return err
	}
	err = u.queries.DeleteUserMessage(ctx, uuid)
	if err != nil {
		return handleNotFound(err, map[string]string{"id": id})
	}
	return nil
}
