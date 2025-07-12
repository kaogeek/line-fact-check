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
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

type TestSuite struct {
	container di.Container
}

func init() {
	slog.Info("handlers_test.level.slog")
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("handlers_test.level.slog=DEBUG")
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
