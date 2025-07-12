//go:build integration_test
// +build integration_test

package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

type TestSuite struct {
	container di.Container
}

func init() {
	// Set slog level to DEBUG and log with json formatter
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
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

func TestHandlerTopic_Stateful(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Create test server
	testServer := httptest.NewServer(app.Server.(*http.Server).Handler)
	defer testServer.Close()

	// Clear all data
	t.Log("WARN: Clearing all data from database")
	app.PostgresConn.Exec(t.Context(), "DELETE FROM topics")
	app.PostgresConn.Exec(t.Context(), "DELETE FROM messages")
	app.PostgresConn.Exec(t.Context(), "DELETE FROM user_messages")
	t.Log("WARN: Cleared all data from database")

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
			ID:           created.ID,
			Name:         name,
			Description:  desc,
			Status:       factcheck.StatusTopicPending,
			Result:       "",
			ResultStatus: factcheck.StatusTopicResultNone,
			CreatedAt:    now,
			UpdatedAt:    nil,
		}
		assertEq(t, created, expected)

		// Assert in database
		actualDB, err := app.Repository.Topic.GetByID(t.Context(), created.ID)
		assertEq(t, err, nil)
		assertEq(t, actualDB, expected)

		t.Log("Testing ListTopics")
		reqList, err := http.NewRequestWithContext(t.Context(), http.MethodGet, testServer.URL+"/topics/", nil)
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
			ID:           created.ID,
			Name:         name,
			Description:  desc,
			Status:       factcheck.StatusTopicResolved,
			ResultStatus: factcheck.StatusTopicResultNone,
			CreatedAt:    now,
			//nolint:unused
			UpdatedAt: nil, // Underlying database will set this to NOW()
		}
		assertEq(t, updatedStatus.ID, expectedUpdateStatus.ID)
		assertEq(t, updatedStatus.Name, expectedUpdateStatus.Name)
		assertEq(t, updatedStatus.Description, expectedUpdateStatus.Description)
		assertEq(t, updatedStatus.Status, expectedUpdateStatus.Status)
		assertEq(t, updatedStatus.Result, expectedUpdateStatus.Result)
		assertEq(t, updatedStatus.ResultStatus, expectedUpdateStatus.ResultStatus)
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
