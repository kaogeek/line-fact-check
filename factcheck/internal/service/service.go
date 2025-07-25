// Package service defines entrypoints for complex business use cases
// If your logic is just getting/listing/deleting stuff, do it directly in the HTTP handler.
package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Service interface {
	AdminSubmit(context.Context, factcheck.MessageV2) (factcheck.MessageV2, error)
}

func New(repo repo.Repository) Service {
	panic("not implemented")
}

type service struct {
	repo repo.Repository
}

func (s *service) AdminSubmit(
	ctx context.Context,
	userID string,
	text string,
	topicID string,
) (
	factcheck.MessageV2,
	factcheck.MessageV2Group,
	error,
) {
	textSHA1, err := factcheck.SHA1Base64(text)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageV2Group{}, err
	}
	tx, err := s.repo.BeginTx(ctx, repo.Serializable)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageV2Group{}, err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err == nil {
			return
		}
		slog.ErrorContext(ctx, "error rolling back AdminSubmit", "err", err)
	}()

	now := time.Now()
	withTx := repo.WithTx(tx)

	if topicID != "" {
		exists, err := s.repo.Topics.Exists(ctx, topicID, withTx)
		if err != nil {
			return factcheck.MessageV2{}, factcheck.MessageV2Group{}, fmt.Errorf("error getting topic '%s' for a new message", topicID)
		}
		if !exists {
			return factcheck.MessageV2{}, factcheck.MessageV2Group{}, fmt.Errorf("non existent topic '%s'", topicID)
		}
	}
	group, err := s.repo.MessagesV2Groups.GetBySHA1(ctx, textSHA1, withTx)
	if err != nil {
		if !repo.IsNotFound(err) {
			return factcheck.MessageV2{}, factcheck.MessageV2Group{}, fmt.Errorf("error finding group based on sha1 hash '%s'", textSHA1)
		}
		group = factcheck.MessageV2Group{
			ID:        uuid.NewString(),
			TopicID:   topicID,
			Text:      text,
			TextSHA1:  textSHA1,
			CreatedAt: now,
		}
		slog.Info("creating new group",
			"uuid", group.ID,
			"name", group.Name,
			"text_sha1", group.SHA1,
		)
		group, err = s.repo.MessagesV2Groups.Create(ctx, group, withTx)
		if err != nil {
			return factcheck.MessageV2{}, factcheck.MessageV2Group{}, fmt.Errorf("error pre-creating group %s", textSHA1)
		}
	}
	m := factcheck.MessageV2{
		ID:          uuid.NewString(),
		TopicID:     topicID,
		UserID:      userID,
		GroupID:     group.ID,
		TypeUser:    factcheck.TypeUserMessageAdmin,
		TypeMessage: factcheck.TypeMessageText,
		Text:        text,
		CreatedAt:   now,
	}
	created, err := s.repo.MessagesV2.Create(ctx, m, withTx)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageV2Group{}, fmt.Errorf("error creating message: %w", err)
	}
	return created, group, nil
}
