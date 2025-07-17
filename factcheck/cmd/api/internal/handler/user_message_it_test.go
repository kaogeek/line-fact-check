//go:build integration_test
// +build integration_test

package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

type newUserMessageResponse struct {
	UserMessage factcheck.UserMessage `json:"user_message"`
	Message     factcheck.Message     `json:"message"`
}

func TestHandlerUserMessage_NewUserMessage(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Create test server
	testServer := httptest.NewServer(app.Server.(*http.Server).Handler)
	defer testServer.Close()

	t.Run("Create user message without topic", func(t *testing.T) {
		now := utils.TimeNow().Round(0) // Postgres timestampz will not preserve monotonic clock
		utils.TimeFreeze(now)
		defer utils.TimeUnfreeze()

		text := "This is a test message without topic"

		t.Log("Testing NewUserMessage without topic")
		body := reqBodyJSON(struct {
			Text    string `json:"text"`
			TopicID string `json:"topic_id"`
		}{
			Text:    text,
			TopicID: "",
		})

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, testServer.URL+"/user-messages/", body)
		assertEq(t, err, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusCreated)

		// Assert response structure
		var response newUserMessageResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assertEq(t, err, nil)

		// Verify user message fields
		assertNeq(t, response.UserMessage.ID, "")
		assertEq(t, response.UserMessage.Type, factcheck.TypeUserMessageAdmin)
		assertEq(t, response.UserMessage.RepliedAt, nil)
		assertNeq(t, len(response.UserMessage.Metadata), 0)
		assertEq(t, response.UserMessage.CreatedAt, now)
		assertEq(t, response.UserMessage.UpdatedAt, nil)

		// Verify message fields
		assertNeq(t, response.Message.ID, "")
		assertEq(t, response.Message.UserMessageID, response.UserMessage.ID)
		assertEq(t, response.Message.Type, factcheck.TypeMessageText)
		assertEq(t, response.Message.Status, factcheck.StatusMessageSubmitted)
		assertEq(t, response.Message.TopicID, "")
		assertEq(t, response.Message.Text, text)
		assertEq(t, response.Message.CreatedAt, now)
		assertEq(t, response.Message.UpdatedAt, nil)

		// Verify in database
		actualUserMessage, err := app.Repository.UserMessages.GetByID(t.Context(), response.UserMessage.ID)
		assertEq(t, err, nil)
		assertEq(t, actualUserMessage.ID, response.UserMessage.ID)
		assertEq(t, actualUserMessage.Type, factcheck.TypeUserMessageAdmin)

		actualMessage, err := app.Repository.Messages.GetByID(t.Context(), response.Message.ID)
		assertEq(t, err, nil)
		assertEq(t, actualMessage.ID, response.Message.ID)
		assertEq(t, actualMessage.UserMessageID, response.UserMessage.ID)
		assertEq(t, actualMessage.Text, text)
		assertEq(t, actualMessage.Status, factcheck.StatusMessageSubmitted)
		assertEq(t, actualMessage.TopicID, "")
	})

	t.Run("Create user message with topic", func(t *testing.T) {
		now := utils.TimeNow().Round(0)
		utils.TimeFreeze(now)
		defer utils.TimeUnfreeze()

		// First create a topic
		topicName := fmt.Sprintf("topic-test-user-message-%s", now.String())
		topicDesc := fmt.Sprintf("topic-test-user-message-desc-%s", now.String())

		topicBody := reqBodyJSON(factcheck.Topic{
			Name:        topicName,
			Description: topicDesc,
		})

		reqTopic, err := http.NewRequestWithContext(t.Context(), http.MethodPost, testServer.URL+"/topics/", topicBody)
		assertEq(t, err, nil)
		reqTopic.Header.Set("Content-Type", "application/json")

		respTopic, err := http.DefaultClient.Do(reqTopic)
		assertEq(t, err, nil)
		defer respTopic.Body.Close()
		assertEq(t, respTopic.StatusCode, http.StatusCreated)

		var createdTopic factcheck.Topic
		err = json.NewDecoder(respTopic.Body).Decode(&createdTopic)
		assertEq(t, err, nil)

		// Now create user message with topic
		text := "This is a test message with topic"

		t.Log("Testing NewUserMessage with topic")
		body := reqBodyJSON(struct {
			Text    string `json:"text"`
			TopicID string `json:"topic_id"`
		}{
			Text:    text,
			TopicID: createdTopic.ID,
		})

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, testServer.URL+"/user-messages/", body)
		assertEq(t, err, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusCreated)

		// Assert response structure
		var response newUserMessageResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assertEq(t, err, nil)

		// Verify message has correct status and topic_id
		assertEq(t, response.Message.Status, factcheck.StatusMessageTopicSubmitted)
		assertEq(t, response.Message.TopicID, createdTopic.ID)
		assertEq(t, response.Message.Text, text)

		// Verify in database
		actualMessage, err := app.Repository.Messages.GetByID(t.Context(), response.Message.ID)
		assertEq(t, err, nil)
		assertEq(t, actualMessage.Status, factcheck.StatusMessageTopicSubmitted)
		assertEq(t, actualMessage.TopicID, createdTopic.ID)
	})

	t.Run("Create user message with invalid JSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{"text": "test", "topic_id": invalid}`)

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, testServer.URL+"/user-messages/", body)
		assertEq(t, err, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusBadRequest)
	})

	t.Run("Create user message with missing text", func(t *testing.T) {
		body := reqBodyJSON(struct {
			TopicID string `json:"topic_id"`
		}{
			TopicID: "",
		})

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, testServer.URL+"/user-messages/", body)
		assertEq(t, err, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		assertEq(t, err, nil)
		defer resp.Body.Close()
		assertEq(t, resp.StatusCode, http.StatusCreated) // Handler accepts empty text
	})
}
