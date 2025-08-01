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
	return CreateTopicParams{
		ID:          id,
		Name:        topic.Name,
		Description: topic.Description,
		Status:      string(topic.Status),
		Result:      result,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
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
	if data.CreatedAt.Valid {
		topic.CreatedAt = data.CreatedAt.Time
	}
	if data.UpdatedAt.Valid {
		topic.UpdatedAt = &data.UpdatedAt.Time
	}
	return topic
}

func MessageV2Creator(m factcheck.MessageV2) (CreateMessageV2Params, error) {
	id, err := UUID(m.ID)
	if err != nil {
		return CreateMessageV2Params{}, err
	}
	// Could be nil
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
		TopicID:   UUIDNullable(m.TopicID),
		GroupID:   UUIDNullable(m.GroupID),
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
		TypeUser:    factcheck.TypeUser(data.TypeUser),
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

func MessageGroupCreator(g factcheck.MessageGroup) (CreateMessageGroupParams, error) {
	id, err := UUID(g.ID)
	if err != nil {
		return CreateMessageGroupParams{}, err
	}
	createdAt, err := Timestamptz(g.CreatedAt)
	if err != nil {
		return CreateMessageGroupParams{}, err
	}
	updatedAt, err := TimestamptzNullable(g.UpdatedAt)
	if err != nil {
		return CreateMessageGroupParams{}, err
	}
	language, err := Text(g.Language)
	if err != nil {
		return CreateMessageGroupParams{}, err
	}
	topicID := UUIDNullable(g.TopicID)

	return CreateMessageGroupParams{
		ID:        id,
		TopicID:   topicID,
		Name:      g.Name,
		Text:      g.Text,
		TextSha1:  g.TextSHA1,
		Language:  language,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func ToMessageGroup(data MessageGroup) (factcheck.MessageGroup, error) {
	id, err := FromUUID(data.ID)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	var topicID string
	if data.TopicID.Valid {
		topicID, err = FromUUID(data.TopicID)
		if err != nil {
			return factcheck.MessageGroup{}, err
		}
	}
	createdAt, err := Time(data.CreatedAt)
	if err != nil {
		return factcheck.MessageGroup{}, err
	}
	group := factcheck.MessageGroup{
		ID:        id,
		Name:      data.Name,
		Text:      data.Text,
		TextSHA1:  data.TextSha1,
		TopicID:   topicID,
		CreatedAt: createdAt,
		UpdatedAt: TimeNullable(data.UpdatedAt),
	}
	return group, nil
}

func ToMessageGroups(data []MessageGroup) ([]factcheck.MessageGroup, error) {
	return utils.Map(data, ToMessageGroup)
}

func UUID(id string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(id)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return uuid, nil
}

func UUIDs(ids []string) ([]pgtype.UUID, error) {
	return utils.Map(ids, UUID)
}

func UUIDNullable(id string) pgtype.UUID {
	uuid, _ := UUID(id) //nolint
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
		slog.Error("bad language text", "language", s, "err", err) //nolint:noctx
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

func AnswerCreator(a factcheck.Answer) (CreateAnswerParams, error) {
	id, err := UUID(a.ID)
	if err != nil {
		return CreateAnswerParams{}, err
	}
	topicID, err := UUID(a.TopicID)
	if err != nil {
		return CreateAnswerParams{}, err
	}
	createdAt, err := Timestamptz(a.CreatedAt)
	if err != nil {
		return CreateAnswerParams{}, err
	}
	return CreateAnswerParams{
		ID:        id,
		TopicID:   topicID,
		Text:      a.Text,
		CreatedAt: createdAt,
	}, nil
}

func ToAnswer(data Answer) (factcheck.Answer, error) {
	id, err := FromUUID(data.ID)
	if err != nil {
		return factcheck.Answer{}, err
	}
	topicID, err := FromUUID(data.TopicID)
	if err != nil {
		return factcheck.Answer{}, err
	}
	createdAt, err := Time(data.CreatedAt)
	if err != nil {
		return factcheck.Answer{}, err
	}
	return factcheck.Answer{
		ID:        id,
		TopicID:   topicID,
		Text:      data.Text,
		CreatedAt: createdAt,
	}, nil
}

func ToAnswers(data []Answer) ([]factcheck.Answer, error) {
	return utils.Map(data, ToAnswer)
}
