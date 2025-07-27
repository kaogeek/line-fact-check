package factcheck_test

import (
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
)

func TestValidate(t *testing.T) {
	shouldOk := []interface{ IsValid() bool }{
		factcheck.StatusTopicPending,
		factcheck.StatusTopicResolved,
		factcheck.TypeMessageText,
	}
	for i := range shouldOk {
		s := shouldOk[i]
		if s.IsValid() {
			continue
		}
		t.Fatalf("unexpected invalid value: %v", s)
	}
}
