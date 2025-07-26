// Package factcheck defines shared symbols for the entire project.
// It is very business-centric and agnostic to tech stack
package factcheck

import (
	"bytes"
	"crypto/sha1" //nolint:gosec
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// nolint:gosec
// The checksum is only used to match text,
// not for cryptographic purpose
var checksum = sha1.New()

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

type MessageV2 struct {
	ID          string          `json:"id"`
	TopicID     string          `json:"topic_id"`
	UserID      string          `json:"user_id"`
	GroupID     string          `json:"group_id"`
	TypeUser    TypeUserMessage `json:"type_user"`
	TypeMessage TypeMessage     `json:"type"`
	Text        string          `json:"text"`
	Metadata    json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at"`
}

type MessageGroup struct {
	ID        string     `json:"id"`
	TopicID   string     `json:"topic_id"`
	Name      string     `json:"name"`
	Text      string     `json:"text"`
	TextSHA1  string     `json:"text_sha1"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Answer struct {
	ID        string     `json:"id"`
	TopicID   string     `json:"topic_id"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
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

func (m MessageV2) SHA1() (string, error) {
	if m.GroupID != "" {
		return "", fmt.Errorf("message was already assigned to group %s", m.GroupID)
	}
	if m.Text == "" {
		return "", errors.New("message has empty text")
	}
	return SHA1Base64(m.Text)
}

func (g MessageGroup) SHA1() (string, error) {
	return SHA1Base64(g.Text)
}

func SHA1Base64(s string) (string, error) {
	s = strings.TrimSpace(s)
	hash := sha1sum([]byte(s))
	buf := bytes.NewBuffer(nil)
	_, err := base64.
		NewEncoder(base64.StdEncoding, buf).
		Write(hash)
	if err != nil {
		return "", fmt.Errorf("base64 error: %w", err)
	}
	b64 := buf.String()
	return strings.ToLower(b64), nil
}

func sha1sum(b []byte) []byte { return checksum.Sum(b) }
