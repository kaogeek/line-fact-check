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
	// Submit handles new message submission by creating the message and assigning it to a group.
	// Submit returns message created, message group assigned to the new message, and topic (if any)
	//
	// Caller could call this Submit, and on success gets all the messages from users for replies.
	Submit(ctx context.Context, user factcheck.UserInfo, text string, topicID string) (factcheck.MessageV2, factcheck.MessageGroup, *factcheck.Topic, error)

	// ResolveTopic resolves topic and returns list of messages associated with the topic.
	ResolveTopic(ctx context.Context, user factcheck.UserInfo, topicID string, answer string) (factcheck.Answer, factcheck.Topic, []factcheck.MessageV2, error)
}

func New(repo repo.Repository) Service { return &service{repo: repo} }

type service struct{ repo repo.Repository }

func (s *service) ResolveTopic(
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

func (s *service) Submit(
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
	textSHA1 := factcheck.SHA1(text)
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
	if err == nil && !utils.Empty(topicID, group.ID) && topicID != group.ID {
		// TODO: what to do?
		// Mismatch topicID
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("mismatch topic '%s': found group %s (%s) has topic '%s'", topicID, group.ID, textSHA1, group.TopicID)
	}
	if err != nil {
		if !repo.IsNotFound(err) {
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error finding group based on sha1 hash '%s'", textSHA1)
		}

		// If not found, we'll create a new group for it.
		// But the group will not have topicID - to be assigned topic by admin
		group = factcheck.MessageGroup{
			ID:        utils.NewID().String(),
			Status:    factcheck.StatusMGroupPending,
			Text:      text,
			TextSHA1:  textSHA1,
			CreatedAt: now,
		}
		slog.InfoContext(ctx, "creating new group without topic",
			"gid", group.ID,
			"name", group.Name,
			"text_sha1", group.SHA1,
		)
		group, err = s.repo.MessageGroups.Create(ctx, group, withTx)
		if err != nil {
			return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error pre-creating group %s", textSHA1)
		}
	}

	meta := factcheck.Metadata[factcheck.UserInfo]{
		Type: factcheck.TypeMetadataUserInfo,
		Data: user,
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return factcheck.MessageV2{}, factcheck.MessageGroup{}, nil, fmt.Errorf("error pre-creating group %s", textSHA1)
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
