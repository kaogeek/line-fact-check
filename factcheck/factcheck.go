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
	StatusTopicPending  StatusTopic = "TOPIC_PENDING"  // topic automatically created, no answer yet
	StatusTopicResolved StatusTopic = "TOPIC_RESOLVED" // topic resolved by human admins

	StatusTopicResultNone        StatusTopicResult = "TOPIC_RESULT_NONE"       // no prior answer
	StatusTopicResultAnswered    StatusTopicResult = "TOPIC_RESULT_ANSWERED"   // answered at least once
	StatusTopicResultChanllenged StatusTopicResult = "TOPIC_RESULT_CHALLENGED" // the last answer was challenged by the public

	TypeMessageText TypeMessage = "TYPE_TEXT"
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

func (s StatusTopic) IsValid() bool {
	switch s {
	case "":
		panic("got empty topic status")
	case
		StatusTopicPending,
		StatusTopicResolved:
		return true
	}
	return false
}

func (s StatusTopicResult) IsValid() bool {
	switch s {
	case "":
		panic("got empty topic result status")
	case
		StatusTopicResultNone,
		StatusTopicResultAnswered,
		StatusTopicResultChanllenged:
		return true
	}
	return false
}

func (t TypeMessage) IsValid() bool {
	switch t {
	case "":
		panic("got empty message type")
	case TypeMessageText:
		return true
	}
	return false
}
