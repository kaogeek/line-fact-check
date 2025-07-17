package postgres

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TopicCreator(topic factcheck.Topic) (CreateTopicParams, error) {
	id, err := UUID(topic.ID)
	if err != nil {
		return CreateTopicParams{}, err
	}
	createdAt, err := Timestamptz(topic.CreatedAt)
	if err != nil {
		return CreateTopicParams{}, err
	}
	updatedAt, err := TimestamptzNullable(topic.UpdatedAt)
	if err != nil {
		return CreateTopicParams{}, err
	}
	result, err := Text(topic.Result)
	if err != nil {
		return CreateTopicParams{}, err
	}
	resultStatus, err := Text(string(topic.ResultStatus))
	if err != nil {
		return CreateTopicParams{}, err
	}
	return CreateTopicParams{
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

func ToTopic(data Topic) factcheck.Topic {
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

func ToTopics(topics []Topic) []factcheck.Topic {
	return utils.MapSliceNoError(topics, ToTopic)
}

// ToTopicFromRow converts a ListTopicsRow to factcheck.Topic
func ToTopicFromRow(data ListTopicsRow) factcheck.Topic {
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

// ToTopicFromStatusRow converts a ListTopicsByStatusRow to factcheck.Topic
func ToTopicFromStatusRow(data ListTopicsByStatusRow) factcheck.Topic {
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

// ToTopicFromMessageTextRow converts a ListTopicsLikeMessageTextRow to factcheck.Topic
func ToTopicFromMessageTextRow(data ListTopicsLikeMessageTextRow) factcheck.Topic {
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

// ToTopicFromIDRow converts a ListTopicsLikeIDRow to factcheck.Topic
func ToTopicFromIDRow(data ListTopicsLikeIDRow) factcheck.Topic {
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

// ToTopicFromStatusLikeMessageTextRow converts a ListTopicsByStatusLikeMessageTextRow to factcheck.Topic
func ToTopicFromStatusLikeMessageTextRow(data ListTopicsByStatusLikeMessageTextRow) factcheck.Topic {
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

// ToTopicFromStatusLikeIDRow converts a ListTopicsByStatusLikeIDRow to factcheck.Topic
func ToTopicFromStatusLikeIDRow(data ListTopicsByStatusLikeIDRow) factcheck.Topic {
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

// ToTopicFromIDLikeMessageTextRow converts a ListTopicsLikeIDLikeMessageTextRow to factcheck.Topic
func ToTopicFromIDLikeMessageTextRow(data ListTopicsLikeIDLikeMessageTextRow) factcheck.Topic {
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

// ToTopicFromStatusLikeIDLikeMessageTextRow converts a ListTopicsByStatusLikeIDLikeMessageTextRow to factcheck.Topic
func ToTopicFromStatusLikeIDLikeMessageTextRow(data ListTopicsByStatusLikeIDLikeMessageTextRow) factcheck.Topic {
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

func MessageCreator(m factcheck.Message) (CreateMessageParams, error) {
	id, err := UUID(m.ID)
	if err != nil {
		return CreateMessageParams{}, err
	}
	userMessageID, err := UUID(m.UserMessageID)
	if err != nil {
		return CreateMessageParams{}, err
	}
	// Could be nil
	topicID, _ := UUID(m.TopicID)
	createdAt, err := Timestamptz(m.CreatedAt)
	if err != nil {
		return CreateMessageParams{}, err
	}
	updatedAt, err := TimestamptzNullable(m.UpdatedAt)
	if err != nil {
		return CreateMessageParams{}, err
	}
	return CreateMessageParams{
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

func MessageUpdater(m factcheck.Message) (UpdateMessageParams, error) {
	id, err := UUID(m.ID)
	if err != nil {
		return UpdateMessageParams{}, err
	}
	updatedAt, err := TimestamptzNullable(m.UpdatedAt)
	if err != nil {
		return UpdateMessageParams{}, err
	}
	return UpdateMessageParams{
		ID:        id,
		Text:      m.Text,
		Type:      string(m.Type),
		Status:    string(m.Status),
		UpdatedAt: updatedAt,
	}, nil
}

func ToMessage(data Message) factcheck.Message {
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

func UserMessageCreator(u factcheck.UserMessage) (CreateUserMessageParams, error) {
	id, err := UUID(u.ID)
	if err != nil {
		return CreateUserMessageParams{}, err
	}
	createdAt, err := Timestamptz(u.CreatedAt)
	if err != nil {
		return CreateUserMessageParams{}, err
	}
	updatedAt, err := TimestamptzNullable(u.UpdatedAt)
	if err != nil {
		return CreateUserMessageParams{}, err
	}
	repliedAt, err := TimestamptzNullable(u.RepliedAt)
	if err != nil {
		return CreateUserMessageParams{}, err
	}
	metadata, err := json.Marshal(u.Metadata)
	if err != nil {
		return CreateUserMessageParams{}, err
	}
	return CreateUserMessageParams{
		ID:        id,
		Type:      string(u.Type),
		RepliedAt: repliedAt,
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func UserMessageUpdater(um factcheck.UserMessage) (UpdateUserMessageParams, error) {
	id, err := UUID(um.ID)
	if err != nil {
		return UpdateUserMessageParams{}, err
	}
	updatedAt, err := TimestamptzNullable(um.UpdatedAt)
	if err != nil {
		return UpdateUserMessageParams{}, err
	}
	repliedAt, err := TimestamptzNullable(um.RepliedAt)
	if err != nil {
		return UpdateUserMessageParams{}, err
	}
	metadata, err := json.Marshal(um.Metadata)
	if err != nil {
		return UpdateUserMessageParams{}, err
	}
	return UpdateUserMessageParams{
		ID:        id,
		Type:      string(um.Type),
		RepliedAt: repliedAt,
		Metadata:  metadata,
		UpdatedAt: updatedAt,
	}, nil
}

func ToUserMessage(data UserMessage) (factcheck.UserMessage, error) {
	id, err := FromUUID(data.ID)
	if err != nil {
		return factcheck.UserMessage{}, err
	}
	createdAt, err := Time(data.CreatedAt)
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
		RepliedAt: TimeNullable(data.RepliedAt),
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: TimeNullable(data.UpdatedAt),
	}, nil
}

func ToUserMessages(data []UserMessage) ([]factcheck.UserMessage, error) {
	return utils.MapSlice(data, ToUserMessage)
}

func UUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	return uuid, err
}

func UUIDs(ids []string) ([]pgtype.UUID, error) {
	return utils.MapSlice(ids, UUID)
}

func UUIDNullable(id string) pgtype.UUID {
	//nolint
	uuid, err := UUID(id)
	if err != nil && id != "" {
		slog.Error("unexpected bad uuid '%s' in our system: %w", id, err.Error())
	}
	return uuid
}

func FromUUID(id pgtype.UUID) (string, error) {
	if !id.Valid {
		return "", fmt.Errorf("invalid uuid string: %+v", id)
	}
	return id.String(), nil
}

func Text(s string) (pgtype.Text, error) {
	var text pgtype.Text
	err := text.Scan(s)
	return text, err
}

func FromText(s pgtype.Text) (string, error) {
	if s.Valid {
		return s.String, nil
	}
	return "", fmt.Errorf("invalid database text: %+v", s)
}

func Timestamptz(t time.Time) (pgtype.Timestamptz, error) {
	var timestamptz pgtype.Timestamptz
	err := timestamptz.Scan(t)
	return timestamptz, err
}

func Time(t pgtype.Timestamptz) (time.Time, error) {
	if !t.Valid {
		return time.Time{}, fmt.Errorf("invalid time %+v", t)
	}
	return t.Time, nil
}

func TimestamptzNullable(t *time.Time) (pgtype.Timestamptz, error) {
	var timestamptz pgtype.Timestamptz
	if t != nil {
		err := timestamptz.Scan(*t)
		return timestamptz, err
	}
	return timestamptz, nil
}

func TimeNullable(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
