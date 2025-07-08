//go:build integration_test
// +build integration_test

package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func reqBodyJSON(data any) io.Reader {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		panic(err)
	}
	return buf
}

func TestHandlerTopic(t *testing.T) {
	const timeout = time.Millisecond * 200
	const url = "localhost:8778/"

	t.Run("create", func(t *testing.T) {
		app, cleanup, err := di.InitializeContainer()
		if err != nil {
			panic(err)
		}
		defer cleanup()

		name := fmt.Sprintf("topic-test-normal-%s", time.Now().String())
		topic := factcheck.Topic{
			Name: name,
		}

		req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, url, reqBodyJSON(topic))
		assertEq(t, err, nil)
		resp := httptest.NewRecorder()
		app.Handler.CreateTopic(resp, req)

		// Assert response
		actual1 := factcheck.Topic{}
		err = json.Unmarshal(resp.Body.Bytes(), &actual1)
		assertEq(t, err, nil)
		assertNeq(t, actual1.ID, "")
		assertEq(t, actual1.Name, name)
		assertEq(t, actual1.Status, factcheck.StatusTopicPending)

		// Assert in database
		actual2, err := app.Repository.Topic.GetByID(t.Context(), actual1.ID)
		assertEq(t, err, nil)
		assertEq(t, actual2, actual1)
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
