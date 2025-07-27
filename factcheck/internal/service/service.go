// Package service defines entrypoints for complex business use cases
// If your logic is just getting/listing/deleting stuff, do it directly in the HTTP handler.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

type Service interface {
	// Submit handles new message submission
	// by creating the message and assigning it to a group
	Submit(ctx context.Context, user factcheck.UserInfo, text string, topicID string) (factcheck.MessageV2, factcheck.MessageGroup, error)
}

func New(repo repo.Repository) Service { return &service{repo: repo} }

type service struct{ repo repo.Repository }

func (s *service) Submit(
	ctx context.Context,
	user factcheck.UserInfo,
	text string,
	topicID string, // If empty, the new message will be topic-less
) (
	factcheck.MessageV2,
	factcheck.MessageGroup,
	error,
) {
	textSHA1, err := factcheck.SHA1Base64(text)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, err
	}
	tx, err := s.repo.BeginTx(ctx, repo.Serializable)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, err
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
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, fmt.Errorf("error getting topic '%s' for a new message", topicID)
		}
		if !exists {
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, fmt.Errorf("non existent topic '%s'", topicID)
		}
	}
	group, err := s.repo.MessageGroups.GetBySHA1(ctx, textSHA1, withTx)
	if err != nil {
		if !repo.IsNotFound(err) {
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, fmt.Errorf("error finding group based on sha1 hash '%s'", textSHA1)
		}

		// If not found, we'll create a new group for it.
		// But the group will not have topicID - to be assigned topic by admin
		group = factcheck.MessageGroup{
			ID:        utils.NewID().String(),
			Text:      text,
			TextSHA1:  textSHA1,
			CreatedAt: now,
		}
		slog.Info("creating new group without topic",
			"gid", group.ID,
			"name", group.Name,
			"text_sha1", group.SHA1,
		)
		group, err = s.repo.MessageGroups.Create(ctx, group, withTx)
		if err != nil {
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, fmt.Errorf("error pre-creating group %s", textSHA1)
		}
	}

	meta := factcheck.Metadata[factcheck.UserInfo]{
		Type: factcheck.TypeMetadataUserInfo,
		Data: user,
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, fmt.Errorf("error pre-creating group %s", textSHA1)
	}
	message := factcheck.MessageV2{
		ID:          utils.NewID().String(),
		TopicID:     group.TopicID,
		UserID:      user.UserID,
		GroupID:     group.ID,
		TypeUser:    factcheck.TypeUserMessageAdmin,
		TypeMessage: factcheck.TypeMessageText,
		Text:        text,
		Metadata:    metaJSON,
		CreatedAt:   now,
	}

	created, err := s.repo.MessagesV2.Create(ctx, message, withTx)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, fmt.Errorf("error creating message: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("error committing admin submission",
			"err", err,
			"mid", message.ID,
			"gid", group.ID,
			"sha1", textSHA1,
		)
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, fmt.Errorf("error committing message: %w", err)
	}
	return created, group, nil
}
