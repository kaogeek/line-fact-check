//go:build integration_test
// +build integration_test

package repo_test

import (
	"testing"
	"time"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	// "github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestMessageGroupRepository_ListDynamic(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()
	ctx := t.Context()

	// Create test data
	baseTime := time.Now()

	// Create message groups
	messageGroup1 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440001",
		Status:    factcheck.StatusMGroupPending,
		TopicID:   "",
		Name:      "COVID Vaccine Group",
		Text:      "COVID-19 vaccine is effective against new variants",
		TextSHA1:  "sha1_hash_1",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: baseTime,
		UpdatedAt: nil,
	}

	messageGroup2 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440002",
		Status:    factcheck.StatusMGroupApproved,
		TopicID:   "",
		Name:      "Election Results Group",
		Text:      "Election results show clear victory",
		TextSHA1:  "sha1_hash_2",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: baseTime.Add(1 * time.Millisecond),
		UpdatedAt: nil,
	}

	messageGroup3 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440003",
		Status:    factcheck.StatusMGroupRejected,
		TopicID:   "",
		Name:      "AI Technology Group",
		Text:      "New AI technology breakthrough",
		TextSHA1:  "sha1_hash_3",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: baseTime.Add(2 * time.Millisecond),
		UpdatedAt: nil,
	}

	messageGroup4 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440004",
		Status:    factcheck.StatusMGroupPending,
		TopicID:   "",
		Name:      "Sports Results Group",
		Text:      "World Cup final results announced",
		TextSHA1:  "sha1_hash_4",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: baseTime.Add(3 * time.Millisecond),
		UpdatedAt: nil,
	}

	// Create topics in database
	createdMessageGroup1, err := app.Repository.MessageGroups.Create(ctx, messageGroup1)
	if err != nil {
		t.Fatalf("Failed to create messageGroup1: %v", err)
	}

	createdMessageGroup2, err := app.Repository.MessageGroups.Create(ctx, messageGroup2)
	if err != nil {
		t.Fatalf("Failed to create messageGroup2: %v", err)
	}

	createdMessageGroup3, err := app.Repository.MessageGroups.Create(ctx, messageGroup3)
	if err != nil {
		t.Fatalf("Failed to create messageGroup3: %v", err)
	}

	createdMessageGroup4, err := app.Repository.MessageGroups.Create(ctx, messageGroup4)
	if err != nil {
		t.Fatalf("Failed to create messageGroup4: %v", err)
	}

	createdMessageGroups := []factcheck.MessageGroup{
		createdMessageGroup1,
		createdMessageGroup2,
		createdMessageGroup3,
		createdMessageGroup4,
	}

	// Helper function to create dynamic options
	createDynamicOpts := func(likeMessageText string, idIn []string, idNotIn []string, statuses []factcheck.StatusMGroup) []repo.OptionMessageGroup {
		return []repo.OptionMessageGroup{
			repo.MessageGroupLikeMessageText(likeMessageText),
			repo.MessageGroupIDIn(idIn),
			repo.MessageGroupIDNotIn(idNotIn),
			repo.MessageGroupStatusesIn(statuses),
		}
	}

	t.Run("MessageGroupListDynamic - no options (all message groups)", func(t *testing.T) {
		ctx := t.Context()
		// Debug: Check if topics exist in database
		allMgs, err := app.Repository.MessageGroups.ListDynamic(ctx, 0, 0)
		if err != nil {
			t.Fatalf("Failed to list all message groups: %v", err)
		}

		t.Logf("Total message groups in database: %d", len(allMgs))

		for i, messageGroup := range allMgs {
			t.Logf("Message Group %d: ID=%s, Name=%s, Status=%s", i+1, messageGroup.ID, messageGroup.Name, messageGroup.Status)
		}

		if len(allMgs) != 4 {
			t.Fatalf("Expected 4 message groups, got %d", len(allMgs))
		}
	})

	t.Run("MessageGroupListDynamic - message group text filter", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("New AI technology break", []string{}, []string{}, []factcheck.StatusMGroup{})
		mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with English text filter failed: %v", err)
		}
		if len(mgs) != 1 {
			t.Fatalf("Expected 1 message group with 'New AI technology break' in message group, got %d", len(mgs))
		}
		if mgs[0].ID != createdMessageGroup3.ID {
			t.Errorf("Expected message group3 to be returned, got %s", mgs[0].ID)
		}
	})

	t.Run("MessageGroupListDynamic - message group text filter (not found)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("Unknown key word", []string{}, []string{}, []factcheck.StatusMGroup{})
		mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with English text filter failed: %v", err)
		}
		if len(mgs) != 0 {
			t.Fatalf("Expected 0 message group with 'Unknown key word' in message group, got %d", len(mgs))
		}
	})

	t.Run("MessageGroupListDynamic - filter by ID (inclusion)", func(t *testing.T) {
		ctx := t.Context()
		// Filter for specific message group IDs
		opts := createDynamicOpts("", []string{createdMessageGroup1.ID, createdMessageGroup2.ID}, []string{}, []factcheck.StatusMGroup{})
		mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with ID filter failed: %v", err)
		}
		if len(mgs) != 2 {
			t.Fatalf("Expected 2 message groups with specified IDs, got %d", len(mgs))
		}

		// Verify we got exactly the message groups we asked for
		foundIDs := make(map[string]bool)
		for _, mg := range mgs {
			foundIDs[mg.ID] = true
		}

		if !foundIDs[createdMessageGroup1.ID] || !foundIDs[createdMessageGroup2.ID] {
			t.Error("Expected to find both specified message groups in results")
		}
	})

	t.Run("MessageGroupListDynamic - filter by ID (exclusion)", func(t *testing.T) {
		ctx := t.Context()
		// Exclude two message groups
		opts := createDynamicOpts("", []string{}, []string{createdMessageGroup1.ID, createdMessageGroup2.ID}, []factcheck.StatusMGroup{})
		mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with ID not in filter failed: %v", err)
		}

		// We have 4 total message groups, excluding 2 should leave us with 2
		if len(mgs) != 2 {
			t.Fatalf("Expected 2 message groups after exclusion, got %d", len(mgs))
		}

		// Verify the excluded message groups are not in the results
		for _, mg := range mgs {
			if mg.ID == createdMessageGroup1.ID || mg.ID == createdMessageGroup2.ID {
				t.Errorf("Found excluded message group ID in results: %s", mg.ID)
			}
		}
	})

	t.Run("MessageGroupListDynamic - combine text search with ID exclusion", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("COVID Vaccine", []string{}, []string{createdMessageGroup1.ID}, []factcheck.StatusMGroup{})
		mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with text search and ID exclusion failed: %v", err)
		}

		if len(mgs) != 0 {
			t.Fatalf("Expected 0 message groups with 'vaccine' after exclusion, got %d", len(mgs))
		}

		opts = createDynamicOpts("show clear victory", []string{}, []string{createdMessageGroup3.ID}, []factcheck.StatusMGroup{})
		mgs, err = app.Repository.MessageGroups.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with text search and combined ID filters failed: %v", err)
		}

		if len(mgs) != 1 {
			t.Fatalf("Expected 1 message groups with 'show clear victory' after exclusion, got %d", len(mgs))
		}

		if mgs[0].ID != createdMessageGroup2.ID {
			t.Errorf("Expected message group2 to be returned, got %s", mgs[0].ID)
		}
	})

	t.Run("MessageGroupListDynamic - filter by Status (inclusion)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("", []string{}, []string{}, []factcheck.StatusMGroup{factcheck.StatusMGroupPending, factcheck.StatusMGroupRejected})
		mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with status filter failed: %v", err)
		}
		if len(mgs) != 3 {
			t.Fatalf("Expected 3 message group with 'pending' status, got %d", len(mgs))
		}

		for _, mg := range mgs {
			if mg.Status != factcheck.StatusMGroupPending && mg.Status != factcheck.StatusMGroupRejected {
				t.Errorf("Expected message group with 'pending' or 'rejected' status, got %s", mg.Status)
			}
		}
	})

	t.Run("MessageGroupListDynamic - pagination", func(t *testing.T) {
		ctx := t.Context()
		limit := 1

		for i := 0; i < len(createdMessageGroups); i++ {
			offset := i * limit

			mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, limit, offset)
			if err != nil {
				t.Fatalf("ListDynamic with pagination failed: %v", err)
			}
			if len(mgs) != 1 {
				t.Fatalf("Expected 1 message group with page %d in message group, got %d", i+1, len(mgs))
			}

			// use reverseIndex for message createMg cause by query is order createAt desc
			reverseIndex := len(createdMessageGroups) - 1 - i
			if mgs[0].ID != createdMessageGroups[reverseIndex].ID {
				t.Fatalf("Expected message group %d to be returned, got %s", i+1, mgs[0].ID)
			}
		}
	})
}

func TestMessageGroupRepository_CountDynamic(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()
	ctx := t.Context()

	// Create test data
	baseTime := time.Now()

	// Create message groups
	messageGroup1 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440001",
		Status:    factcheck.StatusMGroupPending,
		TopicID:   "",
		Name:      "COVID Vaccine Group",
		Text:      "COVID-19 vaccine is effective against new variants",
		TextSHA1:  "sha1_hash_1",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: baseTime,
		UpdatedAt: nil,
	}

	messageGroup2 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440002",
		Status:    factcheck.StatusMGroupApproved,
		TopicID:   "",
		Name:      "Election Results Group",
		Text:      "Election results show clear victory",
		TextSHA1:  "sha1_hash_2",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: baseTime.Add(1 * time.Millisecond),
		UpdatedAt: nil,
	}

	messageGroup3 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440003",
		Status:    factcheck.StatusMGroupRejected,
		TopicID:   "",
		Name:      "AI Technology Group",
		Text:      "New AI technology breakthrough",
		TextSHA1:  "sha1_hash_3",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: baseTime.Add(2 * time.Millisecond),
		UpdatedAt: nil,
	}

	messageGroup4 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440004",
		Status:    factcheck.StatusMGroupPending,
		TopicID:   "",
		Name:      "Sports Results Group",
		Text:      "World Cup final results announced",
		TextSHA1:  "sha1_hash_4",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: baseTime.Add(3 * time.Millisecond),
		UpdatedAt: nil,
	}

	// Create topics in database
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

	_, err = app.Repository.MessageGroups.Create(ctx, messageGroup4)
	if err != nil {
		t.Fatalf("Failed to create messageGroup4: %v", err)
	}

	// Helper function to create dynamic options
	createDynamicOpts := func(likeMessageText string, idIn []string, idNotIn []string) []repo.OptionMessageGroup {
		return []repo.OptionMessageGroup{
			repo.MessageGroupLikeMessageText(likeMessageText),
			repo.MessageGroupIDIn(idIn),
			repo.MessageGroupIDNotIn(idNotIn),
		}
	}

	t.Run("MessageGroupCountDynamic - no options (all message groups)", func(t *testing.T) {
		ctx := t.Context()
		count, err := app.Repository.MessageGroups.CountDynamic(ctx)
		if err != nil {
			t.Fatalf("Failed to count all message groups: %v", err)
		}

		t.Logf("Total message groups in database: %v", count)

		if count[factcheck.StatusMGroupPending] != 2 {
			t.Fatalf("Expected 1 message group with 'pending' status, got %d", count[factcheck.StatusMGroupPending])
		}
		if count[factcheck.StatusMGroupApproved] != 1 {
			t.Fatalf("Expected 1 message group with 'approved' status, got %d", count[factcheck.StatusMGroupPending])
		}
		if count[factcheck.StatusMGroupRejected] != 1 {
			t.Fatalf("Expected 1 message group with 'rejected' status, got %d", count[factcheck.StatusMGroupPending])
		}
	})

	t.Run("MessageGroupCountDynamic - message group text filter", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("New AI technology break", []string{}, []string{})
		count, err := app.Repository.MessageGroups.CountDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountDynamic with text filter failed: %v", err)
		}
		total := count[factcheck.StatusMGroupPending] + count[factcheck.StatusMGroupApproved] + count[factcheck.StatusMGroupRejected]
		if total != 1 {
			t.Fatalf("Expected 1 message group with 'New AI technology break' in text, got %d", total)
		}
	})

	t.Run("MessageGroupCountDynamic - message group text filter (not found)", func(t *testing.T) {
		ctx := t.Context()
		opts := createDynamicOpts("Nonexistent text that won't be found", []string{}, []string{})
		count, err := app.Repository.MessageGroups.CountDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountDynamic with text filter failed: %v", err)
		}
		total := count[factcheck.StatusMGroupPending] + count[factcheck.StatusMGroupApproved] + count[factcheck.StatusMGroupRejected]
		if total != 0 {
			t.Fatalf("Expected 0 message groups with 'Nonexistent text' in text, got %d", total)
		}
	})

	t.Run("MessageGroupCountDynamic - filter by ID (inclusion)", func(t *testing.T) {
		ctx := t.Context()
		// Get all message groups to get their IDs
		_, err := app.Repository.MessageGroups.CountDynamic(ctx)
		if err != nil {
			t.Fatalf("Failed to get all message groups: %v", err)
		}

		// Get two message groups to test inclusion
		opts := createDynamicOpts("", []string{messageGroup1.ID, messageGroup2.ID}, []string{})
		count, err := app.Repository.MessageGroups.CountDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountDynamic with ID inclusion filter failed: %v", err)
		}
		total := count[factcheck.StatusMGroupPending] + count[factcheck.StatusMGroupApproved] + count[factcheck.StatusMGroupRejected]
		if total != 2 {
			t.Fatalf("Expected 2 message groups with specified IDs, got %d", total)
		}
	})

	t.Run("MessageGroupCountDynamic - filter by ID (exclusion)", func(t *testing.T) {
		ctx := t.Context()
		// Exclude two message groups
		opts := createDynamicOpts("", []string{}, []string{messageGroup1.ID, messageGroup2.ID})
		count, err := app.Repository.MessageGroups.CountDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountDynamic with ID exclusion filter failed: %v", err)
		}

		total := count[factcheck.StatusMGroupPending] + count[factcheck.StatusMGroupApproved] + count[factcheck.StatusMGroupRejected]
		if total != 2 {
			t.Fatalf("Expected 2 message groups after exclusion, got %d", total)
		}
	})

	t.Run("MessageGroupCountDynamic - combine text search with ID exclusion", func(t *testing.T) {
		ctx := t.Context()
		// First test with text search that should match one message group
		opts := createDynamicOpts("COVID Vaccine", []string{}, []string{messageGroup1.ID})
		count, err := app.Repository.MessageGroups.CountDynamic(ctx, opts...)
		if err != nil {
			t.Fatalf("CountDynamic with text search and ID exclusion failed: %v", err)
		}

		total := count[factcheck.StatusMGroupPending] + count[factcheck.StatusMGroupApproved] + count[factcheck.StatusMGroupRejected]
		if total != 0 {
			t.Fatalf("Expected 0 message groups with 'COVID Vaccine' after exclusion, got %d", total)
		}
	})
}
