package repo

import (
	"github.com/kaogeek/line-fact-check/factcheck"
)

// TopicOption represents topic-specific options
type TopicOption func(*TopicOptions)

// TopicOptions contains topic operation configuration
type TopicOptions struct {
	TxOptions
	// Topic-specific filters
	LikeID          string
	LikeMessageText string
	Status          factcheck.StatusTopic
}

// WithTopicLikeID sets the topic ID pattern filter
func WithTopicLikeID(id string) TopicOption {
	return func(opts *TopicOptions) { opts.LikeID = id }
}

// WithTopicLikeMessageText sets the message text pattern filter
func WithTopicLikeMessageText(text string) TopicOption {
	return func(opts *TopicOptions) { opts.LikeMessageText = text }
}

// WithTopicStatus sets the topic status filter
func WithTopicStatus(status factcheck.StatusTopic) TopicOption {
	return func(opts *TopicOptions) { opts.Status = status }
}

// WithTopicTx sets the transaction for topic operations
func WithTopicTx(tx Tx) TopicOption {
	return func(opts *TopicOptions) { opts.Tx = tx }
}
