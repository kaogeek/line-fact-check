package repo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	gotime "time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TopicCreator(topic factcheck.Topic) (postgres.CreateTopicParams, error) {
	id, err := uuid(topic.ID)
	if err != nil {
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
		ID:           id,
		Name:         topic.Name,
		Description:  topic.Description,
		Status:       string(topic.Status),
		Result:       result,
		ResultStatus: resultStatus,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func ToTopic(data postgres.Topic) factcheck.Topic {
	topic := factcheck.Topic{
		Name:        data.Name,
		Description: data.Description,
		Status:      factcheck.StatusTopic(data.Status),
	}
	if data.ID.Valid {
		topic.ID = data.ID.String()
	}
	if data.Result.Valid {
		topic.Result = data.Result.String
	}
	if data.ResultStatus.Valid {
		topic.ResultStatus = factcheck.StatusTopicResult(data.ResultStatus.String)
	}
	if data.CreatedAt.Valid {
		topic.CreatedAt = data.CreatedAt.Time
	}
	if data.UpdatedAt.Valid {
		topic.UpdatedAt = &data.UpdatedAt.Time
	}
	return topic
}

func ToTopics(topics []postgres.Topic) []factcheck.Topic {
	return utils.MapSliceNoError(topics, ToTopic)
}

func MessageCreator(m factcheck.Message) (postgres.CreateMessageParams, error) {
	id, err := uuid(m.ID)
	if err != nil {
		return postgres.CreateMessageParams{}, err
	}
	userMessageID, err := uuid(m.UserMessageID)
	if err != nil {
		return postgres.CreateMessageParams{}, err
	}
	// Could be nil
	topicID, _ := uuid(m.TopicID)
	createdAt, err := timestamptz(m.CreatedAt)
	if err != nil {
		return postgres.CreateMessageParams{}, err
	}
	updatedAt, err := timestamptzNullable(m.UpdatedAt)
	if err != nil {
		return postgres.CreateMessageParams{}, err
	}
	return postgres.CreateMessageParams{
		ID:            id,
		UserMessageID: userMessageID,
		Type:          string(m.Type),
		Status:        string(m.Status),
		TopicID:       topicID,
		Text:          m.Text,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}, nil
}

func MessageUpdater(m factcheck.Message) (postgres.UpdateMessageParams, error) {
	id, err := uuid(m.ID)
	if err != nil {
		return postgres.UpdateMessageParams{}, err
	}
	updatedAt, err := timestamptzNullable(m.UpdatedAt)
	if err != nil {
		return postgres.UpdateMessageParams{}, err
	}
	return postgres.UpdateMessageParams{
		ID:        id,
		Text:      m.Text,
		Type:      string(m.Type),
		Status:    string(m.Status),
		UpdatedAt: updatedAt,
	}, nil
}

func ToMessage(data postgres.Message) factcheck.Message {
	message := factcheck.Message{
		Text:   data.Text,
		Type:   factcheck.TypeMessage(data.Type),
		Status: factcheck.StatusMessage(data.Status),
	}
	if data.ID.Valid {
		message.ID = data.ID.String()
	}
	if data.UserMessageID.Valid {
		message.UserMessageID = data.UserMessageID.String()
	}
	if data.TopicID.Valid {
		message.TopicID = data.TopicID.String()
	}
	if data.CreatedAt.Valid {
		message.CreatedAt = data.CreatedAt.Time
	}
	if data.UpdatedAt.Valid {
		message.UpdatedAt = &data.UpdatedAt.Time
	}
	return message
}

func UserMessageCreator(u factcheck.UserMessage) (postgres.CreateUserMessageParams, error) {
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
		Type:      string(u.Type),
		RepliedAt: repliedAt,
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func UserMessageUpdater(um factcheck.UserMessage) (postgres.UpdateUserMessageParams, error) {
	var userMessageID pgtype.UUID
	if err := userMessageID.Scan(um.ID); err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}
	updatedAt, err := timestamptzNullable(um.UpdatedAt)
	if err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}
	repliedAt, err := timestamptzNullable(um.RepliedAt)
	if err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}
	metadata, err := json.Marshal(um.Metadata)
	if err != nil {
		return postgres.UpdateUserMessageParams{}, err
	}
	return postgres.UpdateUserMessageParams{
		ID:        userMessageID,
		Type:      string(um.Type),
		RepliedAt: repliedAt,
		Metadata:  metadata,
		UpdatedAt: updatedAt,
	}, nil
}

func ToUserMessage(data postgres.UserMessage) (factcheck.UserMessage, error) {
	id, err := uuidString(data.ID)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	createdAt, err := time(data.CreatedAt)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	var metadata json.RawMessage
	if len(data.Metadata) > 0 {
		metadata = json.RawMessage(data.Metadata)
	}
	return factcheck.UserMessage{
		ID:        id,
		Type:      factcheck.TypeUserMessage(data.Type),
		RepliedAt: timeNullable(data.RepliedAt),
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: timeNullable(data.UpdatedAt),
	}, nil
}

func ToUserMessages(data []postgres.UserMessage) ([]factcheck.UserMessage, error) {
	return utils.MapSlice(data, ToUserMessage)
}

func uuid(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	return uuid, err
}

func uuidStringNullable(id string) pgtype.UUID {
	//nolint
	uuid, err := uuid(id)
	if err != nil && id != "" {
		slog.Error("unexpected bad uuid '%s' in our system: %w", id, err.Error())
	}
	return uuid
}

func uuidString(id pgtype.UUID) (string, error) {
	if !id.Valid {
		return "", fmt.Errorf("invalid uuid string: %+v", id)
	}
	return id.String(), nil
}

func text(s string) (pgtype.Text, error) {
	var text pgtype.Text
	err := text.Scan(s)
	return text, err
}

func textString(s pgtype.Text) (string, error) {
	if s.Valid {
		return s.String, nil
	}
	return "", fmt.Errorf("invalid database text: %+v", s)
}

func timestamptz(t gotime.Time) (pgtype.Timestamptz, error) {
	var timestamptz pgtype.Timestamptz
	err := timestamptz.Scan(t)
	return timestamptz, err
}

func time(t pgtype.Timestamptz) (gotime.Time, error) {
	if !t.Valid {
		return gotime.Time{}, fmt.Errorf("invalid time %+v", t)
	}
	return t.Time, nil
}

func timestamptzNullable(t *gotime.Time) (pgtype.Timestamptz, error) {
	var timestamptz pgtype.Timestamptz
	if t != nil {
		err := timestamptz.Scan(*t)
		return timestamptz, err
	}
	return timestamptz, nil
}

func timeNullable(t pgtype.Timestamptz) *gotime.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
