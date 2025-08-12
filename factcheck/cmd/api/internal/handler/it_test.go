package handler_test

import (
	"bytes"
	"encoding/json"
	"testing"
)

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
