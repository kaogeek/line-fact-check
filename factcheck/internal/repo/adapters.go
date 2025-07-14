package repo

import (
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func topic(topic factcheck.Topic) (postgres.CreateTopicParams, error) {
	var topicID pgtype.UUID
	if err := topicID.Scan(topic.ID); err != nil {
		return postgres.CreateTopicParams{}, err
	}
	createdAt, err := timestamptz(topic.CreatedAt)
	if err != nil {
		return postgres.CreateTopicParams{}, err
	}
	updatedAt, err := timestamptzNullable(topic.UpdatedAt)
	if err != nil {
		return postgres.CreateTopicParams{}, err
	}
	result, err := text(topic.Result)
	if err != nil {
		return postgres.CreateTopicParams{}, err
	}
	resultStatus, err := text(string(topic.ResultStatus))
	if err != nil {
		return postgres.CreateTopicParams{}, err
	}
	return postgres.CreateTopicParams{
		ID:           topicID,
		Name:         topic.Name,
		Description:  topic.Description,
		Status:       string(topic.Status),
		Result:       result,
		ResultStatus: resultStatus,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func topicDomain(dbTopic postgres.Topic) factcheck.Topic {
	topic := factcheck.Topic{
		Name:        dbTopic.Name,
		Description: dbTopic.Description,
		Status:      factcheck.StatusTopic(dbTopic.Status),
	}
	if dbTopic.ID.Valid {
		topic.ID = dbTopic.ID.String()
	}
	if dbTopic.Result.Valid {
		topic.Result = dbTopic.Result.String
	}
	if dbTopic.ResultStatus.Valid {
		topic.ResultStatus = factcheck.StatusTopicResult(dbTopic.ResultStatus.String)
	}
	if dbTopic.CreatedAt.Valid {
		topic.CreatedAt = dbTopic.CreatedAt.Time
	}
	if dbTopic.UpdatedAt.Valid {
		topic.UpdatedAt = &dbTopic.UpdatedAt.Time
	}
	return topic
}

func topicsDomain(topics []postgres.Topic) []factcheck.Topic {
	return utils.MapSlice(topics, topicDomain)
}

func message(m factcheck.Message) (postgres.CreateMessageParams, error) {
	var messageID pgtype.UUID
	if err := messageID.Scan(m.ID); err != nil {
		return postgres.CreateMessageParams{}, err
	}
	var topicID pgtype.UUID
	if err := topicID.Scan(m.TopicID); err != nil {
		return postgres.CreateMessageParams{}, err
	}
	createdAt, err := timestamptz(m.CreatedAt)
	if err != nil {
		return postgres.CreateMessageParams{}, err
	}
	updatedAt, err := timestamptzNullable(m.UpdatedAt)
	if err != nil {
		return postgres.CreateMessageParams{}, err
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

func messageUpdate(m factcheck.Message) (postgres.UpdateMessageParams, error) {
	var messageID pgtype.UUID
	if err := messageID.Scan(m.ID); err != nil {
		return postgres.UpdateMessageParams{}, err
	}
	updatedAt, err := timestamptzNullable(m.UpdatedAt)
	if err != nil {
		return postgres.UpdateMessageParams{}, err
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
	if dbMessage.ID.Valid {
		message.ID = dbMessage.ID.String()
	}
	if dbMessage.TopicID.Valid {
		message.TopicID = dbMessage.TopicID.String()
	}
	if dbMessage.CreatedAt.Valid {
		message.CreatedAt = dbMessage.CreatedAt.Time
	}
	if dbMessage.UpdatedAt.Valid {
		message.UpdatedAt = &dbMessage.UpdatedAt.Time
	}
	return message
}

func userMessage(u factcheck.UserMessage) (postgres.CreateUserMessageParams, error) {
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(u.ID); err != nil {
		return postgres.CreateUserMessageParams{}, err
	}
	createdAt, err := timestamptz(u.CreatedAt)
	if err != nil {
		return postgres.CreateUserMessageParams{}, err
	}
	updatedAt, err := timestamptzNullable(u.UpdatedAt)
	if err != nil {
		return postgres.CreateUserMessageParams{}, err
	}
	repliedAt, err := timestamptzNullable(u.RepliedAt)
	if err != nil {
		return postgres.CreateUserMessageParams{}, err
	}
	metadata, err := json.Marshal(u.Metadata)
	if err != nil {
		return postgres.CreateUserMessageParams{}, err
	}
	return postgres.CreateUserMessageParams{
		ID:        userMessageID,
		RepliedAt: repliedAt,
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func userMessageUpdate(userMessage factcheck.UserMessage) (postgres.UpdateUserMessageParams, error) {
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(userMessage.ID); err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}
	updatedAt, err := timestamptzNullable(userMessage.UpdatedAt)
	if err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}
	repliedAt, err := timestamptzNullable(userMessage.RepliedAt)
	if err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}
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

func userMessageDomain(dbUserMessage postgres.UserMessage) (factcheck.UserMessage, error) {
	userMessage := factcheck.UserMessage{}
	if dbUserMessage.ID.Valid {
		userMessage.ID = dbUserMessage.ID.String()
	}
	if dbUserMessage.CreatedAt.Valid {
		userMessage.CreatedAt = dbUserMessage.CreatedAt.Time
	}
	if dbUserMessage.UpdatedAt.Valid {
		userMessage.UpdatedAt = &dbUserMessage.UpdatedAt.Time
	}
	if dbUserMessage.RepliedAt.Valid {
		userMessage.RepliedAt = &dbUserMessage.RepliedAt.Time
	}
	if len(dbUserMessage.Metadata) > 0 {
		userMessage.Metadata = json.RawMessage(dbUserMessage.Metadata)
	}
	return userMessage, nil
}

func uuid(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	return uuid, err
}

func text(s string) (pgtype.Text, error) {
	var text pgtype.Text
	err := text.Scan(s)
	return text, err
}

func timestamptz(t time.Time) (pgtype.Timestamptz, error) {
	var timestamptz pgtype.Timestamptz
	err := timestamptz.Scan(t)
	return timestamptz, err
}

func timestamptzNullable(t *time.Time) (pgtype.Timestamptz, error) {
	var timestamptz pgtype.Timestamptz
	if t != nil {
		err := timestamptz.Scan(*t)
		return timestamptz, err
	}
	return timestamptz, nil
}
