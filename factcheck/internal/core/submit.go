package core

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

func (s ServiceFactcheck) Submit(
	ctx context.Context,
	user factcheck.UserInfo,
	text string,
	topicID string, // Users can submit with topic_id, but this will be pending approval for inclusion into topic
) (
	factcheck.MessageV2,
	factcheck.MessageGroup,
	*factcheck.Topic,
	error,
) {
	meta := factcheck.Metadata[factcheck.UserInfo]{
		Type: factcheck.TypeMetadataUserInfo,
		Data: user,
	}
	textSHA1 := factcheck.SHA1(text)
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error creating metadata %s: %w", textSHA1, err)
	}
	tx, err := s.repo.BeginTx(ctx, repo.RepeatableRead)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, err
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

	var topic *factcheck.Topic
	if topicID != "" {
		topicDB, err := s.repo.Topics.GetByID(ctx, topicID, withTx)
		if err != nil {
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error getting topic '%s' for a new message: %w", topicID, err)
		}
		err = topicDB.Validate()
		if err != nil {
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error validating topic '%s' for a new message: %w", topicID, err)
		}
		topic = &topicDB
	}

	group, err := s.repo.MessageGroups.GetBySHA1(ctx, textSHA1, withTx)
	if err != nil {
		if !repo.IsNotFound(err) {
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error finding group based on sha1 hash '%s': %s", textSHA1, err)
		}
		// If not found, we'll create a new group for it.
		// But the group will not have topicID - to be assigned topic by admin
		group := factcheck.MessageGroup{
			ID:        utils.NewID().String(),
			Status:    factcheck.StatusMGroupPending,
			Text:      text,
			TextSHA1:  textSHA1,
			CreatedAt: now,
		}
		slog.InfoContext(ctx, "creating new group without topic",
			"gid", group.ID,
			"name", group.Name,
			"sha1", group.SHA1,
		)
		group, err = s.repo.MessageGroups.Create(ctx, group, withTx)
		if err != nil {
			slog.Error("error pre-creating group",
				"gid", group.ID,
				"sha1", textSHA1,
				"err", err,
			)
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error pre-creating group %s: %w", textSHA1, err)
		}
	}
	if !utils.Empty(topicID, group.ID) && topicID != group.ID {
		// TODO: what to do?
		// Mismatch topicID
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("mismatch topic '%s': found group %s (%s) has topic '%s'", topicID, group.ID, textSHA1, group.TopicID)
	}

	message := factcheck.MessageV2{
		ID:          utils.NewID().String(),
		TopicID:     group.TopicID,
		GroupID:     group.ID,
		UserID:      user.UserID,
		TypeUser:    user.UserType,
		TypeMessage: factcheck.TypeMessageText,
		Text:        text,
		Metadata:    metaJSON,
		CreatedAt:   now,
	}

	created, err := s.repo.MessagesV2.Create(ctx, message, withTx)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error creating message: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error committing admin submission",
			"err", err,
			"mid", message.ID,
			"gid", group.ID,
			"sha1", textSHA1,
		)
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error committing message: %w", err)
	}
	return created, group, topic, nil
}
