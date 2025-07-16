package repo

import (
	"github.com/kaogeek/line-fact-check/factcheck"
)

// OptionTopic represents topic-specific options
type OptionTopic func(*OptionsTopic)

// OptionsTopic contains topic operation configuration
type OptionsTopic struct {
	TxOptions
	// Topic-specific filters
	LikeID          string
	LikeMessageText string
	Status          factcheck.StatusTopic
}

// WithTopicLikeID sets the topic ID pattern filter
func WithTopicLikeID(id string) OptionTopic {
	return func(opts *OptionsTopic) { opts.LikeID = id }
}

// WithTopicLikeMessageText sets the message text pattern filter
func WithTopicLikeMessageText(text string) OptionTopic {
	return func(opts *OptionsTopic) { opts.LikeMessageText = text }
}

// WithTopicStatus sets the topic status filter
func WithTopicStatus(status factcheck.StatusTopic) OptionTopic {
	return func(opts *OptionsTopic) { opts.Status = status }
}

// WithTopicTx sets the transaction for topic operations
func WithTopicTx(tx Tx) OptionTopic {
	return func(opts *OptionsTopic) { opts.Tx = tx }
}
