package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// UserMessages defines the interface for user message data operations
type UserMessages interface {
	Create(ctx context.Context, userMessage factcheck.UserMessage, opts ...Option) (factcheck.UserMessage, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.UserMessage, error)
	Delete(ctx context.Context, id string, opts ...Option) error
}

// userMessages implements RepositoryUserMessage
type userMessages struct {
	queries *postgres.Queries
}

// NewUserMessages creates a new user message repository
func NewUserMessages(queries *postgres.Queries) UserMessages {
	return &userMessages{queries: queries}
}

// Create creates a new user message using the userMessage adapter
func (u *userMessages) Create(ctx context.Context, um factcheck.UserMessage, opts ...Option) (factcheck.UserMessage, error) {
	queries := queries(u.queries, options(opts...))
	params, err := postgres.UserMessageCreator(um)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	dbUserMessage, err := queries.CreateUserMessage(ctx, params)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	return postgres.ToUserMessage(dbUserMessage)
}

// GetByID retrieves a user message by ID using the userMessageDomain adapter
func (u *userMessages) GetByID(ctx context.Context, id string, opts ...Option) (factcheck.UserMessage, error) {
	queries := queries(u.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	dbUserMessage, err := queries.GetUserMessage(ctx, uuid)
	if err != nil {
		return factcheck.UserMessage{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToUserMessage(dbUserMessage)
}

// Delete deletes a user message by ID using the stringToUUID adapter
func (u *userMessages) Delete(ctx context.Context, id string, opts ...Option) error {
	queries := queries(u.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return err
	}
	err = queries.DeleteUserMessage(ctx, uuid)
	if err != nil {
		return handleNotFound(err, map[string]string{"id": id})
	}
	return nil
}
