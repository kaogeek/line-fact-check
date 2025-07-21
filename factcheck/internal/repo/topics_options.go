package repo

import (
	"github.com/kaogeek/line-fact-check/factcheck"
)

// OptionTopic represents topic-specific options
type OptionTopic func(*OptionsTopic)

type OptionsTopic struct {
	Options
	LikeID          string
	LikeMessageText string
	Statuses        []factcheck.StatusTopic
}

func TopicLikeID(id string) OptionTopic {
	return func(opts *OptionsTopic) {
		opts.LikeID = substringAuto(id)
	}
}

func TopicLikeMessageText(text string) OptionTopic {
	return func(opts *OptionsTopic) {
		opts.LikeMessageText = substringAuto(text)
	}
}

func TopicInStatuses(statuses []factcheck.StatusTopic) OptionTopic {
	return func(opts *OptionsTopic) {
		opts.Statuses = statuses
	}
}
