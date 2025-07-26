package factcheck

import (
	"encoding/json"
	"time"
)

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
	Type      TypeUser        `json:"type"`
	RepliedAt *time.Time      `json:"replied_at"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}
