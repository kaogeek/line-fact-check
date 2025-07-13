//go:build integration_test
// +build integration_test

package repo_test

import (
	"context"
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestRepository_TopicFiltering(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()

	ctx := context.Background()

	// Clear all data
	t.Log("Clearing all data from database")
	_, err = app.PostgresConn.Exec(ctx, "DELETE FROM user_messages")
	if err != nil {
		t.Fatalf("Failed to clear user_messages: %v", err)
	}
	_, err = app.PostgresConn.Exec(ctx, "DELETE FROM messages")
	if err != nil {
		t.Fatalf("Failed to clear messages: %v", err)
	}
	_, err = app.PostgresConn.Exec(ctx, "DELETE FROM topics")
	if err != nil {
		t.Fatalf("Failed to clear topics: %v", err)
	}

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
	createdTopic1, err := app.Repository.Topic.Create(ctx, topic1)
	if err != nil {
		t.Fatalf("Failed to create topic1: %v", err)
	}

	createdTopic2, err := app.Repository.Topic.Create(ctx, topic2)
	if err != nil {
		t.Fatalf("Failed to create topic2: %v", err)
	}

	createdTopic3, err := app.Repository.Topic.Create(ctx, topic3)
	if err != nil {
		t.Fatalf("Failed to create topic3: %v", err)
	}

	// Create messages
	message1 := factcheck.Message{
		ID:        "660e8400-e29b-41d4-a716-446655440001",
		TopicID:   createdTopic1.ID,
		Text:      "COVID-19 vaccine is effective against new variants",
		Type:      factcheck.TypeMessageText,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	message2 := factcheck.Message{
		ID:        "660e8400-e29b-41d4-a716-446655440002",
		TopicID:   createdTopic1.ID,
		Text:      "COVID-19 cases are increasing in winter",
		Type:      factcheck.TypeMessageText,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	message3 := factcheck.Message{
		ID:        "660e8400-e29b-41d4-a716-446655440003",
		TopicID:   createdTopic2.ID,
		Text:      "Election results show clear victory",
		Type:      factcheck.TypeMessageText,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	message4 := factcheck.Message{
		ID:        "660e8400-e29b-41d4-a716-446655440004",
		TopicID:   createdTopic3.ID,
		Text:      "New AI technology breakthrough",
		Type:      factcheck.TypeMessageText,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	// Create messages in database
	// These messages are used indirectly through topic filtering tests
	_, err = app.Repository.Message.Create(ctx, message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	_, err = app.Repository.Message.Create(ctx, message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	_, err = app.Repository.Message.Create(ctx, message3)
	if err != nil {
		t.Fatalf("Failed to create message3: %v", err)
	}

	_, err = app.Repository.Message.Create(ctx, message4)
	if err != nil {
		t.Fatalf("Failed to create message4: %v", err)
	}

	t.Run("ListInIDs", func(t *testing.T) {
		// Test filtering by specific IDs
		ids := []string{createdTopic1.ID, createdTopic3.ID}
		topics, err := app.Repository.Topic.ListInIDs(ctx, ids)
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
		topics, err := app.Repository.Topic.ListByMessageText(ctx, "COVID")
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

	t.Run("ListInIDsAndMessageText", func(t *testing.T) {
		// Test filtering by both IDs and message text
		ids := []string{createdTopic1.ID, createdTopic2.ID, createdTopic3.ID}
		topics, err := app.Repository.Topic.ListInIDsAndMessageText(ctx, ids, "COVID")
		if err != nil {
			t.Fatalf("ListInIDsAndMessageText failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with COVID messages from the specified IDs, got %d", len(topics))
		}

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (COVID topic), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListFiltered - IDs only", func(t *testing.T) {
		// Test wrapper method with IDs only
		ids := []string{createdTopic1.ID, createdTopic2.ID}
		topics, err := app.Repository.Topic.ListFiltered(ctx, ids, "")
		if err != nil {
			t.Fatalf("ListFiltered with IDs only failed: %v", err)
		}

		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics, got %d", len(topics))
		}
	})

	t.Run("ListFiltered - message text only", func(t *testing.T) {
		// Test wrapper method with message text only
		topics, err := app.Repository.Topic.ListFiltered(ctx, nil, "Election")
		if err != nil {
			t.Fatalf("ListFiltered with message text only failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with Election messages, got %d", len(topics))
		}

		if topics[0].ID != createdTopic2.ID {
			t.Errorf("Expected topic2 (Politics topic), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListFiltered - both filters", func(t *testing.T) {
		// Test wrapper method with both filters
		ids := []string{createdTopic1.ID, createdTopic2.ID, createdTopic3.ID}
		topics, err := app.Repository.Topic.ListFiltered(ctx, ids, "COVID")
		if err != nil {
			t.Fatalf("ListFiltered with both filters failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with COVID messages from the specified IDs, got %d", len(topics))
		}

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (COVID topic), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListFiltered - no filters", func(t *testing.T) {
		// Test wrapper method with no filters
		topics, err := app.Repository.Topic.ListFiltered(ctx, nil, "")
		if err != nil {
			t.Fatalf("ListFiltered with no filters failed: %v", err)
		}

		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics (all), got %d", len(topics))
		}
	})

	t.Run("ListFiltered - empty IDs", func(t *testing.T) {
		// Test wrapper method with empty IDs
		topics, err := app.Repository.Topic.ListFiltered(ctx, []string{}, "COVID")
		if err != nil {
			t.Fatalf("ListFiltered with empty IDs failed: %v", err)
		}

		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with COVID messages, got %d", len(topics))
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

	// Clear all data
	t.Log("Clearing all data from database")
	_, err = app.PostgresConn.Exec(ctx, "DELETE FROM user_messages")
	if err != nil {
		t.Fatalf("Failed to clear user_messages: %v", err)
	}
	_, err = app.PostgresConn.Exec(ctx, "DELETE FROM messages")
	if err != nil {
		t.Fatalf("Failed to clear messages: %v", err)
	}
	_, err = app.PostgresConn.Exec(ctx, "DELETE FROM topics")
	if err != nil {
		t.Fatalf("Failed to clear topics: %v", err)
	}

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
	createdTopic1, err := app.Repository.Topic.Create(ctx, topic1)
	if err != nil {
		t.Fatalf("Failed to create topic1: %v", err)
	}

	createdTopic2, err := app.Repository.Topic.Create(ctx, topic2)
	if err != nil {
		t.Fatalf("Failed to create topic2: %v", err)
	}

	// Create a message in topic1
	message := factcheck.Message{
		ID:        "660e8400-e29b-41d4-a716-446655440001",
		TopicID:   createdTopic1.ID,
		Text:      "This message should be moved to topic2",
		Type:      factcheck.TypeMessageText,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	createdMessage, err := app.Repository.Message.Create(ctx, message)
	if err != nil {
		t.Fatalf("Failed to create message: %v", err)
	}

	// Verify message is initially in topic1
	if createdMessage.TopicID != createdTopic1.ID {
		t.Fatalf("Message should initially be in topic1, but is in %s", createdMessage.TopicID)
	}

	// Verify topic1 has the message
	topic1Messages, err := app.Repository.Message.ListByTopic(ctx, createdTopic1.ID)
	if err != nil {
		t.Fatalf("Failed to list messages for topic1: %v", err)
	}
	if len(topic1Messages) != 1 {
		t.Fatalf("Topic1 should have 1 message, but has %d", len(topic1Messages))
	}

	// Verify topic2 has no messages initially
	topic2Messages, err := app.Repository.Message.ListByTopic(ctx, createdTopic2.ID)
	if err != nil {
		t.Fatalf("Failed to list messages for topic2: %v", err)
	}
	if len(topic2Messages) != 0 {
		t.Fatalf("Topic2 should have 0 messages initially, but has %d", len(topic2Messages))
	}

	t.Run("AssignMessageToTopic", func(t *testing.T) {
		// Assign message to topic2
		updatedMessage, err := app.Repository.Message.AssignToTopic(ctx, createdMessage.ID, createdTopic2.ID)
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
		topic1MessagesAfter, err := app.Repository.Message.ListByTopic(ctx, createdTopic1.ID)
		if err != nil {
			t.Fatalf("Failed to list messages for topic1 after assignment: %v", err)
		}
		if len(topic1MessagesAfter) != 0 {
			t.Fatalf("Topic1 should have 0 messages after assignment, but has %d", len(topic1MessagesAfter))
		}

		// Verify topic2 now has the message
		topic2MessagesAfter, err := app.Repository.Message.ListByTopic(ctx, createdTopic2.ID)
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
		_, err := app.Repository.Message.AssignToTopic(ctx, "non-existent-id", createdTopic2.ID)
		if err == nil {
			t.Fatalf("Expected error when assigning non-existent message")
		}
	})

	t.Run("AssignMessageToTopic - Invalid Topic ID", func(t *testing.T) {
		// Create another message to test with
		message2 := factcheck.Message{
			ID:        "660e8400-e29b-41d4-a716-446655440002",
			TopicID:   createdTopic1.ID,
			Text:      "Another test message",
			Type:      factcheck.TypeMessageText,
			CreatedAt: now,
			UpdatedAt: nil,
		}

		createdMessage2, err := app.Repository.Message.Create(ctx, message2)
		if err != nil {
			t.Fatalf("Failed to create message2: %v", err)
		}

		// Try to assign to a non-existent topic
		_, err = app.Repository.Message.AssignToTopic(ctx, createdMessage2.ID, "non-existent-topic-id")
		if err == nil {
			t.Fatalf("Expected error when assigning to non-existent topic")
		}
	})
}
