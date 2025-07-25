// Package factcheck defines shared symbols for the entire project.
// It is very business-centric and agnostic to tech stack
package factcheck

import (
	"encoding/json"
	"time"
)

type (
	TopicResult       string
	StatusTopic       string
	StatusTopicResult string
	StatusMessage     string
	TypeMessage       string
	TypeUserMessage   string
	Language          string
)

const (
	StatusTopicPending  StatusTopic = "TOPIC_PENDING"  // topic automatically created, no answer yet
	StatusTopicResolved StatusTopic = "TOPIC_RESOLVED" // topic resolved by human admins

	StatusTopicResultNone        StatusTopicResult = "TOPIC_RESULT_NONE"       // no prior answer
	StatusTopicResultAnswered    StatusTopicResult = "TOPIC_RESULT_ANSWERED"   // answered at least once
	StatusTopicResultChanllenged StatusTopicResult = "TOPIC_RESULT_CHALLENGED" // the last answer was challenged by the public

	StatusMessageSubmitted      StatusMessage = "MSG_SUBMITTED"
	StatusMessageTopicSubmitted StatusMessage = "MSG_TOPIC_SUBMITTED"
	StatusMessageTopicAssigned  StatusMessage = "MSG_TOPIC_ASSIGNED"

	TypeMessageText TypeMessage = "TYPE_TEXT"

	TypeUserMessageLINEChat      TypeUserMessage = "CHAT"
	TypeUserMessageLINEGroupChat TypeUserMessage = "GROUPCHAT"
	TypeUserMessageAdmin         TypeUserMessage = "ADMIN"

	LanguageEnglish Language = "en"
	LanguageThai    Language = "th"
)

type Topic struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Status       StatusTopic       `json:"status"`
	Result       string            `json:"result"`
	ResultStatus StatusTopicResult `json:"result_status"`
	RepliedAt    *time.Time        `json:"replied_at"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    *time.Time        `json:"updated_at"`
}

// Message is to be deprecated and replaced with MessageV2
type Message struct {
	ID            string        `json:"id"`
	UserMessageID string        `json:"user_message_id"`
	Type          TypeMessage   `json:"type"`
	Status        StatusMessage `json:"status"`
	TopicID       string        `json:"topic_id"`
	Text          string        `json:"text"`
	Language      Language      `json:"language"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     *time.Time    `json:"updated_at"`
}

// UserMessage is to be deprecated and replaced with MessageV2
type UserMessage struct {
	ID        string          `json:"id"`
	Type      TypeUserMessage `json:"type"`
	RepliedAt *time.Time      `json:"replied_at"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}

type MessageV2 struct {
	ID        string          `json:"id"`
	TopicID   string          `json:"topic_id"`
	UserID    string          `json:"user_id"`
	Type      TypeUserMessage `json:"type"`
	Text      string          `json:"text"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}

func (s StatusTopic) IsValid() bool {
	switch s {
	case
		StatusTopicPending,
		StatusTopicResolved:
		return true
	}
	return false
}

func (s StatusTopicResult) IsValid() bool {
	switch s {
	case
		StatusTopicResultNone,
		StatusTopicResultAnswered,
		StatusTopicResultChanllenged:
		return true
	}
	return false
}

func (t TypeMessage) IsValid() bool {
	return t == TypeMessageText
}
