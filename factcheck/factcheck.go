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

type shared struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Topic struct {
	ID           string
	Name         string
	Status       StatusTopic
	Result       string
	ResultStatus StatusTopicResult
	shared
}

type Message struct {
	ID      string
	TopicID string
	Text    string
	Type    TypeMessage
	shared
}

type UserMessage[T any] struct {
	RepliedAt *time.Time
	MessageID string
	Metadata  T
	shared
}

func Bar() {}
