//go:build integration_test
// +build integration_test

package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestHandlerMessageGroup_ListMessageGroupDynamic(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Create test server
	testServer := httptest.NewServer(app.Server.(*http.Server).Handler)
	defer testServer.Close()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create message groups with different statuses and languages
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
		Name:      "Election News Group",
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
		Name:      "Thai News Group",
		Text:      "ข่าวปลอมเกี่ยวกับการเมืองไทย",
		TextSHA1:  "sha1_hash_4",
		Language:  factcheck.LanguageThai,
		CreatedAt: now,
		UpdatedAt: nil,
	}

	// Create message groups in database
	createdMessageGroup1, err := app.Repository.MessageGroups.Create(t.Context(), messageGroup1)
	if err != nil {
		t.Fatalf("Failed to create messageGroup1: %v", err)
	}

	createdMessageGroup2, err := app.Repository.MessageGroups.Create(t.Context(), messageGroup2)
	if err != nil {
		t.Fatalf("Failed to create messageGroup2: %v", err)
	}

	createdMessageGroup3, err := app.Repository.MessageGroups.Create(t.Context(), messageGroup3)
	if err != nil {
		t.Fatalf("Failed to create messageGroup3: %v", err)
	}

	createdMessageGroup4, err := app.Repository.MessageGroups.Create(t.Context(), messageGroup4)
	if err != nil {
		t.Fatalf("Failed to create messageGroup4: %v", err)
	}

	t.Run("ListMessageGroupDynamic - no query parameters (all message groups)", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/message-groups", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 4)
	})

	t.Run("ListMessageGroupDynamic - like_message_text filter (English)", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/message-groups?like_message_text=vaccine", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 1)
		assertEq(t, messageGroups[0].ID, createdMessageGroup1.ID)
	})

	t.Run("ListMessageGroupDynamic - like_message_text filter (Thai)", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/message-groups?like_message_text=การเมือง", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 1)
		assertEq(t, messageGroups[0].ID, createdMessageGroup4.ID)
	})

	t.Run("ListMessageGroupDynamic - in_id filter", func(t *testing.T) {
		ids := messageGroup1.ID + "," + messageGroup2.ID
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/message-groups?in_id="+ids, nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 2)

		// Verify we got the expected message groups
		foundIDs := make(map[string]bool)
		for _, mg := range messageGroups {
			foundIDs[mg.ID] = true
		}
		assertEq(t, foundIDs[createdMessageGroup1.ID], true)
		assertEq(t, foundIDs[createdMessageGroup2.ID], true)
	})

	t.Run("ListMessageGroupDynamic - not_in_id filter", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/message-groups?not_in_id="+messageGroup1.ID, nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 3) // 4 total - 1 excluded = 3

		// Verify the excluded message group is not in the results
		for _, mg := range messageGroups {
			if mg.ID == createdMessageGroup1.ID {
				t.Error("Excluded message group found in results")
			}
		}
	})

	t.Run("ListMessageGroupDynamic - statuses_in filter", func(t *testing.T) {
		statuses := string(factcheck.StatusMGroupPending) + "," + string(factcheck.StatusMGroupApproved)
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/message-groups?statuses_in="+statuses, nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 3)

		// Verify we got the expected message groups
		foundIDs := make(map[string]bool)
		for _, mg := range messageGroups {
			foundIDs[mg.ID] = true
		}

		assertEq(t, foundIDs[createdMessageGroup1.ID], true)
		assertEq(t, foundIDs[createdMessageGroup2.ID], true)
		assertEq(t, foundIDs[createdMessageGroup4.ID], true)
	})

	t.Run("ListMessageGroupDynamic - combined filters", func(t *testing.T) {
		// Search for English messages but exclude messageGroup1
		req, err := http.NewRequestWithContext(
			t.Context(),
			http.MethodGet,
			testServer.URL+"/message-groups?like_message_text=technology&not_in_id="+messageGroup1.ID,
			nil,
		)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 1)
		assertEq(t, messageGroups[0].ID, createdMessageGroup3.ID)
	})

	t.Run("ListMessageGroupDynamic - pagination", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/message-groups?limit=2&offset=1", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 2)
	})

	t.Run("ListMessageGroupDynamic - no results", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/message-groups?like_message_text=nonexistent", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var messageGroups []factcheck.MessageGroup
		err = json.NewDecoder(resp.Body).Decode(&messageGroups)
		assertEq(t, err, nil)
		assertEq(t, len(messageGroups), 0)
	})
}
