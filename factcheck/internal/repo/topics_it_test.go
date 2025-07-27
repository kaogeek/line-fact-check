//go:build integration_test
// +build integration_test

package repo_test

import (
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestRepository_CountByStatusDynamic(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()
	ctx := t.Context()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create test topics
	topic1 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440001",
		Name:        "Topic 1 - COVID Pending",
		Description: "COVID-19 related topic",
		Status:      factcheck.StatusTopicPending,
		CreatedAt:   now,
		UpdatedAt:   nil,
	}
	topic2 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440002",
		Name:        "Topic 2 - Politics Resolved",
		Description: "Politics related topic",
		Status:      factcheck.StatusTopicResolved,
		CreatedAt:   now,
		UpdatedAt:   nil,
	}
	topic3 := factcheck.Topic{
		ID:          "660e8400-e29b-41d4-a716-446655440003",
		Name:        "Topic 3 - Technology Pending",
		Description: "Technology related topic",
		Status:      factcheck.StatusTopicPending,
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

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

	// Create test messages
	userMessage1 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440001",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "user1"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}
	userMessage2 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440002",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "user2"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}
	userMessage3 := factcheck.UserMessage{
		ID:        "770e8400-e29b-41d4-a716-446655440003",
		Type:      factcheck.TypeUserMessageAdmin,
		RepliedAt: nil,
		Metadata:  []byte(`{"user_id": "user3"}`),
		CreatedAt: now,
		UpdatedAt: nil,
	}

	_, err = app.Repository.UserMessages.Create(ctx, userMessage1)
	if err != nil {
		t.Fatalf("Failed to create userMessage1: %v", err)
	}
	_, err = app.Repository.UserMessages.Create(ctx, userMessage2)
	if err != nil {
		t.Fatalf("Failed to create userMessage2: %v", err)
	}
	_, err = app.Repository.UserMessages.Create(ctx, userMessage3)
	if err != nil {
		t.Fatalf("Failed to create userMessage3: %v", err)
	}

	message1 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440001",
		UserMessageID: userMessage1.ID,
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		TopicID:       createdTopic1.ID,
		Text:          "This is a COVID-19 related message",
		Language:      "en",
		CreatedAt:     now,
		UpdatedAt:     nil,
	}
	message2 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440002",
		UserMessageID: userMessage2.ID,
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		TopicID:       createdTopic2.ID,
		Text:          "นี่คือข่าวปลอมเกี่ยวกับการเมือง",
		Language:      "th",
		CreatedAt:     now,
		UpdatedAt:     nil,
	}
	message3 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440003",
		UserMessageID: userMessage3.ID,
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		TopicID:       createdTopic3.ID,
		Text:          "This is a technology related message",
		Language:      "en",
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
	_, err = app.Repository.Messages.Create(ctx, message3)
	if err != nil {
		t.Fatalf("Failed to create message3: %v", err)
	}

	// Helper function to create dynamic options
	createDynamicOpts := func(likeID string, likeMessageText string) []repo.OptionTopic {
		return []repo.OptionTopic{
			repo.TopicLikeID(likeID),
			repo.TopicLikeMessageText(likeMessageText),
		}
	}

	t.Run("CountByStatusDynamic - no options (all topics)", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusDynamic(ctx)
		if err != nil {
			t.Fatalf("CountByStatusDynamic with no options failed: %v", err)
		}
		if counts[factcheck.StatusTopicPending] != 2 {
			t.Errorf("Expected 2 pending topics, got %d", counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != 1 {
			t.Errorf("Expected 1 resolved topic, got %d", counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusDynamic - message text filter only (English)", func(t *testing.T) {
		opts := createDynamicOpts("", "COVID")
		counts, err := app.Repository.Topics.CountByStatusDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamic with English text filter failed: %v", err)
		}
		if counts[factcheck.StatusTopicPending] != 1 {
			t.Errorf("Expected 1 pending topic with 'COVID' in message, got %d", counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != 0 {
			t.Errorf("Expected 0 resolved topics with 'COVID' in message, got %d", counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusDynamic - message text filter only (Thai)", func(t *testing.T) {
		opts := createDynamicOpts("", "ข่าวปลอม")
		counts, err := app.Repository.Topics.CountByStatusDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamic with Thai text filter failed: %v", err)
		}
		if counts[factcheck.StatusTopicResolved] != 1 {
			t.Errorf("Expected 1 resolved topic with 'ข่าวปลอม' in message, got %d", counts[factcheck.StatusTopicResolved])
		}
		if counts[factcheck.StatusTopicPending] != 0 {
			t.Errorf("Expected 0 pending topics with 'ข่าวปลอม' in message, got %d", counts[factcheck.StatusTopicPending])
		}
	})

	t.Run("CountByStatusDynamic - ID filter only", func(t *testing.T) {
		opts := createDynamicOpts("550e8400", "")
		counts, err := app.Repository.Topics.CountByStatusDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamic with ID filter failed: %v", err)
		}
		if counts[factcheck.StatusTopicPending] != 1 {
			t.Errorf("Expected 1 pending topic with ID starting with '550e8400', got %d", counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != 1 {
			t.Errorf("Expected 1 resolved topic with ID starting with '550e8400', got %d", counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusDynamic - combined filters (ID + message text)", func(t *testing.T) {
		opts := createDynamicOpts("550e8400", "COVID")
		counts, err := app.Repository.Topics.CountByStatusDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamic with ID + message text filter failed: %v", err)
		}
		if counts[factcheck.StatusTopicPending] != 1 {
			t.Errorf("Expected 1 pending topic (550e8400 + COVID), got %d", counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != 0 {
			t.Errorf("Expected 0 resolved topics, got %d", counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusDynamic - empty results", func(t *testing.T) {
		opts := createDynamicOpts("", "nonexistent")
		counts, err := app.Repository.Topics.CountByStatusDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamic with no matching results failed: %v", err)
		}
		if counts[factcheck.StatusTopicResolved] != 0 {
			t.Errorf("Expected 0 resolved topics with 'nonexistent' in message, got %d", counts[factcheck.StatusTopicResolved])
		}
		if counts[factcheck.StatusTopicPending] != 0 {
			t.Errorf("Expected 0 pending topics, got %d", counts[factcheck.StatusTopicPending])
		}
	})
}

func TestRepository_CountByStatusDynamicV2(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()
	ctx := t.Context()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create test topics
	topic1 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440001",
		Name:        "Topic 1 - COVID Pending",
		Description: "COVID-19 related topic",
		Status:      factcheck.StatusTopicPending,
		CreatedAt:   now,
		UpdatedAt:   nil,
	}
	topic2 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440002",
		Name:        "Topic 2 - Politics Resolved",
		Description: "Politics related topic",
		Status:      factcheck.StatusTopicResolved,
		CreatedAt:   now,
		UpdatedAt:   nil,
	}
	topic3 := factcheck.Topic{
		ID:          "660e8400-e29b-41d4-a716-446655440003",
		Name:        "Topic 3 - Technology Pending",
		Description: "Technology related topic",
		Status:      factcheck.StatusTopicPending,
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

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

	// Create message groups with different languages
	messageGroup1 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440001",
		TopicID:   createdTopic1.ID,
		Name:      "COVID Vaccine Group",
		Text:      "COVID-19 vaccine is effective against new variants",
		TextSHA1:  "sha1_hash_1",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	messageGroup2 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440002",
		TopicID:   createdTopic2.ID,
		Name:      "Election News Group",
		Text:      "ข่าวปลอมเกี่ยวกับการเลือกตั้ง",
		TextSHA1:  "sha1_hash_2",
		Language:  factcheck.LanguageThai,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	messageGroup3 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440003",
		TopicID:   createdTopic3.ID,
		Name:      "AI Technology Group",
		Text:      "New AI technology breakthrough",
		TextSHA1:  "sha1_hash_3",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	// Create message groups in database
	_, err = app.Repository.MessageGroups.Create(ctx, messageGroup1)
	if err != nil {
		t.Fatalf("Failed to create messageGroup1: %v", err)
	}

	_, err = app.Repository.MessageGroups.Create(ctx, messageGroup2)
	if err != nil {
		t.Fatalf("Failed to create messageGroup2: %v", err)
	}

	_, err = app.Repository.MessageGroups.Create(ctx, messageGroup3)
	if err != nil {
		t.Fatalf("Failed to create messageGroup3: %v", err)
	}

	// Helper function to create dynamic options
	createDynamicOpts := func(likeID string, likeMessageText string) []repo.OptionTopic {
		return []repo.OptionTopic{
			repo.TopicLikeID(likeID),
			repo.TopicLikeMessageText(likeMessageText),
		}
	}

	t.Run("CountByStatusDynamicV2 - no options (all topics)", func(t *testing.T) {
		counts, err := app.Repository.Topics.CountByStatusDynamicV2(ctx)
		if err != nil {
			t.Fatalf("CountByStatusDynamicV2 with no options failed: %v", err)
		}
		if counts[factcheck.StatusTopicPending] != 2 {
			t.Errorf("Expected 2 pending topics, got %d", counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != 1 {
			t.Errorf("Expected 1 resolved topic, got %d", counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusDynamicV2 - message group text filter only (English)", func(t *testing.T) {
		opts := createDynamicOpts("", "COVID")
		counts, err := app.Repository.Topics.CountByStatusDynamicV2(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamicV2 with English text filter failed: %v", err)
		}
		if counts[factcheck.StatusTopicPending] != 1 {
			t.Errorf("Expected 1 pending topic with 'COVID' in message group, got %d", counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != 0 {
			t.Errorf("Expected 0 resolved topics with 'COVID' in message group, got %d", counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusDynamicV2 - message group text filter only (Thai)", func(t *testing.T) {
		opts := createDynamicOpts("", "ข่าวปลอม")
		counts, err := app.Repository.Topics.CountByStatusDynamicV2(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamicV2 with Thai text filter failed: %v", err)
		}
		if counts[factcheck.StatusTopicResolved] != 1 {
			t.Errorf("Expected 1 resolved topic with 'ข่าวปลอม' in message group, got %d", counts[factcheck.StatusTopicResolved])
		}
		if counts[factcheck.StatusTopicPending] != 0 {
			t.Errorf("Expected 0 pending topics with 'ข่าวปลอม' in message group, got %d", counts[factcheck.StatusTopicPending])
		}
	})

	t.Run("CountByStatusDynamicV2 - ID filter only", func(t *testing.T) {
		opts := createDynamicOpts("550e8400", "")
		counts, err := app.Repository.Topics.CountByStatusDynamicV2(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamicV2 with ID filter failed: %v", err)
		}
		if counts[factcheck.StatusTopicPending] != 1 {
			t.Errorf("Expected 1 pending topic with ID starting with '550e8400', got %d", counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != 1 {
			t.Errorf("Expected 1 resolved topic with ID starting with '550e8400', got %d", counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusDynamicV2 - combined filters (ID + message group text)", func(t *testing.T) {
		opts := createDynamicOpts("550e8400", "COVID")
		counts, err := app.Repository.Topics.CountByStatusDynamicV2(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamicV2 with ID + message group text filter failed: %v", err)
		}
		if counts[factcheck.StatusTopicPending] != 1 {
			t.Errorf("Expected 1 pending topic (550e8400 + COVID), got %d", counts[factcheck.StatusTopicPending])
		}
		if counts[factcheck.StatusTopicResolved] != 0 {
			t.Errorf("Expected 0 resolved topics, got %d", counts[factcheck.StatusTopicResolved])
		}
	})

	t.Run("CountByStatusDynamicV2 - empty results", func(t *testing.T) {
		opts := createDynamicOpts("", "nonexistent")
		counts, err := app.Repository.Topics.CountByStatusDynamicV2(ctx, opts...)
		if err != nil {
			t.Fatalf("CountByStatusDynamicV2 with no matching results failed: %v", err)
		}
		if counts[factcheck.StatusTopicResolved] != 0 {
			t.Errorf("Expected 0 resolved topics with 'nonexistent' in message group, got %d", counts[factcheck.StatusTopicResolved])
		}
		if counts[factcheck.StatusTopicPending] != 0 {
			t.Errorf("Expected 0 pending topics, got %d", counts[factcheck.StatusTopicPending])
		}
	})
}

func TestRepository_ListDynamicV2(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()
	ctx := t.Context()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create topics with different statuses
	topic1 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440001",
		Name:        "Topic 1 - COVID Pending",
		Description: "COVID-19 related news (pending)",
		Status:      factcheck.StatusTopicPending,
		Result:      "",
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

	topic2 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440002",
		Name:        "Topic 2 - Politics Resolved",
		Description: "Political news and updates (resolved)",
		Status:      factcheck.StatusTopicResolved,
		Result:      "Verified as true",
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

	topic3 := factcheck.Topic{
		ID:          "660e8400-e29b-41d4-a716-446655440003",
		Name:        "Topic 3 - Technology Pending",
		Description: "Technology news and updates (pending)",
		Status:      factcheck.StatusTopicPending,
		Result:      "",
		CreatedAt:   now,
		UpdatedAt:   nil,
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

	// Create message groups with different languages
	messageGroup1 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440001",
		TopicID:   createdTopic1.ID,
		Name:      "COVID Vaccine Group",
		Text:      "COVID-19 vaccine is effective against new variants",
		TextSHA1:  "sha1_hash_1",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	messageGroup2 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440002",
		TopicID:   createdTopic2.ID,
		Name:      "Election News Group",
		Text:      "ข่าวปลอมเกี่ยวกับการเลือกตั้ง",
		TextSHA1:  "sha1_hash_2",
		Language:  factcheck.LanguageThai,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	messageGroup3 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440003",
		TopicID:   createdTopic3.ID,
		Name:      "AI Technology Group",
		Text:      "New AI technology breakthrough",
		TextSHA1:  "sha1_hash_3",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	// Create message groups in database
	_, err = app.Repository.MessageGroups.Create(ctx, messageGroup1)
	if err != nil {
		t.Fatalf("Failed to create messageGroup1: %v", err)
	}

	_, err = app.Repository.MessageGroups.Create(ctx, messageGroup2)
	if err != nil {
		t.Fatalf("Failed to create messageGroup2: %v", err)
	}

	_, err = app.Repository.MessageGroups.Create(ctx, messageGroup3)
	if err != nil {
		t.Fatalf("Failed to create messageGroup3: %v", err)
	}

	// Helper function to create dynamic options
	createDynamicOpts := func(likeID string, statuses []factcheck.StatusTopic, likeMessageText string) []repo.OptionTopic {
		return []repo.OptionTopic{
			repo.TopicLikeID(likeID),
			repo.TopicInStatuses(statuses),
			repo.TopicLikeMessageText(likeMessageText),
		}
	}

	t.Run("ListDynamicV2 - no options (all topics)", func(t *testing.T) {
		ctx := t.Context()
		// Debug: Check if topics exist in database
		allTopics, err := app.Repository.Topics.List(ctx, 0, 0)
		if err != nil {
			t.Fatalf("Failed to list all topics: %v", err)
		}
		t.Logf("Total topics in database: %d", len(allTopics))
		for i, topic := range allTopics {
			t.Logf("Topic %d: ID=%s, Name=%s, Status=%s", i+1, topic.ID, topic.Name, topic.Status)
		}

		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0)
		if err != nil {
			t.Fatalf("ListDynamicV2 with no options failed: %v", err)
		}
		t.Logf("ListDynamicV2 returned %d topics", len(topics))
		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics, got %d", len(topics))
		}
	})

	t.Run("ListDynamicV2 - status filter only (single status)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicPending}, "")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with status filter failed: %v", err)
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

	t.Run("ListDynamicV2 - status filter only (multiple statuses)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicPending, factcheck.StatusTopicResolved}, "")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with multiple statuses failed: %v", err)
		}
		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics (pending + resolved), got %d", len(topics))
		}
	})

	t.Run("ListDynamicV2 - message group text filter only (English)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", nil, "COVID")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with English text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with 'COVID' in message group, got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamicV2 - message group text filter only (Thai)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", nil, "ข่าวปลอม")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with Thai text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with 'ข่าวปลอม' in message group, got %d", len(topics))
		}
		if topics[0].ID != createdTopic2.ID {
			t.Errorf("Expected topic2 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamicV2 - ID filter only", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("550e8400", nil, "")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with ID filter failed: %v", err)
		}
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics with ID starting with '550e8400', got %d", len(topics))
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
			t.Errorf("Expected topic3 to NOT be in results")
		}
	})

	t.Run("ListDynamicV2 - combined filters (status + message group text)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicPending}, "COVID")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with status + message group text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic (pending + COVID), got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamicV2 - combined filters (ID + message group text)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("550e8400", nil, "COVID")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with ID + message group text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic (550e8400 + COVID), got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamicV2 - all filters combined", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("550e8400", []factcheck.StatusTopic{factcheck.StatusTopicPending}, "COVID")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with all filters failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic (550e8400 + pending + COVID), got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamicV2 - pagination", func(t *testing.T) {
		ctx := t.Context()
		// Test limit
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 2, 0)
		if err != nil {
			t.Fatalf("ListDynamicV2 with limit failed: %v", err)
		}
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics with limit=2, got %d", len(topics))
		}

		// Test offset
		topics, err = app.Repository.Topics.ListDynamicV2(ctx, 1, 1)
		if err != nil {
			t.Fatalf("ListDynamicV2 with offset failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with limit=1, offset=1, got %d", len(topics))
		}
	})

	t.Run("ListDynamicV2 - empty results", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicResolved}, "nonexistent")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with no matching results failed: %v", err)
		}
		if len(topics) != 0 {
			t.Fatalf("Expected 0 topics, got %d", len(topics))
		}
	})

	t.Run("ListDynamicV2 - no pagination (limit=0, offset=0)", func(t *testing.T) {
		ctx := t.Context()
		// Test that when both limit and offset are 0, we get all results without pagination
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0)
		if err != nil {
			t.Fatalf("ListDynamicV2 with no pagination failed: %v", err)
		}
		if len(topics) != 3 {
			t.Fatalf("Expected all 3 topics with no pagination, got %d", len(topics))
		}

		// Verify we got all topics
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
		if !topicIDs[createdTopic3.ID] {
			t.Errorf("Expected topic3 to be in results")
		}
	})

	t.Run("ListDynamicV2 - case insensitive message group text filter", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", nil, "covid")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with case insensitive text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with 'covid' (case insensitive) in message group, got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamicV2 - partial message group text match", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", nil, "technology")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with partial text match failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with 'technology' in message group, got %d", len(topics))
		}
		if topics[0].ID != createdTopic3.ID {
			t.Errorf("Expected topic3 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamicV2 - ID pattern with wildcards", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("550e8400%", nil, "")
		topics, err := app.Repository.Topics.ListDynamicV2(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamicV2 with ID pattern failed: %v", err)
		}
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics with ID pattern '550e8400%%', got %d", len(topics))
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
			t.Errorf("Expected topic3 to NOT be in results")
		}
	})
}

func TestRepository_TopicFiltering(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()
	ctx := t.Context()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create topics
	topic1 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440001",
		Name:        "Topic 1 - COVID",
		Description: "COVID-19 related news",
		Status:      factcheck.StatusTopicPending,
		Result:      "",
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

	topic2 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440002",
		Name:        "Topic 2 - Politics",
		Description: "Political news and updates",
		Status:      factcheck.StatusTopicResolved,
		Result:      "Verified as true",
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

	topic3 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440003",
		Name:        "Topic 3 - Technology",
		Description: "Technology news and updates",
		Status:      factcheck.StatusTopicPending,
		Result:      "",
		CreatedAt:   now,
		UpdatedAt:   nil,
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
		ctx := t.Context()
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
}

func TestRepository_ListDynamic(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()
	ctx := t.Context()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create topics with different statuses
	topic1 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440001",
		Name:        "Topic 1 - COVID Pending",
		Description: "COVID-19 related news (pending)",
		Status:      factcheck.StatusTopicPending,
		Result:      "",
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

	topic2 := factcheck.Topic{
		ID:          "550e8400-e29b-41d4-a716-446655440002",
		Name:        "Topic 2 - Politics Resolved",
		Description: "Political news and updates (resolved)",
		Status:      factcheck.StatusTopicResolved,
		Result:      "Verified as true",
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

	topic3 := factcheck.Topic{
		ID:          "660e8400-e29b-41d4-a716-446655440003",
		Name:        "Topic 3 - Technology Pending",
		Description: "Technology news and updates (pending)",
		Status:      factcheck.StatusTopicPending,
		Result:      "",
		CreatedAt:   now,
		UpdatedAt:   nil,
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

	// Create messages with different languages
	message1 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440001",
		UserMessageID: createdUserMessage1.ID,
		TopicID:       createdTopic1.ID,
		Text:          "COVID-19 vaccine is effective against new variants",
		Language:      factcheck.LanguageEnglish,
		Type:          factcheck.TypeMessageText,
		Status:        factcheck.StatusMessageTopicSubmitted,
		CreatedAt:     now,
		UpdatedAt:     nil,
	}

	message2 := factcheck.Message{
		ID:            "660e8400-e29b-41d4-a716-446655440002",
		UserMessageID: createdUserMessage2.ID,
		TopicID:       createdTopic2.ID,
		Text:          "ข่าวปลอมเกี่ยวกับการเลือกตั้ง",
		Language:      factcheck.LanguageThai,
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
		Language:      factcheck.LanguageEnglish,
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

	// Helper function to create dynamic options
	createDynamicOpts := func(likeID string, statuses []factcheck.StatusTopic, likeMessageText string) []repo.OptionTopic {
		return []repo.OptionTopic{
			repo.TopicLikeID(likeID),
			repo.TopicInStatuses(statuses),
			repo.TopicLikeMessageText(likeMessageText),
		}
	}

	t.Run("ListDynamic - no options (all topics)", func(t *testing.T) {
		ctx := t.Context()
		// Debug: Check if topics exist in database
		allTopics, err := app.Repository.Topics.List(ctx, 0, 0)
		if err != nil {
			t.Fatalf("Failed to list all topics: %v", err)
		}
		t.Logf("Total topics in database: %d", len(allTopics))
		for i, topic := range allTopics {
			t.Logf("Topic %d: ID=%s, Name=%s, Status=%s", i+1, topic.ID, topic.Name, topic.Status)
		}

		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0)
		if err != nil {
			t.Fatalf("ListDynamic with no options failed: %v", err)
		}
		t.Logf("ListDynamic returned %d topics", len(topics))
		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics, got %d", len(topics))
		}
	})

	t.Run("ListDynamic - status filter only (single status)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicPending}, "")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with status filter failed: %v", err)
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

	t.Run("ListDynamic - status filter only (multiple statuses)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicPending, factcheck.StatusTopicResolved}, "")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with multiple statuses failed: %v", err)
		}
		if len(topics) != 3 {
			t.Fatalf("Expected 3 topics (pending + resolved), got %d", len(topics))
		}
	})

	t.Run("ListDynamic - message text filter only (English)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", nil, "COVID")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with English text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with 'COVID' in message, got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamic - message text filter only (Thai)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", nil, "ข่าวปลอม")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with Thai text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with 'ข่าวปลอม' in message, got %d", len(topics))
		}
		if topics[0].ID != createdTopic2.ID {
			t.Errorf("Expected topic2 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamic - ID filter only", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("550e8400", nil, "")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with ID filter failed: %v", err)
		}
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics with ID starting with '550e8400', got %d", len(topics))
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
			t.Errorf("Expected topic3 to NOT be in results")
		}
	})

	t.Run("ListDynamic - combined filters (status + message text)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicPending}, "COVID")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with status + message text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic (pending + COVID), got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamic - combined filters (ID + message text)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("550e8400", nil, "COVID")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with ID + message text filter failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic (550e8400 + COVID), got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamic - all filters combined", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("550e8400", []factcheck.StatusTopic{factcheck.StatusTopicPending}, "COVID")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with all filters failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic (550e8400 + pending + COVID), got %d", len(topics))
		}
		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 to be returned, got %s", topics[0].ID)
		}
	})

	t.Run("ListDynamic - pagination", func(t *testing.T) {
		ctx := t.Context()
		// Test limit
		topics, err := app.Repository.Topics.ListDynamic(ctx, 2, 0)
		if err != nil {
			t.Fatalf("ListDynamic with limit failed: %v", err)
		}
		if len(topics) != 2 {
			t.Fatalf("Expected 2 topics with limit=2, got %d", len(topics))
		}

		// Test offset
		topics, err = app.Repository.Topics.ListDynamic(ctx, 1, 1)
		if err != nil {
			t.Fatalf("ListDynamic with offset failed: %v", err)
		}
		if len(topics) != 1 {
			t.Fatalf("Expected 1 topic with limit=1, offset=1, got %d", len(topics))
		}
	})

	t.Run("ListDynamic - empty results", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicResolved}, "nonexistent")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with no matching results failed: %v", err)
		}
		if len(topics) != 0 {
			t.Fatalf("Expected 0 topics, got %d", len(topics))
		}
	})

	t.Run("ListDynamic - no pagination (limit=0, offset=0)", func(t *testing.T) {
		ctx := t.Context()
		// Test that when both limit and offset are 0, we get all results without pagination
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0)
		if err != nil {
			t.Fatalf("ListDynamic with no pagination failed: %v", err)
		}
		if len(topics) != 3 {
			t.Fatalf("Expected all 3 topics with no pagination, got %d", len(topics))
		}

		// Verify we got all topics
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
		if !topicIDs[createdTopic3.ID] {
			t.Errorf("Expected topic3 to be in results")
		}
	})

	t.Run("ListDynamic - no pagination with filters (limit=0, offset=0)", func(t *testing.T) {
		ctx := t.Context()
		// Test no pagination with filters
		opts := createDynamicOpts("", []factcheck.StatusTopic{factcheck.StatusTopicPending}, "")
		topics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with no pagination and filters failed: %v", err)
		}
		if len(topics) != 2 {
			t.Fatalf("Expected 2 pending topics with no pagination, got %d", len(topics))
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

	t.Run("ListDynamic - pagination vs no pagination comparison", func(t *testing.T) {
		ctx := t.Context()
		// Test that pagination works differently from no pagination
		//
		// With pagination (limit=2, offset=0)
		paginatedTopics, err := app.Repository.Topics.ListDynamic(ctx, 2, 0)
		if err != nil {
			t.Fatalf("ListDynamic with pagination failed: %v", err)
		}
		if len(paginatedTopics) != 2 {
			t.Fatalf("Expected 2 topics with pagination limit=2, got %d", len(paginatedTopics))
		}

		// Without pagination (limit=0, offset=0)
		allTopics, err := app.Repository.Topics.ListDynamic(ctx, 0, 0)
		if err != nil {
			t.Fatalf("ListDynamic with no pagination failed: %v", err)
		}
		if len(allTopics) != 3 {
			t.Fatalf("Expected all 3 topics with no pagination, got %d", len(allTopics))
		}

		// Verify that paginated results are a subset of all results
		paginatedIDs := make(map[string]bool)
		for _, topic := range paginatedTopics {
			paginatedIDs[topic.ID] = true
		}

		allIDs := make(map[string]bool)
		for _, topic := range allTopics {
			allIDs[topic.ID] = true
		}

		// All paginated topics should be in all topics
		for id := range paginatedIDs {
			if !allIDs[id] {
				t.Errorf("Paginated topic ID %s not found in all topics", id)
			}
		}
	})
}
