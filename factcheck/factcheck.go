// Package factcheck defines shared symbols for the entire project.
// It is very business-centric and agnostic to tech stack
package factcheck

import (
	"crypto/sha1" //nolint:gosec
	"encoding/hex"
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
	TopicResult string
	StatusTopic string
	StatusUser  string
	TypeMessage string
	TypeUser    string
	Language    string

	StatusMGroup string
)

const (

	// StatusTopicPending is used when topic is automatically created.
	// Note that since only admins can create topics,
	// topics do not need approval or rejection like message groups
	StatusTopicPending  StatusTopic = "TOPIC_PENDING"
	StatusTopicResolved StatusTopic = "TOPIC_RESOLVED"

	TypeMessageText TypeMessage = "MSG_TEXT"

	TypeUserMessageLINEChat      TypeUser = "USER_CHAT"
	TypeUserMessageLINEGroupChat TypeUser = "USER_GROUPCHAT"
	TypeUserMessageAdmin         TypeUser = "USER_ADMIN"

	StatusMGroupPending  StatusMGroup = "PENDING"
	StatusMGroupApproved StatusMGroup = "APPROVED"
	StatusMGroupRejected StatusMGroup = "REJECTED"

	LanguageEnglish Language = "en"
	LanguageThai    Language = "th"
)

type Topic struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Status      StatusTopic `json:"status"`
	Result      string      `json:"result"`
	RepliedAt   *time.Time  `json:"replied_at"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   *time.Time  `json:"updated_at"`
}

type MessageV2 struct {
	ID          string          `json:"id"`
	GroupID     string          `json:"group_id"`
	TopicID     string          `json:"topic_id"`
	UserID      string          `json:"user_id"`
	TypeUser    TypeUser        `json:"type_user"`
	TypeMessage TypeMessage     `json:"type"`
	Text        string          `json:"text"`
	Metadata    json.RawMessage `json:"metadata"`
	RepliedAt   *time.Time      `json:"replied_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   *time.Time      `json:"updated_at"`
}

type MessageGroup struct {
	ID        string       `json:"id"`
	Status    StatusMGroup `json:"status"`
	TopicID   string       `json:"topic_id"`
	Name      string       `json:"name"`
	Text      string       `json:"text"`
	TextSHA1  string       `json:"text_sha1"`
	Language  Language     `json:"language"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
}

type Answer struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TopicID   string    `json:"topic_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
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

func (t TypeUser) IsValid() bool {
	switch t {
	case
		TypeUserMessageAdmin,
		TypeUserMessageLINEChat,
		TypeUserMessageLINEGroupChat:
		return true
	}
	return false
}

func (t TypeMessage) IsValid() bool {
	return t == TypeMessageText
}

func (t Topic) Validate() error {
	switch t.Status {
	case StatusTopicResolved:
		if t.Result == "" {
			return fmt.Errorf("unexpected empty topic result of resolved topic '%s'", t.ID)
		}
	case StatusTopicPending:
		if t.Result != "" {
			return fmt.Errorf("unexpected non-empty result of pending topic '%s': '%s'", t.ID, t.Result)
		}
	}
	return nil
}

func (m MessageV2) Validate() error {
	if m.Text == "" {
		return errors.New("empty text")
	}
	if !m.TypeUser.IsValid() {
		return fmt.Errorf("invalid typeUser '%s'", m.TypeUser)
	}
	if !m.TypeMessage.IsValid() {
		return fmt.Errorf("invalid typeMessage '%s'", m.TypeMessage)
	}
	return nil
}

func (m MessageV2) SHA1() (string, error) {
	if m.GroupID != "" {
		return "", fmt.Errorf("message was already assigned to group %s", m.GroupID)
	}
	if m.Text == "" {
		return "", errors.New("message has empty text")
	}
	return SHA1(m.Text), nil
}

func (g MessageGroup) SHA1() (string, error) {
	return SHA1(g.Text), nil
}

func SHA1(s string) string {
	hash := sha1sum([]byte(strings.TrimSpace(s)))
	return hex.EncodeToString(hash)
}

func sha1sum(b []byte) []byte { return checksum.Sum(b) }
