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
	return utils.MapNoError(topics, ToTopic)
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
		Language:      TextNullable(m.Language),
		Status:        string(m.Status),
		TopicID:       topicID,
		Text:          m.Text,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
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
	return utils.Map(data, ToUserMessage)
}

func MessageV2Creator(m factcheck.MessageV2) (CreateMessageV2Params, error) {
	id, err := UUID(m.ID)
	if err != nil {
		return CreateMessageV2Params{}, err
	}
	// Could be nil
	topicID, _ := UUID(m.TopicID)
	createdAt, err := Timestamptz(m.CreatedAt)
	if err != nil {
		return CreateMessageV2Params{}, err
	}
	updatedAt, err := TimestamptzNullable(m.UpdatedAt)
	if err != nil {
		return CreateMessageV2Params{}, err
	}
	metadata, err := json.Marshal(m.Metadata)
	if err != nil {
		return CreateMessageV2Params{}, err
	}
	return CreateMessageV2Params{
		ID:        id,
		UserID:    m.UserID,
		TopicID:   topicID,
		TypeUser:  string(m.TypeUser),
		Type:      string(m.TypeMessage),
		Text:      m.Text,
		Language:  pgtype.Text{}, // MessageV2 doesn't have Language field
		Metadata:  metadata,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func ToMessageV2(data MessagesV2) (factcheck.MessageV2, error) {
	id, err := FromUUID(data.ID)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	createdAt, err := Time(data.CreatedAt)
	if err != nil {
		return factcheck.MessageV2{}, err
	}
	var metadata json.RawMessage
	if len(data.Metadata) > 0 {
		metadata = json.RawMessage(data.Metadata)
	}
	message := factcheck.MessageV2{
		ID:          id,
		UserID:      data.UserID,
		TypeUser:    factcheck.TypeUserMessage(data.TypeUser),
		TypeMessage: factcheck.TypeMessage(data.Type),
		Text:        data.Text,
		Metadata:    metadata,
		CreatedAt:   createdAt,
		UpdatedAt:   TimeNullable(data.UpdatedAt),
	}
	if data.TopicID.Valid {
		message.TopicID = data.TopicID.String()
	}
	if data.GroupID.Valid {
		message.GroupID = data.GroupID.String()
	}
	// MessageV2 doesn't have Language field
	return message, nil
}

func ToMessagesV2(data []MessagesV2) ([]factcheck.MessageV2, error) {
	return utils.Map(data, ToMessageV2)
}

func MessageV2GroupCreator(g factcheck.MessageGroup) (CreateMessageV2GroupParams, error) {
	id, err := UUID(g.ID)
	if err != nil {
		return CreateMessageV2GroupParams{}, err
	}
	topicID, err := UUID(g.TopicID)
	if err != nil {
		return CreateMessageV2GroupParams{}, err
	}
	createdAt, err := Timestamptz(g.CreatedAt)
	if err != nil {
		return CreateMessageV2GroupParams{}, err
	}
	updatedAt, err := TimestamptzNullable(g.UpdatedAt)
	if err != nil {
		return CreateMessageV2GroupParams{}, err
	}
	return CreateMessageV2GroupParams{
		ID:        id,
		TopicID:   topicID,
		Name:      TextNullable(g.Name),
		Text:      TextNullable(g.Text),
		TextSha1:  TextNullable(g.TextSHA1),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func ToMessageV2Group(data MessagesV2Group) (factcheck.MessageGroup, error) {
	id, err := FromUUID(data.ID)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	topicID, err := FromUUID(data.TopicID)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	createdAt, err := Time(data.CreatedAt)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	group := factcheck.MessageGroup{
		ID:        id,
		TopicID:   topicID,
		CreatedAt: createdAt,
		UpdatedAt: TimeNullable(data.UpdatedAt),
	}
	if data.Name.Valid {
		group.Name = data.Name.String
	}
	if data.Text.Valid {
		group.Text = data.Text.String
	}
	if data.TextSha1.Valid {
		group.TextSHA1 = data.TextSha1.String
	}
	return group, nil
}

func ToMessageV2Groups(data []MessagesV2Group) ([]factcheck.MessageGroup, error) {
	return utils.Map(data, ToMessageV2Group)
}

func UUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	return uuid, err
}

func UUIDs(ids []string) ([]pgtype.UUID, error) {
	return utils.Map(ids, UUID)
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

func Text[S ~string](s S) (pgtype.Text, error) {
	var text pgtype.Text
	err := text.Scan(string(s))
	return text, err
}

func TextNullable[S ~string](s S) pgtype.Text {
	var text pgtype.Text
	if s == "" {
		return text
	}
	err := text.Scan(string(s))
	if err != nil {
		text = pgtype.Text{}
		slog.Error("bad language text", "language", s, "err", err)
	}
	return text
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
