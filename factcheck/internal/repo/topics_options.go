package repo

import (
	"strings"

	"github.com/kaogeek/line-fact-check/factcheck"
)

// OptionTopicDynamic represents topic-specific options
type OptionTopicDynamic func(*OptionsTopicDynamic)

type OptionsTopicDynamic struct {
	Options
	LikeID          string
	LikeMessageText string
	Statuses        []factcheck.StatusTopic
}

func substringAuto(pattern string) string {
	if pattern != "" && !strings.Contains(pattern, "%") {
		pattern = substring(pattern)
	}
	return pattern
}

func WithTopicDynamicLikeID(id string) OptionTopicDynamic {
	return func(opts *OptionsTopicDynamic) {
		opts.LikeID = substringAuto(id)
	}
}

func WithTopicDynamicLikeMessageText(text string) OptionTopicDynamic {
	return func(opts *OptionsTopicDynamic) {
		opts.LikeMessageText = substringAuto(text)
	}
}

func WithTopicDynamicStatuses(statuses []factcheck.StatusTopic) OptionTopicDynamic {
	return func(opts *OptionsTopicDynamic) {
		opts.Statuses = statuses
	}
}

// WithTopicTx sets the transaction for topic operations
func WithTopicDynamicTx(tx Tx) OptionTopicDynamic {
	return func(opts *OptionsTopicDynamic) { opts.tx = tx }
}

func (o OptionsTopicDynamic) Clone() []OptionTopic {
	opts := []OptionTopic{}
	if o.LikeID != "" {
		opts = append(opts, WithTopicLikeID(o.LikeID))
	}
	if o.LikeMessageText != "" {
		opts = append(opts, WithTopicLikeMessageText(o.LikeMessageText))
	}
	if len(o.Statuses) > 0 {
		opts = append(opts, WithTopicStatus(o.Statuses[0]))
	}
	if o.tx != nil {
		opts = append(opts, WithTopicTx(o.tx))
	}
	return opts
}

// ------------------------------------------------------------
// OptionTopic to be deprecated
// ------------------------------------------------------------
// OptionTopic represents topic-specific options
type OptionTopic func(*OptionsTopic)

// OptionsTopic contains topic operation configuration
type OptionsTopic struct {
	Options
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
	return func(opts *OptionsTopic) { opts.tx = tx }
}
