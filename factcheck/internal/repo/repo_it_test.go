//go:build integration_test
// +build integration_test

package repo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestRepository_TopicFiltering(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()

	ctx := context.Background()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create topics
	topic1 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440001",
		Name:         "Topic 1 - COVID",
		Description:  "COVID-19 related news",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	topic2 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440002",
		Name:         "Topic 2 - Politics",
		Description:  "Political news and updates",
		Status:       factcheck.StatusTopicResolved,
		Result:       "Verified as true",
		ResultStatus: factcheck.StatusTopicResultAnswered,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	topic3 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440003",
		Name:         "Topic 3 - Technology",
		Description:  "Technology news and updates",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	// Create topics in database
	createdTopic1, err := app.Repository.Topics.Create(ctx, topic1)
	if err != nil {
		t.Fatalf("Failed to create topic1: %v", err)
	}

	createdTopic2, err := app.Repository.Topics.Create(ctx, topic2)
	if err != nil {
		t.Fatalf("Failed to create topic2: %v", err)
	}

	createdTopic3, err := app.Repository.Topics.Create(ctx, topic3)
	if err != nil {
		t.Fatalf("Failed to create topic3: %v", err)
	}

	// Create user messages first
	userMessage1 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440001",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-1"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage2 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440002",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-2"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage3 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440003",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-3"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage4 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440004",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-4"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	// Create user messages in database
	createdUserMessage1, err := app.Repository.UserMessages.Create(ctx, userMessage1)
	if err != nil {
		t.Fatalf("Failed to create userMessage1: %v", err)
	}

	createdUserMessage2, err := app.Repository.UserMessages.Create(ctx, userMessage2)
	if err != nil {
		t.Fatalf("Failed to create userMessage2: %v", err)
	}

	createdUserMessage3, err := app.Repository.UserMessages.Create(ctx, userMessage3)
	if err != nil {
		t.Fatalf("Failed to create userMessage3: %v", err)
	}

	createdUserMessage4, err := app.Repository.UserMessages.Create(ctx, userMessage4)
	if err != nil {
		t.Fatalf("Failed to create userMessage4: %v", err)
	}

	// Create messages
	message1 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440001",
		UserMessageID: createdUserMessage1.ID,
		TopicID:       createdTopic1.ID,
		Text:          "COVID-19 vaccine is effective against new variants",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message2 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440002",
		UserMessageID: createdUserMessage2.ID,
		TopicID:       createdTopic1.ID,
		Text:          "COVID-19 cases are increasing in winter",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message3 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440003",
		UserMessageID: createdUserMessage3.ID,
		TopicID:       createdTopic2.ID,
		Text:          "Election results show clear victory",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message4 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440004",
		UserMessageID: createdUserMessage4.ID,
		TopicID:       createdTopic3.ID,
		Text:          "New AI technology breakthrough",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	// Create messages in database
	// These messages are used indirectly through topic filtering tests
	_, err = app.Repository.Messages.Create(ctx, message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message3)
	if err != nil {
		t.Fatalf("Failed to create message3: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message4)
	if err != nil {
		t.Fatalf("Failed to create message4: %v", err)
	}

	t.Run("ListInIDs", func(t *testing.T) {
		// Test filtering by specific IDs
		ids := []string{createdTopic1.ID, createdTopic3.ID}
		topics, err := app.Repository.Topics.ListInIDs(ctx, ids)
		if err != nil {
			t.Fatalf("ListInIDs failed: %v", err)
		}

		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics, got %d", len(topics))
		}

		// Verify we got the expected topics
		topicIDs := make(map[string]bool)
		for _, topic := range topics {
			topicIDs[topic.ID] = true
		}

		if !topicIDs[createdTopic1.ID] {
			t.Errorf("Expected topic1 to be in results")
		}
		if !topicIDs[createdTopic3.ID] {
			t.Errorf("Expected topic3 to be in results")
		}
		if topicIDs[createdTopic2.ID] {
			t.Errorf("Expected topic2 to NOT be in results")
		}
	})

	t.Run("ListByMessageText", func(t *testing.T) {
		// Test filtering by message text containing "COVID"
		topics, err := app.Repository.Topics.ListLikeMessageText(ctx, "COVID", 0, 0)
		if err != nil {
			t.Fatalf("ListByMessageText failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with COVID messages, got %d", len(topics))
		}

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (COVID topic), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListLikeID", func(t *testing.T) {
		// Test filtering by ID pattern - look for topics with "001" in their ID
		topics, err := app.Repository.Topics.ListLikeID(ctx, "001", 0, 0)
		if err != nil {
			t.Fatalf("ListLikeID failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with '001' in ID, got %d", len(topics))
		}

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (ID contains '001'), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListLikeID - multiple matches", func(t *testing.T) {
		// Test filtering by ID pattern - look for topics with "550e8400" in their ID
		topics, err := app.Repository.Topics.ListLikeID(ctx, "550e8400", 0, 0)
		if err != nil {
			t.Fatalf("ListLikeID with multiple matches failed: %v", err)
		}

		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics with '550e8400' in ID, got %d", len(topics))
		}
	})

	t.Run("ListLikeID - no matches", func(t *testing.T) {
		// Test filtering by ID pattern that doesn't match any topics
		topics, err := app.Repository.Topics.ListLikeID(ctx, "nonexistent", 0, 0)
		if err != nil {
			t.Fatalf("ListLikeID with no matches failed: %v", err)
		}

		if len(topics) != 0 {
			t.Fatalf("Expected 0 topics with 'nonexistent' in ID, got %d", len(topics))
		}
	})

	t.Run("ListLikeIDAndMessageText", func(t *testing.T) {
		// Test filtering by both ID pattern and message text
		topics, err := app.Repository.Topics.ListLikeIDLikeMessageText(ctx, "001", "COVID", 0, 0)
		if err != nil {
			t.Fatalf("ListLikeIDAndMessageText failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with '001' in ID and COVID messages, got %d", len(topics))
		}

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (ID contains '001' and has COVID messages), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListLikeIDAndMessageText - no ID matches", func(t *testing.T) {
		// Test filtering by ID pattern that doesn't match and message text that does
		topics, err := app.Repository.Topics.ListLikeIDLikeMessageText(ctx, "nonexistent", "COVID", 0, 0)
		if err != nil {
			t.Fatalf("ListLikeIDAndMessageText with no ID matches failed: %v", err)
		}

		if len(topics) != 0 {
			t.Fatalf("Expected 0 topics with 'nonexistent' in ID and COVID messages, got %d", len(topics))
		}
	})

	t.Run("ListLikeIDAndMessageText - no message matches", func(t *testing.T) {
		// Test filtering by ID pattern that matches and message text that doesn't
		topics, err := app.Repository.Topics.ListLikeIDLikeMessageText(ctx, "001", "nonexistent", 0, 0)
		if err != nil {
			t.Fatalf("ListLikeIDAndMessageText with no message matches failed: %v", err)
		}

		if len(topics) != 0 {
			t.Fatalf("Expected 0 topics with '001' in ID and 'nonexistent' messages, got %d", len(topics))
		}
	})

	t.Run("ListLikeIDAndMessageText - both patterns match different topics", func(t *testing.T) {
		// Test filtering by ID pattern that matches topic1 and message text that matches topic2
		topics, err := app.Repository.Topics.ListLikeIDLikeMessageText(ctx, "001", "Election", 0, 0)
		if err != nil {
			t.Fatalf("ListLikeIDAndMessageText with different topic matches failed: %v", err)
		}

		if len(topics) != 0 {
			t.Fatalf("Expected 0 topics when ID pattern and message text match different topics, got %d", len(topics))
		}
	})
}

func TestRepository_AssignMessageToTopic(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()

	ctx := context.Background()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create two topics
	topic1 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440001",
		Name:         "Topic 1 - Original",
		Description:  "Original topic",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	topic2 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440002",
		Name:         "Topic 2 - Target",
		Description:  "Target topic for assignment",
		Status:       factcheck.StatusTopicResolved,
		Result:       "Verified",
		ResultStatus: factcheck.StatusTopicResultAnswered,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	// Create topics in database
	createdTopic1, err := app.Repository.Topics.Create(ctx, topic1)
	if err != nil {
		t.Fatalf("Failed to create topic1: %v", err)
	}

	createdTopic2, err := app.Repository.Topics.Create(ctx, topic2)
	if err != nil {
		t.Fatalf("Failed to create topic2: %v", err)
	}

	// Create a user message first
	userMessage := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440001",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-1"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	createdUserMessage, err := app.Repository.UserMessages.Create(ctx, userMessage)
	if err != nil {
		t.Fatalf("Failed to create userMessage: %v", err)
	}

	// Create a message in topic1
	message := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440001",
		UserMessageID: createdUserMessage.ID,
		TopicID:       createdTopic1.ID,
		Text:          "This message should be moved to topic2",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	createdMessage, err := app.Repository.Messages.Create(ctx, message)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Verify message is initially in topic1
	if createdMessage.TopicID != createdTopic1.ID {
		t.Fatalf("Message should initially be in topic1, but is in %s", createdMessage.TopicID)
	}

	// Verify topic1 has the message
	topic1Messages, err := app.Repository.Messages.ListByTopic(ctx, createdTopic1.ID)
	if err != nil {
		t.Fatalf("Failed to list messages for topic1: %v", err)
	}
	if len(topic1Messages) != 1 {
		t.Fatalf("Topic1 should have 1 message, but has %d", len(topic1Messages))
	}

	// Verify topic2 has no messages initially
	topic2Messages, err := app.Repository.Messages.ListByTopic(ctx, createdTopic2.ID)
	if err != nil {
		t.Fatalf("Failed to list messages for topic2: %v", err)
	}
	if len(topic2Messages) != 0 {
		t.Fatalf("Topic2 should have 0 messages initially, but has %d", len(topic2Messages))
	}

	t.Run("AssignMessageToTopic", func(t *testing.T) {
		// Assign message to topic2
		updatedMessage, err := app.Repository.Messages.AssignTopic(ctx, createdMessage.ID, createdTopic2.ID)
		if err != nil {
			t.Fatalf("Failed to assign message to topic2: %v", err)
		}

		// Verify the returned message has the new topic ID
		if updatedMessage.TopicID != createdTopic2.ID {
			t.Errorf("Updated message should have topic2 ID, but has %s", updatedMessage.TopicID)
		}

		// Verify other fields remain the same
		if updatedMessage.ID != createdMessage.ID {
			t.Errorf("Message ID should remain the same")
		}
		if updatedMessage.Text != createdMessage.Text {
			t.Errorf("Message text should remain the same")
		}
		if updatedMessage.Type != createdMessage.Type {
			t.Errorf("Message type should remain the same")
		}

		// Verify topic1 no longer has the message
		topic1MessagesAfter, err := app.Repository.Messages.ListByTopic(ctx, createdTopic1.ID)
		if err != nil {
			t.Fatalf("Failed to list messages for topic1 after assignment: %v", err)
		}
		if len(topic1MessagesAfter) != 0 {
			t.Fatalf("Topic1 should have 0 messages after assignment, but has %d", len(topic1MessagesAfter))
		}

		// Verify topic2 now has the message
		topic2MessagesAfter, err := app.Repository.Messages.ListByTopic(ctx, createdTopic2.ID)
		if err != nil {
			t.Fatalf("Failed to list messages for topic2 after assignment: %v", err)
		}
		if len(topic2MessagesAfter) != 1 {
			t.Fatalf("Topic2 should have 1 message after assignment, but has %d", len(topic2MessagesAfter))
		}

		// Verify the message in topic2 is the correct one
		if topic2MessagesAfter[0].ID != createdMessage.ID {
			t.Errorf("Topic2 should contain the assigned message")
		}
	})

	t.Run("AssignMessageToTopic - Invalid Message ID", func(t *testing.T) {
		// Try to assign a non-existent message
		_, err := app.Repository.Messages.AssignTopic(ctx, "non-existent-id", createdTopic2.ID)
		if err == nil {
			t.Fatalf("Expected error when assigning non-existent message")
		}
	})

	t.Run("AssignMessageToTopic - Invalid Topic ID", func(t *testing.T) {
		// Create another user message to test with
		userMessage2 := factcheck.UserMessage{
			ID:        "770e8400-e29b-41d4-a716-446655440002",
			Type:      factcheck.TypeUserMessageAdmin,
			RepliedAt: nil,
			Metadata:  []byte(`{"user_id": "test-user-2"}`),
			CreatedAt: now,
			UpdatedAt: nil,
		}

		createdUserMessage2, err := app.Repository.UserMessages.Create(ctx, userMessage2)
		if err != nil {
			t.Fatalf("Failed to create userMessage2: %v", err)
		}

		// Create another message to test with
		message2 := factcheck.Message{
			ID:            "660e8400-e29b-41d4-a716-446655440002",
			UserMessageID: createdUserMessage2.ID,
			TopicID:       createdTopic1.ID,
			Text:          "Another test message",
			Type:          factcheck.TypeMessageText,
			Status:        factcheck.StatusMessageTopicSubmitted,
			CreatedAt:     now,
			UpdatedAt:     nil,
		}

		createdMessage2, err := app.Repository.Messages.Create(ctx, message2)
		if err != nil {
			t.Fatalf("Failed to create message2: %v", err)
		}

		// Try to assign to a non-existent topic
		_, err = app.Repository.Messages.AssignTopic(ctx, createdMessage2.ID, "non-existent-topic-id")
		if err == nil {
			t.Fatalf("Expected error when assigning to non-existent topic")
		}
	})
}

func TestRepository_ListHome(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()

	ctx := context.Background()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create topics with different statuses
	topic1 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440001",
		Name:         "Topic 1 - COVID Pending",
		Description:  "COVID-19 related news (pending)",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	topic2 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440002",
		Name:         "Topic 2 - Politics Resolved",
		Description:  "Political news and updates (resolved)",
		Status:       factcheck.StatusTopicResolved,
		Result:       "Verified as true",
		ResultStatus: factcheck.StatusTopicResultAnswered,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	topic3 := factcheck.Topic{
		ID:           "660e8400-e29b-41d4-a716-446655440003",
		Name:         "Topic 3 - Technology Pending",
		Description:  "Technology news and updates (pending)",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	// Create topics in database
	createdTopic1, err := app.Repository.Topics.Create(ctx, topic1)
	if err != nil {
		t.Fatalf("Failed to create topic1: %v", err)
	}

	createdTopic2, err := app.Repository.Topics.Create(ctx, topic2)
	if err != nil {
		t.Fatalf("Failed to create topic2: %v", err)
	}

	createdTopic3, err := app.Repository.Topics.Create(ctx, topic3)
	if err != nil {
		t.Fatalf("Failed to create topic3: %v", err)
	}

	// Create user messages first
	userMessage1 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440001",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-1"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage2 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440002",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-2"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage3 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440003",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-3"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	// Create user messages in database
	createdUserMessage1, err := app.Repository.UserMessages.Create(ctx, userMessage1)
	if err != nil {
		t.Fatalf("Failed to create userMessage1: %v", err)
	}

	createdUserMessage2, err := app.Repository.UserMessages.Create(ctx, userMessage2)
	if err != nil {
		t.Fatalf("Failed to create userMessage2: %v", err)
	}

	createdUserMessage3, err := app.Repository.UserMessages.Create(ctx, userMessage3)
	if err != nil {
		t.Fatalf("Failed to create userMessage3: %v", err)
	}

	// Create messages
	message1 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440001",
		UserMessageID: createdUserMessage1.ID,
		TopicID:       createdTopic1.ID,
		Text:          "COVID-19 vaccine is effective against new variants",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message2 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440002",
		UserMessageID: createdUserMessage2.ID,
		TopicID:       createdTopic2.ID,
		Text:          "Election results show clear victory",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message3 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440003",
		UserMessageID: createdUserMessage3.ID,
		TopicID:       createdTopic3.ID,
		Text:          "New AI technology breakthrough",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	// Create messages in database
	_, err = app.Repository.Messages.Create(ctx, message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message3)
	if err != nil {
		t.Fatalf("Failed to create message3: %v", err)
	}

	t.Run("ListHome - no options (all topics)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0)
		if err != nil {
			t.Fatalf("ListHome with no options failed: %v", err)
		}
		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - status filter only", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicStatus(factcheck.StatusTopicPending),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with status filter failed: %v", err)
		}
		if len(topics) != 2 {
			t.Fatalf("Expected 2 pending topics, got %d", len(topics))
		}
		// Verify we got the pending topics
		topicIDs := make(map[string]bool)
		for _, topic := range topics {
			topicIDs[topic.ID] = true
		}
		if !topicIDs[createdTopic1.ID] {
			t.Errorf("Expected topic1 (pending) to be in results")
		}
		if !topicIDs[createdTopic3.ID] {
			t.Errorf("Expected topic3 (pending) to be in results")
		}
		if topicIDs[createdTopic2.ID] {
			t.Errorf("Expected topic2 (resolved) to NOT be in results")
		}
	})

	t.Run("ListHome - message text filter only", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with message text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with COVID messages, got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (COVID topic), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListHome - ID pattern filter only", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with ID pattern filter failed: %v", err)
		}
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics with '550e8400' in ID, got %d", len(topics))
		}
		// Verify we got the expected topics
		topicIDs := make(map[string]bool)
		for _, topic := range topics {
			topicIDs[topic.ID] = true
		}
		if !topicIDs[createdTopic1.ID] {
			t.Errorf("Expected topic1 to be in results")
		}
		if !topicIDs[createdTopic2.ID] {
			t.Errorf("Expected topic2 to be in results")
		}
		if topicIDs[createdTopic3.ID] {
			t.Errorf("Expected topic3 to NOT be in results (different ID prefix)")
		}
	})

	t.Run("ListHome - ID pattern and message text filters", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeMessageText("COVID"),
		}
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with ID pattern and message text filters failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with '550e8400' in ID and COVID messages, got %d", len(topics))
		}

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (matches both filters), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListHome - status and message text filters", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicStatus(factcheck.StatusTopicPending),
			repo.WithTopicLikeMessageText("COVID"),
		}
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with status and message text filters failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 pending topic with COVID messages, got %d", len(topics))
		}

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (pending with COVID messages), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListHome - status and ID pattern filters", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicStatus(factcheck.StatusTopicResolved),
			repo.WithTopicLikeID("550e8400"),
		}
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with status and ID pattern filters failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 resolved topic with '550e8400' in ID, got %d", len(topics))
		}

		if topics[0].ID != createdTopic2.ID {
			t.Errorf("Expected topic2 (resolved with '550e8400' in ID), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListHome - all three filters", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicStatus(factcheck.StatusTopicPending),
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeMessageText("COVID"),
		}
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with all three filters failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 pending topic with '550e8400' in ID and COVID messages, got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (matches all three filters), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListHome - no matches for combined filters", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicStatus(factcheck.StatusTopicResolved),
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeMessageText("COVID"),
		}
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with no matches failed: %v", err)
		}
		if len(topics) != 0 {
			t.Fatalf("Expected 0 topics (no resolved topics with '550e8400' in ID and COVID messages), got %d", len(topics))
		}
	})

	t.Run("ListHome - empty string filters", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID(""),
			repo.WithTopicLikeMessageText(""),
		}
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with empty string filters failed: %v", err)
		}

		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics (empty filters should return all), got %d", len(topics))
		}
	})

	t.Run("ListHome - multiple options of same type (last one wins)", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeID("660e8400"),
		}
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListHome with multiple ID filters failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with '660e8400' in ID (last filter wins), got %d", len(topics))
		}
		if topics[0].ID != createdTopic3.ID {
			t.Errorf("Expected topic3 (matches last ID filter), got topic with ID %s", topics[0].ID)
		}
	})
}

func TestRepository_CountByStatusesHome(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()

	ctx := context.Background()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create topics with different statuses
	topic1 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440001",
		Name:         "Topic 1 - COVID Pending",
		Description:  "COVID-19 related news (pending)",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	topic2 := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440002",
		Name:         "Topic 2 - Politics Resolved",
		Description:  "Political news and updates (resolved)",
		Status:       factcheck.StatusTopicResolved,
		Result:       "Verified as true",
		ResultStatus: factcheck.StatusTopicResultAnswered,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	topic3 := factcheck.Topic{
		ID:           "660e8400-e29b-41d4-a716-446655440003",
		Name:         "Topic 3 - Technology Pending",
		Description:  "Technology news and updates (pending)",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	topic4 := factcheck.Topic{
		ID:           "660e8400-e29b-41d4-a716-446655440004",
		Name:         "Topic 4 - Sports Resolved",
		Description:  "Sports news and updates (resolved)",
		Status:       factcheck.StatusTopicResolved,
		Result:       "Verified as false",
		ResultStatus: factcheck.StatusTopicResultAnswered,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	// Create topics in database
	createdTopic1, err := app.Repository.Topics.Create(ctx, topic1)
	if err != nil {
		t.Fatalf("Failed to create topic1: %v", err)
	}

	createdTopic2, err := app.Repository.Topics.Create(ctx, topic2)
	if err != nil {
		t.Fatalf("Failed to create topic2: %v", err)
	}

	createdTopic3, err := app.Repository.Topics.Create(ctx, topic3)
	if err != nil {
		t.Fatalf("Failed to create topic3: %v", err)
	}

	createdTopic4, err := app.Repository.Topics.Create(ctx, topic4)
	if err != nil {
		t.Fatalf("Failed to create topic4: %v", err)
	}

	// Create user messages first
	userMessage1 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440001",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-1"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage2 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440002",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-2"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage3 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440003",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-3"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage4 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440004",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-4"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	// Create user messages in database
	createdUserMessage1, err := app.Repository.UserMessages.Create(ctx, userMessage1)
	if err != nil {
		t.Fatalf("Failed to create userMessage1: %v", err)
	}

	createdUserMessage2, err := app.Repository.UserMessages.Create(ctx, userMessage2)
	if err != nil {
		t.Fatalf("Failed to create userMessage2: %v", err)
	}

	createdUserMessage3, err := app.Repository.UserMessages.Create(ctx, userMessage3)
	if err != nil {
		t.Fatalf("Failed to create userMessage3: %v", err)
	}

	createdUserMessage4, err := app.Repository.UserMessages.Create(ctx, userMessage4)
	if err != nil {
		t.Fatalf("Failed to create userMessage4: %v", err)
	}

	// Create messages
	message1 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440001",
		UserMessageID: createdUserMessage1.ID,
		TopicID:       createdTopic1.ID,
		Text:          "COVID-19 vaccine is effective against new variants",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message2 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440002",
		UserMessageID: createdUserMessage2.ID,
		TopicID:       createdTopic2.ID,
		Text:          "Election results show clear victory",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message3 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440003",
		UserMessageID: createdUserMessage3.ID,
		TopicID:       createdTopic3.ID,
		Text:          "New AI technology breakthrough",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message4 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440004",
		UserMessageID: createdUserMessage4.ID,
		TopicID:       createdTopic4.ID,
		Text:          "World Cup final results announced",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	// Create messages in database
	_, err = app.Repository.Messages.Create(ctx, message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message3)
	if err != nil {
		t.Fatalf("Failed to create message3: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message4)
	if err != nil {
		t.Fatalf("Failed to create message4: %v", err)
	}

	t.Run("CountByStatusesHome - no filters (all topics)", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx)
		if err != nil {
			t.Fatalf("CountByStatusesHome with no filters failed: %v", err)
		}

		expectedPending := int64(2)  // topic1, topic3
		expectedResolved := int64(2) // topic2, topic4

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - ID pattern filter only", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeID("550e8400"),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with ID pattern filter failed: %v", err)
		}

		expectedPending := int64(1)  // topic1
		expectedResolved := int64(1) // topic2

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics with '550e8400' in ID, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics with '550e8400' in ID, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - message text filter only", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeMessageText("COVID"),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with message text filter failed: %v", err)
		}

		expectedPending := int64(1)  // topic1
		expectedResolved := int64(0) // none

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics with COVID messages, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics with COVID messages, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - ID pattern and message text filters", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeMessageText("COVID"),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with ID pattern and message text filters failed: %v", err)
		}

		expectedPending := int64(1)  // topic1
		expectedResolved := int64(0) // none

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics with '550e8400' in ID and COVID messages, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics with '550e8400' in ID and COVID messages, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - message text filter for resolved topics", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeMessageText("Election"),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with Election message filter failed: %v", err)
		}

		expectedPending := int64(0)  // none
		expectedResolved := int64(1) // topic2

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics with Election messages, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics with Election messages, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - ID pattern filter for different prefix", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeID("660e8400"),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with '660e8400' ID filter failed: %v", err)
		}

		expectedPending := int64(1)  // topic3
		expectedResolved := int64(1) // topic4

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics with '660e8400' in ID, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics with '660e8400' in ID, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - combined filters for technology topics", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeID("660e8400"),
			repo.WithTopicLikeMessageText("technology"),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with technology filters failed: %v", err)
		}

		expectedPending := int64(1)  // topic3
		expectedResolved := int64(0) // none

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics with '660e8400' in ID and technology messages, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics with '660e8400' in ID and technology messages, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - no matches for combined filters", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeMessageText("techonology"),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with no matches filter failed: %v", err)
		}

		expectedPending := int64(0)  // none
		expectedResolved := int64(0) // none

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics with '550e8400' in ID and technology messages, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics with '550e8400' in ID and technology messages, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - empty string filters", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeID(""),
			repo.WithTopicLikeMessageText(""),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with empty string filters failed: %v", err)
		}

		expectedPending := int64(2)  // topic1, topic3
		expectedResolved := int64(2) // topic2, topic4

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics (empty filters should return all), got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics (empty filters should return all), got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusesHome - case insensitive message text filter", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusHome(ctx,
			repo.WithTopicLikeMessageText("covid"),
		)
		if err != nil {
			t.Fatalf("CountByStatusesHome with lowercase COVID filter failed: %v", err)
		}

		expectedPending := int64(1)  // topic1
		expectedResolved := int64(0) // none

		if counts[factcheck.StatusTopicPending] != expectedPending {
			t.Errorf("Expected %d pending topics with lowercase COVID messages, got %d", expectedPending, counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != expectedResolved {
			t.Errorf("Expected %d resolved topics with lowercase COVID messages, got %d", expectedResolved, counts[factcheck.StatusTopicResolved])
		}
	})
}

func TestRepository_TopicPagination(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create 10 topics for pagination testing
	ctx := context.Background()
	createdTopics := make([]factcheck.Topic, 10)
	for i := 0; i < 10; i++ {
		n := i + 1
		topic := factcheck.Topic{
			ID:           fmt.Sprintf("550e8400-e29b-41d4-a716-446655440%03d", n),
			Name:         fmt.Sprintf("Topic %d - Test Pagination", n),
			Description:  fmt.Sprintf("Description for topic %d", n),
			Status:       factcheck.StatusTopicPending,
			Result:       "",
			ResultStatus: factcheck.StatusTopicResultNone,
			CreatedAt:    now,
			UpdatedAt:    nil,
		}

		createdTopic, err := app.Repository.Topics.Create(ctx, topic)
		if err != nil {
			t.Fatalf("Failed to create topics[%d] (#%d): %v", i, n, err)
		}
		createdTopics[i] = createdTopic
	}

	// Create user messages and messages for filtering tests
	userMessage1 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440001",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-1"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	userMessage2 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440002",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "test-user-2"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	createdUserMessage1, err := app.Repository.UserMessages.Create(ctx, userMessage1)
	if err != nil {
		t.Fatalf("Failed to create userMessage1: %v", err)
	}

	createdUserMessage2, err := app.Repository.UserMessages.Create(ctx, userMessage2)
	if err != nil {
		t.Fatalf("Failed to create userMessage2: %v", err)
	}

	// Create messages with specific text for filtering
	message1 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440001",
		UserMessageID: createdUserMessage1.ID,
		TopicID:       createdTopics[0].ID,
		Text:          "COVID-19 vaccine is effective",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message2 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440002",
		UserMessageID: createdUserMessage2.ID,
		TopicID:       createdTopics[1].ID,
		Text:          "Election results show victory",
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	_, err = app.Repository.Messages.Create(ctx, message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	_, err = app.Repository.Messages.Create(ctx, message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	t.Run("List - no pagination (limit=0)", func(t *testing.T) {
		topics, err := app.Repository.Topics.List(ctx, 0, 0)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return all 10 topics when limit=0
		if len(topics) != 10 {
			t.Errorf("Expected 10 topics, got %d", len(topics))
		}
	})

	t.Run("List - normal pagination (limit=3, offset=0)", func(t *testing.T) {
		topics, err := app.Repository.Topics.List(ctx, 3, 0)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		if len(topics) != 3 {
			t.Errorf("Expected 3 topics, got %d", len(topics))
		}
	})

	t.Run("List - pagination with offset (limit=3, offset=3)", func(t *testing.T) {
		topics, err := app.Repository.Topics.List(ctx, 3, 3)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		if len(topics) != 3 {
			t.Errorf("Expected 3 topics, got %d", len(topics))
		}
	})

	t.Run("List - pagination with large offset", func(t *testing.T) {
		topics, err := app.Repository.Topics.List(ctx, 5, 8)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return remaining 2 topics (10 total - 8 offset)
		if len(topics) != 2 {
			t.Errorf("Expected 2 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - no pagination (limit=0)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return all 10 topics when limit=0
		if len(topics) != 10 {
			t.Errorf("Expected 10 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - no pagination (limit=0)", func(t *testing.T) {
		tx, err := app.Repository.BeginTx(ctx, repo.Serializable)
		if err != nil {
			t.Fatalf("Failed to begin transaction: %v", err)
		}
		opts := []repo.TopicOption{
			repo.WithTopicTx(tx),
		}
		topics, err := app.Repository.Topics.ListHome(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return all 10 topics when limit=0
		if len(topics) != 10 {
			t.Errorf("Expected 10 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - normal pagination (limit=3, offset=0)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 3, 0)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		if len(topics) != 3 {
			t.Errorf("Expected 3 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - pagination with filter and limit", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return 5 topics (all match the ID pattern)
		if len(topics) != 5 {
			t.Errorf("Expected 5 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - pagination with message text filter", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 10, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return 1 topic (only topic1 has COVID message)
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
		if topics[0].ID != createdTopics[0].ID {
			t.Errorf("Expected topic ID %s, got %s", createdTopics[0].ID, topics[0].ID)
		}
	})

	t.Run("ListHome - pagination with both filters", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return 1 topic (matches both filters)
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
		if topics[0].ID != createdTopics[0].ID {
			t.Errorf("Expected topic ID %s, got %s", createdTopics[0].ID, topics[0].ID)
		}
	})

	t.Run("ListHome - negative limit (should return all)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, -1, 0)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return all 10 topics when limit is negative
		if len(topics) != 10 {
			t.Errorf("Expected 10 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - negative offset (should start from beginning)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 3, -2)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return first 3 topics (negative offset treated as 0)
		if len(topics) != 3 {
			t.Errorf("Expected 3 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - large limit (should cap at available records)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 100, 0)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return all 10 topics (limit capped at available records)
		if len(topics) != 10 {
			t.Errorf("Expected 10 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - pagination order consistency", func(t *testing.T) {
		// Get first page
		topics1, err := app.Repository.Topics.ListHome(ctx, 3, 0)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		if len(topics1) != 3 {
			t.Errorf("Expected 3 topics, got %d", len(topics1))
		}
		// Get second page
		topics2, err := app.Repository.Topics.ListHome(ctx, 3, 3)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		if len(topics2) != 3 {
			t.Errorf("Expected 3 topics, got %d", len(topics2))
		}
		// Verify no overlap between pages
		topicIDs1 := make(map[string]bool)
		for _, topic := range topics1 {
			topicIDs1[topic.ID] = true
		}
		for _, topic := range topics2 {
			if topicIDs1[topic.ID] {
				t.Errorf("Topic %s appears in both pages", topic.ID)
			}
		}
	})

	t.Run("ListHome - pagination with empty result set", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("nonexistent"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics: %v", err)
		}
		// Should return empty array
		if len(topics) != 0 {
			t.Errorf("Expected 0 topics, got %d", len(topics))
		}
	})

	t.Run("ListByStatus - pagination", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListByStatus(ctx, factcheck.StatusTopicPending, 3, 0)
		if err != nil {
			t.Fatalf("Failed to list topics by status: %v", err)
		}
		if len(topics) != 3 {
			t.Errorf("Expected 3 topics, got %d", len(topics))
		}
		// Verify all topics have the correct status
		for _, topic := range topics {
			if topic.Status != factcheck.StatusTopicPending {
				t.Errorf("Expected status %s, got %s", factcheck.StatusTopicPending, topic.Status)
			}
		}
	})

	t.Run("ListLikeMessageText - pagination", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListLikeMessageText(ctx, "COVID", 5, 0)
		if err != nil {
			t.Fatalf("Failed to list topics by message text: %v", err)
		}
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
		if topics[0].ID != createdTopics[0].ID {
			t.Errorf("Expected topic ID %s, got %s", createdTopics[0].ID, topics[0].ID)
		}
	})

	t.Run("ListLikeID - pagination", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListLikeID(ctx, "550e8400", 5, 0)
		if err != nil {
			t.Fatalf("Failed to list topics by ID pattern: %v", err)
		}
		if len(topics) != 5 {
			t.Errorf("Expected 5 topics, got %d", len(topics))
		}
	})

	t.Run("ListLikeIDLikeMessageText - pagination", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListLikeIDLikeMessageText(ctx, "550e8400", "COVID", 5, 0)
		if err != nil {
			t.Fatalf("Failed to list topics by ID and message text: %v", err)
		}
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
		if topics[0].ID != createdTopics[0].ID {
			t.Errorf("Expected topic ID %s, got %s", createdTopics[0].ID, topics[0].ID)
		}
	})

	// Edge Case 3: Very large values
	t.Run("ListHome - very large limit", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 999999, 0)
		if err != nil {
			t.Fatalf("Failed to list topics with large limit: %v", err)
		}
		// Should return all 10 topics (capped at available records)
		if len(topics) != 10 {
			t.Errorf("Expected 10 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - very large offset", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 5, 999999)
		if err != nil {
			t.Fatalf("Failed to list topics with large offset: %v", err)
		}
		// Should return empty array (offset beyond available records)
		if len(topics) != 0 {
			t.Errorf("Expected 0 topics, got %d", len(topics))
		}
	})

	// Edge Case 4: Single record scenarios
	t.Run("ListHome - single record limit", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 1, 0)
		if err != nil {
			t.Fatalf("Failed to list topics with limit=1: %v", err)
		}
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
	})

	t.Run("ListHome - single record with offset", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 1, 1)
		if err != nil {
			t.Fatalf("Failed to list topics with limit=1, offset=1: %v", err)
		}
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
	})

	// Edge Case 5: Empty result sets with pagination
	t.Run("ListHome - empty result with pagination", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("nonexistent"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics with empty result: %v", err)
		}
		if len(topics) != 0 {
			t.Errorf("Expected 0 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - empty result with offset", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("nonexistent"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 10, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics with empty result and offset: %v", err)
		}
		if len(topics) != 0 {
			t.Errorf("Expected 0 topics, got %d", len(topics))
		}
	})

	// Edge Case 6: Single page scenarios
	t.Run("ListHome - single page (limit equals total)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 10, 0)
		if err != nil {
			t.Fatalf("Failed to list topics with limit=total: %v", err)
		}
		if len(topics) != 10 {
			t.Errorf("Expected 10 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - single page with filter", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 1, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics with single page filter: %v", err)
		}
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
	})

	// Edge Case 7: Last page scenarios
	t.Run("ListHome - last page (offset + limit > total)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 5, 8)
		if err != nil {
			t.Fatalf("Failed to list topics on last page: %v", err)
		}
		// Should return remaining 2 topics (10 total - 8 offset)
		if len(topics) != 2 {
			t.Errorf("Expected 2 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - last page with filter", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 5, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics on last page with filter: %v", err)
		}
		// Should return remaining topics after offset
		if len(topics) != 5 {
			t.Errorf("Expected 5 topics, got %d", len(topics))
		}
	})

	// Edge Case 8: Exact boundary scenarios
	t.Run("ListHome - exact boundary (offset + limit = total)", func(t *testing.T) {
		topics, err := app.Repository.Topics.ListHome(ctx, 5, 5)
		if err != nil {
			t.Fatalf("Failed to list topics at exact boundary: %v", err)
		}
		// Should return exactly 5 topics (10 total - 5 offset)
		if len(topics) != 5 {
			t.Errorf("Expected 5 topics, got %d", len(topics))
		}
	})

	t.Run("ListHome - exact boundary with filter", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 1, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics at exact boundary with filter: %v", err)
		}
		// Should return exactly 1 topic
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
	})

	// Edge Case 9: Multiple filters with pagination
	t.Run("ListHome - multiple filters with pagination", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 3, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics with multiple filters: %v", err)
		}
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
	})

	t.Run("ListHome - multiple filters with offset", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 3, 1, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics with multiple filters and offset: %v", err)
		}
		// Should return 0 topics (only 1 matches, offset=1)
		if len(topics) != 0 {
			t.Errorf("Expected 0 topics, got %d", len(topics))
		}
	})

	// Edge Case 10: Filter that returns exactly limit records
	t.Run("ListHome - filter returns exactly limit records", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeID("550e8400"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics with exact limit: %v", err)
		}
		// Should return exactly 5 topics
		if len(topics) != 5 {
			t.Errorf("Expected 5 topics, got %d", len(topics))
		}
	})

	// Edge Case 11: Filter that returns fewer than limit records
	t.Run("ListHome - filter returns fewer than limit", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 0, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics with fewer than limit: %v", err)
		}
		// Should return 1 topic (fewer than limit of 5)
		if len(topics) != 1 {
			t.Errorf("Expected 1 topic, got %d", len(topics))
		}
	})

	t.Run("ListHome - filter returns fewer than limit with offset", func(t *testing.T) {
		opts := []repo.TopicOption{
			repo.WithTopicLikeMessageText("COVID"),
		}

		topics, err := app.Repository.Topics.ListHome(ctx, 5, 1, opts...)
		if err != nil {
			t.Fatalf("Failed to list topics with fewer than limit and offset: %v", err)
		}
		// Should return 0 topics (only 1 matches, offset=1)
		if len(topics) != 0 {
			t.Errorf("Expected 0 topics, got %d", len(topics))
		}
	})
}
