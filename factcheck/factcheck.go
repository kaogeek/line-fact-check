// Package factcheck defines shared symbols for the entire project.
// It is very business-centric and agnostic to tech stack
package factcheck

import "time"

type (
	TopicResult       string
	StatusTopic       string
	StatusTopicResult string
	TypeMessage       string
)

const (
	StatusTopicPending  StatusTopic = "TOPIC_PENDING"
	StatusTopicResolved StatusTopic = "TOPIC_RESOLVED"
	TypeMessageText     TypeMessage = "TYPE_TEXT"
)

type Topic struct {
	ID           string
	Name         string
	Status       StatusTopic
	Result       string
	ResultStatus StatusTopicResult // TODO: wat?
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type Message struct {
	ID        string
	TopicID   string
	Text      string
	Type      TypeMessage
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type UserMessage[T any] struct {
	RepliedAt *time.Time
	MessageID string
	Metadata  T
	CreatedAt time.Time
	UpdatedAt *time.Time
}

func Bar() {}
