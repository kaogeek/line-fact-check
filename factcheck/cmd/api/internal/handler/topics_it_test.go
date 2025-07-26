//go:build integration_test
// +build integration_test

package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestHandlerTopic_Stateful(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Create test server
	testServer := httptest.NewServer(app.Server.(*http.Server).Handler)
	defer testServer.Close()

	t.Run("CRUD 1 topic", func(t *testing.T) {
		now := utils.TimeNow().Round(0) // Postgres timestampz will not preserve monotonic clock
		utils.TimeFreeze(now)
		defer utils.TimeUnfreeze()

		name := fmt.Sprintf("topic-test-normal-%s", now.String())
		desc := fmt.Sprintf("topic-test-normal-%s-desc", now.String())

		t.Log("Testing CreateTopic")
		topic := factcheck.Topic{
			Name:        name,
			Description: desc,
		}
		body := reqBodyJSON(topic)
		reqCreate, err := http.NewRequestWithContext(t.Context(), http.MethodPost, testServer.URL+"/topics/", body)
		assertEq(t, err, nil)
		reqCreate.Header.Set("Content-Type", "application/json")
		respCreate, err := http.DefaultClient.Do(reqCreate)
		assertEq(t, err, nil)
		defer respCreate.Body.Close()
		assertEq(t, respCreate.StatusCode, http.StatusCreated)

		// Assert response
		created := factcheck.Topic{}
		err = json.NewDecoder(respCreate.Body).Decode(&created)
		assertEq(t, err, nil)
		expected := factcheck.Topic{
			ID:          created.ID,
			Name:        name,
			Description: desc,
			Status:      factcheck.StatusTopicPending,
			Result:      "",
			CreatedAt:   now,
			UpdatedAt:   nil,
		}
		assertEq(t, created, expected)

		// Assert in database
		actualDB, err := app.Repository.Topics.GetByID(t.Context(), created.ID)
		assertEq(t, err, nil)
		assertEq(t, actualDB, expected)

		t.Log("Testing ListAllTopics")
		reqList, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/all", nil)
		assertEq(t, err, nil)
		respList, err := http.DefaultClient.Do(reqList)
		assertEq(t, err, nil)
		defer respList.Body.Close()
		assertEq(t, respList.StatusCode, http.StatusOK)

		// Assert response
		actualList := []factcheck.Topic{}
		err = json.NewDecoder(respList.Body).Decode(&actualList)
		assertEq(t, err, nil)
		assertEq(t, len(actualList), 1)
		assertEq(t, actualList[0], created)

		t.Log("Testing GetTopicByID")
		reqGetByID, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/"+created.ID, nil)
		assertEq(t, err, nil)
		respGetByID, err := http.DefaultClient.Do(reqGetByID)
		assertEq(t, err, nil)
		defer respGetByID.Body.Close()
		assertEq(t, respGetByID.StatusCode, http.StatusOK)

		// Assert response
		actualGetByID := factcheck.Topic{}
		err = json.NewDecoder(respGetByID.Body).Decode(&actualGetByID)
		assertEq(t, err, nil)
		assertEq(t, actualGetByID, created)

		// Test UpdateTopicStatus
		t.Log("Testing UpdateTopicStatus")
		updateStatusBody := reqBodyJSON(struct {
			Status string `json:"status"`
		}{
			Status: string(factcheck.StatusTopicResolved),
		})
		reqUpdateStatus, err := http.NewRequestWithContext(t.Context(), http.MethodPut, testServer.URL+"/topics/"+created.ID+"/status", updateStatusBody)
		assertEq(t, err, nil)
		reqUpdateStatus.Header.Set("Content-Type", "application/json")
		respUpdateStatus, err := http.DefaultClient.Do(reqUpdateStatus)
		assertEq(t, err, nil)
		defer respUpdateStatus.Body.Close()
		assertEq(t, respUpdateStatus.StatusCode, http.StatusOK)

		// Assert UpdateTopicStatus response
		updatedStatus := factcheck.Topic{}
		err = json.NewDecoder(respUpdateStatus.Body).Decode(&updatedStatus)
		assertEq(t, err, nil)
		expectedUpdateStatus := factcheck.Topic{
			ID:          created.ID,
			Name:        name,
			Description: desc,
			Status:      factcheck.StatusTopicResolved,
			CreatedAt:   now,
			//nolint:unused
			UpdatedAt: nil, // Underlying database will set this to NOW()
		}
		assertEq(t, updatedStatus.ID, expectedUpdateStatus.ID)
		assertEq(t, updatedStatus.Name, expectedUpdateStatus.Name)
		assertEq(t, updatedStatus.Description, expectedUpdateStatus.Description)
		assertEq(t, updatedStatus.Status, expectedUpdateStatus.Status)
		assertEq(t, updatedStatus.Result, expectedUpdateStatus.Result)
		assertEq(t, updatedStatus.CreatedAt, expectedUpdateStatus.CreatedAt)
		assertNeq(t, updatedStatus.UpdatedAt, nil)
		// On fast computers, Postgres discarding monotonic clock can actually
		// make it so that updated_at is before created_at.
		// TODO: figure out how to test this.
		// assertEq(t, updatedStatus.UpdatedAt.After(updatedStatus.CreatedAt), true)

		// Verify status update in database via GetByID
		reqGetAfterStatusUpdate, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/"+created.ID, nil)
		assertEq(t, err, nil)
		respGetAfterStatusUpdate, err := http.DefaultClient.Do(reqGetAfterStatusUpdate)
		assertEq(t, err, nil)
		defer respGetAfterStatusUpdate.Body.Close()
		assertEq(t, respGetAfterStatusUpdate.StatusCode, http.StatusOK)

		actualAfterStatusUpdate := factcheck.Topic{}
		err = json.NewDecoder(respGetAfterStatusUpdate.Body).Decode(&actualAfterStatusUpdate)
		assertEq(t, err, nil)
		assertEq(t, actualAfterStatusUpdate.Status, factcheck.StatusTopicResolved)

		// Test UpdateTopicName
		t.Log("Testing UpdateTopicName")
		newName := "Updated topic name for testing"
		updateNameBody := reqBodyJSON(struct {
			Name string `json:"name"`
		}{
			Name: newName,
		})
		reqUpdateName, err := http.NewRequestWithContext(t.Context(), http.MethodPut, testServer.URL+"/topics/"+created.ID+"/name", updateNameBody)
		assertEq(t, err, nil)
		respUpdateName, err := http.DefaultClient.Do(reqUpdateName)
		assertEq(t, err, nil)
		defer respUpdateName.Body.Close()
		assertEq(t, respUpdateName.StatusCode, http.StatusOK)

		// Assert UpdateTopicName response
		updatedName := factcheck.Topic{}
		err = json.NewDecoder(respUpdateName.Body).Decode(&updatedName)
		assertEq(t, err, nil)
		assertEq(t, updatedName.Name, newName)
		assertEq(t, updatedName.Description, desc)
		assertEq(t, updatedName.Status, factcheck.StatusTopicResolved)

		// Verify name update in database via GetByID
		reqGetAfterNameUpdate, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/"+created.ID, nil)
		assertEq(t, err, nil)
		respGetAfterNameUpdate, err := http.DefaultClient.Do(reqGetAfterNameUpdate)
		assertEq(t, err, nil)
		defer respGetAfterNameUpdate.Body.Close()
		assertEq(t, respGetAfterNameUpdate.StatusCode, http.StatusOK)

		actualAfterNameUpdate := factcheck.Topic{}
		err = json.NewDecoder(respGetAfterNameUpdate.Body).Decode(&actualAfterNameUpdate)
		assertEq(t, err, nil)
		assertEq(t, actualAfterNameUpdate.Name, newName)
		assertEq(t, actualAfterNameUpdate.Description, desc)
		assertEq(t, actualAfterNameUpdate.Status, factcheck.StatusTopicResolved)

		// Test UpdateTopicDescription
		t.Log("Testing UpdateTopicDescription")
		newDesc := "Updated description for testing"
		updateDescBody := reqBodyJSON(struct {
			Description string `json:"description"`
		}{
			Description: newDesc,
		})
		reqUpdateDesc, err := http.NewRequestWithContext(t.Context(), http.MethodPut, testServer.URL+"/topics/"+created.ID+"/description", updateDescBody)
		assertEq(t, err, nil)
		respUpdateDesc, err := http.DefaultClient.Do(reqUpdateDesc)
		assertEq(t, err, nil)
		defer respUpdateDesc.Body.Close()
		assertEq(t, respUpdateDesc.StatusCode, http.StatusOK)

		// Assert UpdateTopicDescription response
		updatedDesc := factcheck.Topic{}
		err = json.NewDecoder(respUpdateDesc.Body).Decode(&updatedDesc)
		assertEq(t, err, nil)
		assertEq(t, updatedDesc.Description, newDesc)
		assertEq(t, updatedDesc.Name, newName)                         // Name should remain unchanged
		assertEq(t, updatedDesc.Status, factcheck.StatusTopicResolved) // Status should remain unchanged

		// Verify description update in database via GetByID
		reqGetAfterDescUpdate, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/"+created.ID, nil)
		assertEq(t, err, nil)
		respGetAfterDescUpdate, err := http.DefaultClient.Do(reqGetAfterDescUpdate)
		assertEq(t, err, nil)
		defer respGetAfterDescUpdate.Body.Close()
		assertEq(t, respGetAfterDescUpdate.StatusCode, http.StatusOK)

		actualAfterDescUpdate := factcheck.Topic{}
		err = json.NewDecoder(respGetAfterDescUpdate.Body).Decode(&actualAfterDescUpdate)
		assertEq(t, err, nil)
		assertEq(t, actualAfterDescUpdate.Description, newDesc)
		assertEq(t, actualAfterDescUpdate.Name, newName)
		assertEq(t, actualAfterDescUpdate.Status, factcheck.StatusTopicResolved)

		t.Log("Testing DeleteTopicByID")
		reqDelete, err := http.NewRequestWithContext(t.Context(), http.MethodDelete, testServer.URL+"/topics/"+created.ID, nil)
		assertEq(t, err, nil)
		respDelete, err := http.DefaultClient.Do(reqDelete)
		assertEq(t, err, nil)
		defer respDelete.Body.Close()
		assertEq(t, respDelete.StatusCode, http.StatusOK)

		// Get by ID should return 404
		reqGetByID, err = http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/"+created.ID, nil)
		assertEq(t, err, nil)
		respGetByID, err = http.DefaultClient.Do(reqGetByID)
		assertEq(t, err, nil)
		defer respGetByID.Body.Close()
		assertEq(t, respGetByID.StatusCode, http.StatusNotFound)
	})

	t.Run("Test not found scenarios", func(t *testing.T) {
		// Test GetByID with non-existent ID
		nonExistentID := "00000000-0000-0000-0000-000000000000"
		reqGetByID, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/"+nonExistentID, nil)
		assertEq(t, err, nil)
		respGetByID, err := http.DefaultClient.Do(reqGetByID)
		assertEq(t, err, nil)
		defer respGetByID.Body.Close()

		// Debug: read the response body to see what error we're getting
		respBody, err := io.ReadAll(respGetByID.Body)
		assertEq(t, err, nil)
		t.Logf("GetByID response status: %d, body: %s", respGetByID.StatusCode, string(respBody))

		assertEq(t, respGetByID.StatusCode, http.StatusNotFound)

		// Test UpdateTopicStatus with non-existent ID
		updateStatusBody := reqBodyJSON(struct {
			Status string `json:"status"`
		}{
			Status: string(factcheck.StatusTopicResolved),
		})
		reqUpdateStatus, err := http.NewRequestWithContext(t.Context(), http.MethodPut, testServer.URL+"/topics/"+nonExistentID+"/status", updateStatusBody)
		assertEq(t, err, nil)
		reqUpdateStatus.Header.Set("Content-Type", "application/json")
		respUpdateStatus, err := http.DefaultClient.Do(reqUpdateStatus)
		assertEq(t, err, nil)
		defer respUpdateStatus.Body.Close()
		assertEq(t, respUpdateStatus.StatusCode, http.StatusNotFound)

		// Test UpdateTopicName with non-existent ID
		updateNameBody := reqBodyJSON(struct {
			Name string `json:"name"`
		}{
			Name: "Updated name",
		})
		reqUpdateName, err := http.NewRequestWithContext(t.Context(), http.MethodPut, testServer.URL+"/topics/"+nonExistentID+"/name", updateNameBody)
		assertEq(t, err, nil)
		reqUpdateName.Header.Set("Content-Type", "application/json")
		respUpdateName, err := http.DefaultClient.Do(reqUpdateName)
		assertEq(t, err, nil)
		defer respUpdateName.Body.Close()
		assertEq(t, respUpdateName.StatusCode, http.StatusNotFound)

		// Test UpdateTopicDescription with non-existent ID
		updateDescBody := reqBodyJSON(struct {
			Description string `json:"description"`
		}{
			Description: "Updated description",
		})
		reqUpdateDesc, err := http.NewRequestWithContext(t.Context(), http.MethodPut, testServer.URL+"/topics/"+nonExistentID+"/description", updateDescBody)
		assertEq(t, err, nil)
		reqUpdateDesc.Header.Set("Content-Type", "application/json")
		respUpdateDesc, err := http.DefaultClient.Do(reqUpdateDesc)
		assertEq(t, err, nil)
		defer respUpdateDesc.Body.Close()
		assertEq(t, respUpdateDesc.StatusCode, http.StatusNotFound)

		// Test that the error response contains detailed information
		body, err := io.ReadAll(respUpdateDesc.Body)
		assertEq(t, err, nil)
		errorMessage := string(body)
		t.Logf("Error message: %s", errorMessage)
		// The error message should contain the filter information from our custom error type
		assertEq(t, strings.Contains(errorMessage, "not found for filter"), true)
		assertEq(t, strings.Contains(errorMessage, nonExistentID), true)
	})
}

func TestHandlerTopic_CountTopicsHome(t *testing.T) {
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

	topic4 := factcheck.Topic{
		ID:          "660e8400-e29b-41d4-a716-446655440004",
		Name:        "Topic 4 - Sports Resolved",
		Description: "Sports news and updates (resolved)",
		Status:      factcheck.StatusTopicResolved,
		Result:      "Verified as false",
		CreatedAt:   now,
		UpdatedAt:   nil,
	}

	// Create topics in database
	createdTopic1, err := app.Repository.Topics.Create(t.Context(), topic1)
	if err != nil {
		t.Fatalf("Failed to create topic1: %v", err)
	}

	createdTopic2, err := app.Repository.Topics.Create(t.Context(), topic2)
	if err != nil {
		t.Fatalf("Failed to create topic2: %v", err)
	}

	createdTopic3, err := app.Repository.Topics.Create(t.Context(), topic3)
	if err != nil {
		t.Fatalf("Failed to create topic3: %v", err)
	}

	createdTopic4, err := app.Repository.Topics.Create(t.Context(), topic4)
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
	createdUserMessage1, err := app.Repository.UserMessages.Create(t.Context(), userMessage1)
	if err != nil {
		t.Fatalf("Failed to create userMessage1: %v", err)
	}

	createdUserMessage2, err := app.Repository.UserMessages.Create(t.Context(), userMessage2)
	if err != nil {
		t.Fatalf("Failed to create userMessage2: %v", err)
	}

	createdUserMessage3, err := app.Repository.UserMessages.Create(t.Context(), userMessage3)
	if err != nil {
		t.Fatalf("Failed to create userMessage3: %v", err)
	}

	createdUserMessage4, err := app.Repository.UserMessages.Create(t.Context(), userMessage4)
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
	_, err = app.Repository.Messages.Create(t.Context(), message1)
	if err != nil {
		t.Fatalf("Failed to create message1: %v", err)
	}

	_, err = app.Repository.Messages.Create(t.Context(), message2)
	if err != nil {
		t.Fatalf("Failed to create message2: %v", err)
	}

	_, err = app.Repository.Messages.Create(t.Context(), message3)
	if err != nil {
		t.Fatalf("Failed to create message3: %v", err)
	}

	_, err = app.Repository.Messages.Create(t.Context(), message4)
	if err != nil {
		t.Fatalf("Failed to create message4: %v", err)
	}

	t.Run("CountTopicsHome - no query parameters (all topics)", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for each status plus total
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(2))
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(2))
		assertEq(t, result["total"], int64(4))
	})

	t.Run("CountTopicsHome - like_id filter only", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_id=550e8400", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for topics with '550e8400' in ID
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(1))  // topic1
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(1)) // topic2
		assertEq(t, result["total"], int64(2))
	})

	t.Run("CountTopicsHome - like_message_text filter only", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_message_text=COVID", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for topics with COVID messages
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(1))  // topic1
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(0)) // none
		assertEq(t, result["total"], int64(1))
	})

	t.Run("CountTopicsHome - both filters", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_id=550e8400&like_message_text=COVID", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for topics with both filters
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(1))  // topic1
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(0)) // none
		assertEq(t, result["total"], int64(1))
	})

	t.Run("CountTopicsHome - case insensitive message text filter", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_message_text=covid", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for topics with COVID messages (case insensitive)
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(1))  // topic1
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(0)) // none
		assertEq(t, result["total"], int64(1))
	})

	t.Run("CountTopicsHome - no matches for filters", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_id=99999999&like_message_text=nonexistent", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have zero counts
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(0))
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(0))
		assertEq(t, result["total"], int64(0))
	})

	t.Run("CountTopicsHome - empty string filters", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_id=&like_message_text=", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for all topics (empty filters ignored)
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(2))
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(2))
		assertEq(t, result["total"], int64(4))
	})

	t.Run("CountTopicsHome - partial message text match", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_message_text=technology", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for topics with technology messages
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(1))  // topic3
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(0)) // none
		assertEq(t, result["total"], int64(1))
	})

	t.Run("CountTopicsHome - ID pattern with wildcards", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_id=550e8400%25", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for topics with '550e8400' prefix in ID
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(1))  // topic1
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(1)) // topic2
		assertEq(t, result["total"], int64(2))
	})

	t.Run("CountTopicsHome - multiple status counts", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count?like_message_text=results", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Should have counts for topics with "results" in messages (topic2 and topic4)
		assertEq(t, result[string(factcheck.StatusTopicPending)], int64(0))  // none
		assertEq(t, result[string(factcheck.StatusTopicResolved)], int64(2)) // topic2, topic4
		assertEq(t, result["total"], int64(2))
	})

	t.Run("CountTopicsHome - response structure validation", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/count", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var result map[string]int64
		err = json.NewDecoder(resp.Body).Decode(&result)
		assertEq(t, err, nil)

		// Verify response structure
		_, hasPending := result[string(factcheck.StatusTopicPending)]
		_, hasResolved := result[string(factcheck.StatusTopicResolved)]
		_, hasTotal := result["total"]

		if !hasPending {
			t.Errorf("Expected '%s' key in response", string(factcheck.StatusTopicPending))
		}
		if !hasResolved {
			t.Errorf("Expected '%s' key in response", string(factcheck.StatusTopicResolved))
		}
		if !hasTotal {
			t.Errorf("Expected 'total' key in response")
		}

		// Verify total is sum of individual counts
		expectedTotal := result[string(factcheck.StatusTopicPending)] + result[string(factcheck.StatusTopicResolved)]
		if result["total"] != expectedTotal {
			t.Errorf("Expected total to be sum of pending and resolved, got total=%d, sum=%d", result["total"], expectedTotal)
		}
	})
}

func TestHandlerTopic_ListTopicsHome(t *testing.T) {
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
	createdTopic1, err := app.Repository.Topics.Create(t.Context(), topic1)
	if err != nil {
		t.Fatalf("Failed to create topic1: %v", err)
	}

	createdTopic2, err := app.Repository.Topics.Create(t.Context(), topic2)
	if err != nil {
		t.Fatalf("Failed to create topic2: %v", err)
	}

	createdTopic3, err := app.Repository.Topics.Create(t.Context(), topic3)
	if err != nil {
		t.Fatalf("Failed to create topic3: %v", err)
	}

	// Create message groups with different languages and texts
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
	_, err = app.Repository.MessageGroups.Create(t.Context(), messageGroup1)
	if err != nil {
		t.Fatalf("Failed to create messageGroup1: %v", err)
	}

	_, err = app.Repository.MessageGroups.Create(t.Context(), messageGroup2)
	if err != nil {
		t.Fatalf("Failed to create messageGroup2: %v", err)
	}

	_, err = app.Repository.MessageGroups.Create(t.Context(), messageGroup3)
	if err != nil {
		t.Fatalf("Failed to create messageGroup3: %v", err)
	}

	t.Run("ListTopicsV2 - no query parameters (all topics)", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 3)
	})

	t.Run("ListTopicsV2 - like_id filter only", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_id=550e8400", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 2)

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

	t.Run("ListTopicsV2 - like_message_text filter only (English)", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_message_text=COVID", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 1)

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (COVID topic), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListTopicsV2 - like_message_text filter only (Thai)", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_message_text=ข่าวปลอม", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 1)

		if topics[0].ID != createdTopic2.ID {
			t.Errorf("Expected topic2 (Thai election topic), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListTopicsV2 - both filters", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_id=550e8400&like_message_text=COVID", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 1)

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (matches both filters), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListTopicsV2 - case insensitive message text filter", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_message_text=covid", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 1)

		if topics[0].ID != createdTopic1.ID {
			t.Errorf("Expected topic1 (case insensitive COVID match), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListTopicsV2 - no matches for filters", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_id=99999999&like_message_text=nonexistent", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 0)
	})

	t.Run("ListTopicsV2 - empty string filters", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_id=&like_message_text=", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 3)
	})

	t.Run("ListTopicsV2 - partial message text match", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_message_text=technology", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 1)

		if topics[0].ID != createdTopic3.ID {
			t.Errorf("Expected topic3 (technology topic), got topic with ID %s", topics[0].ID)
		}
	})

	t.Run("ListTopicsV2 - ID pattern with wildcards", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?like_id=550e8400%25", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 2)

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

	t.Run("ListTopicsV2 - pagination", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?limit=2&offset=0", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 2)
	})

	t.Run("ListTopicsV2 - pagination with offset", func(t *testing.T) {
		req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/?limit=1&offset=1", nil)
		assertEq(t, err, nil)
		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusOK)

		var topics []factcheck.Topic
		err = json.NewDecoder(resp.Body).Decode(&topics)
		assertEq(t, err, nil)
		assertEq(t, len(topics), 1)
	})
}

func assertEq[X comparable](t *testing.T, actual, expected X) {
	if actual != expected {
		t.Logf("actual: %+v", actual)
		t.Logf("expected: %+v", expected)
		t.Fatalf("assertEq: unexpected value for type %T", actual)
	}
}

// nolint:unused
func assertNeq[X comparable](t *testing.T, actual, notExpected X) {
	if actual == notExpected {
		t.Logf("not expected: %+v", notExpected)
		t.Fatalf("assertEq: unexpected value for type %T", actual)
	}
}

func reqBodyJSON(data any) *bytes.Buffer {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return buf
}
