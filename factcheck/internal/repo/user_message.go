package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// UserMessages defines the interface for user message data operations
type UserMessages interface {
	Create(ctx context.Context, userMessage factcheck.UserMessage, opts ...Option) (factcheck.UserMessage, error)
	GetByID(ctx context.Context, id string) (factcheck.UserMessage, error)
	Update(ctx context.Context, userMessage factcheck.UserMessage) (factcheck.UserMessage, error)
	Delete(ctx context.Context, id string) error
}

// userMessages implements RepositoryUserMessage
type userMessages struct {
	queries *postgres.Queries
}

// NewUserMessages creates a new user message repository
func NewUserMessages(
	queries *postgres.Queries,
) UserMessages {
	return &userMessages{
		queries: queries,
	}
}

// Create creates a new user message using the userMessage adapter
func (u *userMessages) Create(ctx context.Context, um factcheck.UserMessage, opts ...Option) (factcheck.UserMessage, error) {
	options := Options{}
	for i := range opts {
		options = opts[i](options)
	}
	params, err := UserMessageCreator(um)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	query := u.queries.CreateUserMessage
	if options.tx != nil {
		query = u.queries.WithTx(options.tx).CreateUserMessage
	}
	dbUserMessage, err := query(ctx, params)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	return ToUserMessage(dbUserMessage)
}

// GetByID retrieves a user message by ID using the userMessageDomain adapter
func (u *userMessages) GetByID(ctx context.Context, id string) (factcheck.UserMessage, error) {
	uuid, err := uuid(id)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	dbUserMessage, err := u.queries.GetUserMessage(ctx, uuid)
	if err != nil {
		return factcheck.UserMessage{}, handleNotFound(err, map[string]string{"id": id})
	}
	return ToUserMessage(dbUserMessage)
}

// Update updates a user message using the userMessageUpdate adapter
func (u *userMessages) Update(ctx context.Context, um factcheck.UserMessage) (factcheck.UserMessage, error) {
	params, err := UserMessageUpdater(um)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	dbUserMessage, err := u.queries.UpdateUserMessage(ctx, params)
	if err != nil {
		return factcheck.UserMessage{}, handleNotFound(err, map[string]string{"id": um.ID})
	}
	return ToUserMessage(dbUserMessage)
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
