// Package service defines entrypoints for complex business use cases
// If your logic is just getting/listing/deleting stuff, do it directly in the HTTP handler.
package service

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Service interface {
	// Submit handles new message submission by creating the message and assigning it to a group.
	// Submit returns message created, message group assigned to the new message, and topic (if any)
	//
	// Caller could call this Submit, and on success gets all the messages from users for replies.
	Submit(ctx context.Context, user factcheck.UserInfo, text string, topicID string) (factcheck.MessageV2, factcheck.MessageGroup, *factcheck.Topic, error)

	// ResolveTopic resolves topic and returns list of messages associated with the topic.
	ResolveTopic(ctx context.Context, user factcheck.UserInfo, topicID string, answer string) (factcheck.Answer, factcheck.Topic, []factcheck.MessageV2, error)
}

func New(repo repo.Repository) ServiceFactcheck {
	return ServiceFactcheck{repo: repo}
}

type ServiceFactcheck struct {
	repo repo.Repository
}
