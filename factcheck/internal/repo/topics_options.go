package repo

import (
	"strings"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

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

// OptionTopicDynamic represents topic-specific options
type OptionTopicDynamic func(*OptionsTopicDynamic)

type OptionsTopicDynamic struct {
	Options
	LikeID          string
	LikeMessageText string
	Statuses        []factcheck.StatusTopic
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

func (o OptionsTopicDynamic) ListDynamicParams(offset, limit int) postgres.ListTopicsDynamicParams {
	// Add wildcards for LIKE queries if not already present
	likeIDPattern := o.LikeID
	if likeIDPattern != "" && !strings.Contains(likeIDPattern, "%") {
		likeIDPattern = substring(likeIDPattern)
	}

	likeMessagePattern := o.LikeMessageText
	if likeMessagePattern != "" && !strings.Contains(likeMessagePattern, "%") {
		likeMessagePattern = substring(likeMessagePattern)
	}

	return postgres.ListTopicsDynamicParams{
		Column1: likeIDPattern,
		Column2: utils.MapSliceNoError(o.Statuses, utils.String[factcheck.StatusTopic, string]),
		Column3: likeMessagePattern,
		Column4: int32(limit),
		Column5: int32(offset),
	}
}
