package repo

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// topic converts a factcheck.Topic to postgres.CreateTopicParams
func topic(topic factcheck.Topic) (postgres.CreateTopicParams, error) {
	// Convert string ID to UUID
	var topicID pgtype.UUID
	if err := topicID.Scan(topic.ID); err != nil {
		return postgres.CreateTopicParams{}, err
	}

	// Convert timestamps
	createdAt := pgtype.Timestamptz{}
	if err := createdAt.Scan(topic.CreatedAt); err != nil {
		return postgres.CreateTopicParams{}, err
	}

	var updatedAt pgtype.Timestamptz
	if topic.UpdatedAt != nil {
		if err := updatedAt.Scan(*topic.UpdatedAt); err != nil {
			return postgres.CreateTopicParams{}, err
		}
	}

	// Convert optional fields
	var result pgtype.Text
	if topic.Result != "" {
		if err := result.Scan(topic.Result); err != nil {
			return postgres.CreateTopicParams{}, err
		}
	}

	var resultStatus pgtype.Text
	if topic.ResultStatus != "" {
		if err := resultStatus.Scan(string(topic.ResultStatus)); err != nil {
			return postgres.CreateTopicParams{}, err
		}
	}

	return postgres.CreateTopicParams{
		ID:           topicID,
		Name:         topic.Name,
		Status:       string(topic.Status),
		Result:       result,
		ResultStatus: resultStatus,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

// topicUpdate converts a factcheck.Topic to postgres.UpdateTopicParams
func topicUpdate(topic factcheck.Topic) (postgres.UpdateTopicParams, error) {
	// Convert string ID to UUID
	var topicID pgtype.UUID
	if err := topicID.Scan(topic.ID); err != nil {
		return postgres.UpdateTopicParams{}, err
	}

	// Convert timestamps
	var updatedAt pgtype.Timestamptz
	if topic.UpdatedAt != nil {
		if err := updatedAt.Scan(*topic.UpdatedAt); err != nil {
			return postgres.UpdateTopicParams{}, err
		}
	}

	// Convert optional fields
	var result pgtype.Text
	if topic.Result != "" {
		if err := result.Scan(topic.Result); err != nil {
			return postgres.UpdateTopicParams{}, err
		}
	}

	var resultStatus pgtype.Text
	if topic.ResultStatus != "" {
		if err := resultStatus.Scan(string(topic.ResultStatus)); err != nil {
			return postgres.UpdateTopicParams{}, err
		}
	}

	return postgres.UpdateTopicParams{
		ID:           topicID,
		Name:         topic.Name,
		Status:       string(topic.Status),
		Result:       result,
		ResultStatus: resultStatus,
		UpdatedAt:    updatedAt,
	}, nil
}

func topicDomain(dbTopic postgres.Topic) factcheck.Topic {
	topic := factcheck.Topic{
		Name:   dbTopic.Name,
		Status: factcheck.StatusTopic(dbTopic.Status),
	}

	// Convert UUID to string
	if dbTopic.ID.Valid {
		topic.ID = dbTopic.ID.String()
	}

	// Convert optional text fields
	if dbTopic.Result.Valid {
		topic.Result = dbTopic.Result.String
	}

	if dbTopic.ResultStatus.Valid {
		topic.ResultStatus = factcheck.StatusTopicResult(dbTopic.ResultStatus.String)
	}

	// Convert timestamps
	if dbTopic.CreatedAt.Valid {
		topic.CreatedAt = dbTopic.CreatedAt.Time
	}

	if dbTopic.UpdatedAt.Valid {
		topic.UpdatedAt = &dbTopic.UpdatedAt.Time
	}

	return topic
}

// message converts a factcheck.Message to postgres.CreateMessageParams
func message(m factcheck.Message) (postgres.CreateMessageParams, error) {
	// Convert string IDs to UUIDs
	var messageID pgtype.UUID
	if err := messageID.Scan(m.ID); err != nil {
		return postgres.CreateMessageParams{}, err
	}

	var topicID pgtype.UUID
	if err := topicID.Scan(m.TopicID); err != nil {
		return postgres.CreateMessageParams{}, err
	}

	// Convert timestamps
	createdAt := pgtype.Timestamptz{}
	if err := createdAt.Scan(m.CreatedAt); err != nil {
		return postgres.CreateMessageParams{}, err
	}

	var updatedAt pgtype.Timestamptz
	if m.UpdatedAt != nil {
		if err := updatedAt.Scan(*m.UpdatedAt); err != nil {
			return postgres.CreateMessageParams{}, err
		}
	}

	return postgres.CreateMessageParams{
		ID:        messageID,
		TopicID:   topicID,
		Text:      m.Text,
		Type:      string(m.Type),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// messageUpdate converts a factcheck.Message to postgres.UpdateMessageParams
func messageUpdate(m factcheck.Message) (postgres.UpdateMessageParams, error) {
	// Convert string ID to UUID
	var messageID pgtype.UUID
	if err := messageID.Scan(m.ID); err != nil {
		return postgres.UpdateMessageParams{}, err
	}

	// Convert timestamps
	var updatedAt pgtype.Timestamptz
	if m.UpdatedAt != nil {
		if err := updatedAt.Scan(*m.UpdatedAt); err != nil {
			return postgres.UpdateMessageParams{}, err
		}
	}

	return postgres.UpdateMessageParams{
		ID:        messageID,
		Text:      m.Text,
		Type:      string(m.Type),
		UpdatedAt: updatedAt,
	}, nil
}

func messageDomain(dbMessage postgres.Message) factcheck.Message {
	message := factcheck.Message{
		Text: dbMessage.Text,
		Type: factcheck.TypeMessage(dbMessage.Type),
	}

	// Convert UUIDs to strings
	if dbMessage.ID.Valid {
		message.ID = dbMessage.ID.String()
	}

	if dbMessage.TopicID.Valid {
		message.TopicID = dbMessage.TopicID.String()
	}

	// Convert timestamps
	if dbMessage.CreatedAt.Valid {
		message.CreatedAt = dbMessage.CreatedAt.Time
	}

	if dbMessage.UpdatedAt.Valid {
		message.UpdatedAt = &dbMessage.UpdatedAt.Time
	}

	return message
}

// UserMessageAdapters converts between postgres.UserMessage and factcheck.UserMessage

// userMessage converts a factcheck.UserMessage[json.RawMessage] to postgres.CreateUserMessageParams
func userMessage(u factcheck.UserMessage[json.RawMessage]) (postgres.CreateUserMessageParams, error) {
	// Convert string IDs to UUIDs
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(u.MessageID); err != nil {
		return postgres.CreateUserMessageParams{}, err
	}

	var messageID pgtype.UUID
	if err := messageID.Scan(u.MessageID); err != nil {
		return postgres.CreateUserMessageParams{}, err
	}

	// Convert timestamps
	createdAt := pgtype.Timestamptz{}
	if err := createdAt.Scan(u.CreatedAt); err != nil {
		return postgres.CreateUserMessageParams{}, err
	}

	var updatedAt pgtype.Timestamptz
	if u.UpdatedAt != nil {
		if err := updatedAt.Scan(*u.UpdatedAt); err != nil {
			return postgres.CreateUserMessageParams{}, err
		}
	}

	var repliedAt pgtype.Timestamptz
	if u.RepliedAt != nil {
		if err := repliedAt.Scan(*u.RepliedAt); err != nil {
			return postgres.CreateUserMessageParams{}, err
		}
	}

	// Convert metadata to JSON
	metadata, err := json.Marshal(u.Metadata)
	if err != nil {
		return postgres.CreateUserMessageParams{}, err
	}

	return postgres.CreateUserMessageParams{
		ID:        userMessageID,
		RepliedAt: repliedAt,
		MessageID: messageID,
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// userMessageUpdate converts a factcheck.UserMessage[json.RawMessage] to postgres.UpdateUserMessageParams
func userMessageUpdate(userMessage factcheck.UserMessage[json.RawMessage]) (postgres.UpdateUserMessageParams, error) {
	// Convert string ID to UUID
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(userMessage.MessageID); err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}

	// Convert timestamps
	var updatedAt pgtype.Timestamptz
	if userMessage.UpdatedAt != nil {
		if err := updatedAt.Scan(*userMessage.UpdatedAt); err != nil {
			return postgres.UpdateUserMessageParams{}, err
		}
	}

	var repliedAt pgtype.Timestamptz
	if userMessage.RepliedAt != nil {
		if err := repliedAt.Scan(*userMessage.RepliedAt); err != nil {
			return postgres.UpdateUserMessageParams{}, err
		}
	}

	// Convert metadata to JSON
	metadata, err := json.Marshal(userMessage.Metadata)
	if err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}

	return postgres.UpdateUserMessageParams{
		ID:        userMessageID,
		RepliedAt: repliedAt,
		Metadata:  metadata,
		UpdatedAt: updatedAt,
	}, nil
}

// userMessageDomain converts a postgres.UserMessage to factcheck.UserMessage[json.RawMessage]
func userMessageDomain(dbUserMessage postgres.UserMessage) (factcheck.UserMessage[json.RawMessage], error) {
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

// Utility functions for common conversions

// stringToUUID converts a string to pgtype.UUID
func stringToUUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	return uuid, err
}

// timestamptz converts a time.Time to pgtype.Timestamptz
func timestamptz(t time.Time) (pgtype.Timestamptz, error) {
	var timestamptz pgtype.Timestamptz
	err := timestamptz.Scan(t)
	return timestamptz, err
}

// timestamptzNullable converts a *time.Time to pgtype.Timestamptz
func timestamptzNullable(t *time.Time) (pgtype.Timestamptz, error) {
	var timestamptz pgtype.Timestamptz
	if t != nil {
		err := timestamptz.Scan(*t)
		return timestamptz, err
	}
	return timestamptz, nil
}
