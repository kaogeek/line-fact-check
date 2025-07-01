package repo

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"

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

func (r *repositoryUserMessage) Create(ctx context.Context, userMessage factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error) {
	// Convert string IDs to UUIDs
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(userMessage.MessageID); err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	var messageID pgtype.UUID
	if err := messageID.Scan(userMessage.MessageID); err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	// Convert timestamps
	createdAt := pgtype.Timestamptz{}
	if err := createdAt.Scan(userMessage.CreatedAt); err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	var updatedAt pgtype.Timestamptz
	if userMessage.UpdatedAt != nil {
		if err := updatedAt.Scan(*userMessage.UpdatedAt); err != nil {
			return factcheck.UserMessage[json.RawMessage]{}, err
		}
	}

	var repliedAt pgtype.Timestamptz
	if userMessage.RepliedAt != nil {
		if err := repliedAt.Scan(*userMessage.RepliedAt); err != nil {
			return factcheck.UserMessage[json.RawMessage]{}, err
		}
	}

	// Convert metadata to JSON
	metadata, err := json.Marshal(userMessage.Metadata)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	params := postgres.CreateUserMessageParams{
		ID:        userMessageID,
		RepliedAt: repliedAt,
		MessageID: messageID,
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	dbUserMessage, err := r.queries.CreateUserMessage(ctx, params)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	return r.convertToDomainUserMessage(dbUserMessage)
}

func (r *repositoryUserMessage) GetByID(ctx context.Context, id string) (factcheck.UserMessage[json.RawMessage], error) {
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(id); err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	dbUserMessage, err := r.queries.GetUserMessage(ctx, userMessageID)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	return r.convertToDomainUserMessage(dbUserMessage)
}

func (r *repositoryUserMessage) ListByMessage(ctx context.Context, messageID string) ([]factcheck.UserMessage[json.RawMessage], error) {
	var messageUUID pgtype.UUID
	if err := messageUUID.Scan(messageID); err != nil {
		return nil, err
	}

	dbUserMessages, err := r.queries.ListUserMessagesByMessage(ctx, messageUUID)
	if err != nil {
		return nil, err
	}

	userMessages := make([]factcheck.UserMessage[json.RawMessage], len(dbUserMessages))
	for i, dbUserMessage := range dbUserMessages {
		userMessage, err := r.convertToDomainUserMessage(dbUserMessage)
		if err != nil {
			return nil, err
		}
		userMessages[i] = userMessage
	}

	return userMessages, nil
}

func (r *repositoryUserMessage) Update(ctx context.Context, userMessage factcheck.UserMessage[json.RawMessage]) (factcheck.UserMessage[json.RawMessage], error) {
	// Convert string ID to UUID
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(userMessage.MessageID); err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	// Convert timestamps
	var updatedAt pgtype.Timestamptz
	if userMessage.UpdatedAt != nil {
		if err := updatedAt.Scan(*userMessage.UpdatedAt); err != nil {
			return factcheck.UserMessage[json.RawMessage]{}, err
		}
	}

	var repliedAt pgtype.Timestamptz
	if userMessage.RepliedAt != nil {
		if err := repliedAt.Scan(*userMessage.RepliedAt); err != nil {
			return factcheck.UserMessage[json.RawMessage]{}, err
		}
	}

	// Convert metadata to JSON
	metadata, err := json.Marshal(userMessage.Metadata)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	params := postgres.UpdateUserMessageParams{
		ID:        userMessageID,
		RepliedAt: repliedAt,
		Metadata:  metadata,
		UpdatedAt: updatedAt,
	}

	dbUserMessage, err := r.queries.UpdateUserMessage(ctx, params)
	if err != nil {
		return factcheck.UserMessage[json.RawMessage]{}, err
	}

	return r.convertToDomainUserMessage(dbUserMessage)
}

func (r *repositoryUserMessage) Delete(ctx context.Context, id string) error {
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(id); err != nil {
		return err
	}

	return r.queries.DeleteUserMessage(ctx, userMessageID)
}

// convertToDomainUserMessage converts a database user message to domain user message
func (r *repositoryUserMessage) convertToDomainUserMessage(dbUserMessage postgres.UserMessage) (factcheck.UserMessage[json.RawMessage], error) {
	userMessage := factcheck.UserMessage[json.RawMessage]{}

	// Convert UUIDs to strings
	if dbUserMessage.ID.Valid {
		userMessage.MessageID = dbUserMessage.ID.String()
	}

	if dbUserMessage.MessageID.Valid {
		userMessage.MessageID = dbUserMessage.MessageID.String()
	}

	// Convert timestamps
	if dbUserMessage.CreatedAt.Valid {
		userMessage.CreatedAt = dbUserMessage.CreatedAt.Time
	}

	if dbUserMessage.UpdatedAt.Valid {
		userMessage.UpdatedAt = &dbUserMessage.UpdatedAt.Time
	}

	if dbUserMessage.RepliedAt.Valid {
		userMessage.RepliedAt = &dbUserMessage.RepliedAt.Time
	}

	// Convert metadata from JSON
	if len(dbUserMessage.Metadata) > 0 {
		userMessage.Metadata = json.RawMessage(dbUserMessage.Metadata)
	}

	return userMessage, nil
}
