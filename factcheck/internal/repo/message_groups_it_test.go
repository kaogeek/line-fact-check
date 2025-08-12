//go:build integration_test
// +build integration_test

package repo_test

import (
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestMessageGroupRepository_ListDynamic(t *testing.T) {
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

	// Create message groups
	messageGroup1 := factcheck.MessageGroup{
		ID:        "880e8400-e29b-41d4-a716-446655440001",
		Status:    factcheck.StatusMGroupPending,
		TopicID:   "",
		Name:      "COVID Vaccine Group",
		Text:      "COVID-19 vaccine is effective against new variants",
		TextSHA1:  "sha1_hash_1",
		Language:  factcheck.LanguageEnglish,
		CreatedAt: now,
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
		CreatedAt: now,
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
		CreatedAt: now,
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
		CreatedAt: now,
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
	createDynamicOpts := func(likeMessageText string, idIn []string, idNotIn []string) []repo.OptionMessageGroup {
		return []repo.OptionMessageGroup{
			repo.MessageGroupLikeMessageText(likeMessageText, idIn, idNotIn),
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
		opts := createDynamicOpts("New AI technology break", []string{}, []string{})
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
		opts := createDynamicOpts("Unknown key word", []string{}, []string{})
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
		opts := createDynamicOpts("", []string{createdMessageGroup1.ID, createdMessageGroup2.ID}, []string{})
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
		opts := createDynamicOpts("", []string{}, []string{createdMessageGroup1.ID, createdMessageGroup2.ID})
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
		opts := createDynamicOpts("COVID Vaccine", []string{}, []string{createdMessageGroup1.ID})
		mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, 0, 0, opts...)
		if err != nil {
			t.Fatalf("ListDynamic with text search and ID exclusion failed: %v", err)
		}

		if len(mgs) != 0 {
			t.Fatalf("Expected 0 message groups with 'vaccine' after exclusion, got %d", len(mgs))
		}

		opts = createDynamicOpts("show clear victory", []string{}, []string{createdMessageGroup3.ID})
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

	// TODO: resolve this later it seem like it be flaky test cause order?
	t.Run("MessageGroupListDynamic - pagination", func(t *testing.T) {
		ctx := t.Context()
		limit := 1

		for i := 0; i < 4; i++ {
			offset := i * limit

			mgs, err := app.Repository.MessageGroups.ListDynamic(ctx, limit, offset)
			if err != nil {
				t.Fatalf("ListDynamic with pagination failed: %v", err)
			}
			if len(mgs) != 1 {
				t.Fatalf("Expected 1 message group with page %d in message group, got %d", i+1, len(mgs))
			}

			if mgs[0].ID != createdMessageGroups[i].ID {
				t.Fatalf("Expected message group %d to be returned, got %s", i+1, mgs[0].ID)
			}
		}
	})
}
