package service

import (
	"context"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func (s ServiceFactcheck) Resolve(
	ctx context.Context,
	user factcheck.UserInfo,
	topicID string,
	answerText string,
) (
	factcheck.Answer,
	factcheck.Topic,
	[]factcheck.MessageV2,
	error,
) {
	tx, err := s.repo.BeginTx(ctx, repo.RepeatableRead)
	if err != nil {
		return factcheck.Answer{}, factcheck.Topic{}, nil, err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err == nil {
			return
		}
		slog.ErrorContext(ctx, "error rolling back after failure to resolve topic",
			"topic_id", topicID,
			"user", user,
		)
	}()

	withTx := repo.WithTx(tx)
	answer := factcheck.Answer{
		ID:        utils.NewID().String(),
		UserID:    user.UserID,
		TopicID:   topicID,
		Text:      answerText,
		CreatedAt: utils.TimeNow(),
	}
	answer, err = s.repo.Answers.Create(ctx, answer, withTx)
	if err != nil {
		return factcheck.Answer{}, factcheck.Topic{}, nil, err
	}
	resolved, err := s.repo.Topics.Resolve(ctx, topicID, answerText, withTx)
	if err != nil {
		return factcheck.Answer{}, factcheck.Topic{}, nil, err
	}
	messages, err := s.repo.MessagesV2.ListByTopic(ctx, topicID, withTx)
	if err != nil {
		return factcheck.Answer{}, factcheck.Topic{}, nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return factcheck.Answer{}, factcheck.Topic{}, nil, err
	}
	return answer, resolved, messages, nil
}
