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
	"time"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
)

type TestSuite struct {
	container di.Container
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
	const timeout = time.Millisecond * 200
	const url = "localhost:8778/"

	app, cleanup, err := di.InitializeContainer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Clear all data
	t.Log("WARN: Clearing all data from database")
	app.PostgresConn.Exec(t.Context(), "DELETE FROM topics")
	app.PostgresConn.Exec(t.Context(), "DELETE FROM messages")
	app.PostgresConn.Exec(t.Context(), "DELETE FROM user_messages")
	t.Log("WARN: Cleared all data from database")

	t.Run("CRUD", func(t *testing.T) {
		now := time.Now()
		name := fmt.Sprintf("topic-test-normal-%s", now.String())
		desc := fmt.Sprintf("topic-test-normal-%s-desc", now.String())
		topic := factcheck.Topic{
			Name:        name,
			Description: desc,
		}

		body := reqBodyJSON(topic)
		reqCreate, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, body)
		assertEq(t, err, nil)
		resp := httptest.NewRecorder()
		app.Handler.CreateTopic(resp, reqCreate)

		// Assert response
		created := factcheck.Topic{}
		err = json.Unmarshal(resp.Body.Bytes(), &created)
		assertEq(t, err, nil)
		assertNeq(t, created.ID, "")
		assertEq(t, created.Name, name)
		assertEq(t, created.Description, desc)
		assertEq(t, created.Status, factcheck.StatusTopicPending)

		// Assert in database
		actual, err := app.Repository.Topic.GetByID(t.Context(), created.ID)
		assertEq(t, err, nil)
		assertEq(t, actual, created)

		reqList, err := http.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
		assertEq(t, err, nil)
		respList := httptest.NewRecorder()
		app.Handler.ListTopics(respList, reqList)

		// Assert response
		actualList := []factcheck.Topic{}
		err = json.Unmarshal(respList.Body.Bytes(), &actualList)
		assertEq(t, err, nil)
		assertEq(t, len(actualList), 1)
		assertEq(t, actualList[0], created)
	})
}

func assertEq[X comparable](t *testing.T, actual, expected X) {
	if actual != expected {
		t.Logf("actual: %+v", actual)
		t.Logf("expected: %+v", expected)
		t.Fatalf("assertEq: unexpected value for type %T", actual)
	}
}

func assertNeq[X comparable](t *testing.T, actual, notExpected X) {
	if actual == notExpected {
		t.Logf("not expected: %+v", notExpected)
		t.Fatalf("assertEq: unexpected value for type %T", actual)
	}
}
