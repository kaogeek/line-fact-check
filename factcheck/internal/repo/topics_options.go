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
