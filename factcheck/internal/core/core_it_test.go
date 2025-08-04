//go:build integration_test
// +build integration_test

package core_test

import (
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
)

func TestSubmit(t *testing.T) {
	t.Run("normal - all new", func(t *testing.T) {
		container, cleanup, err := di.InitializeContainerTest()
		if err != nil {
			panic(err)
		}
		defer cleanup()
		ctx := t.Context()
		service := container.Service
		repo := container.Repository

		msgText := "text-it-test"
		user := factcheck.UserInfo{
			UserType: factcheck.TypeUserMessageLINEChat,
			UserID:   "it-test",
		}

		msg, group, topic, err := service.Submit(ctx, user, msgText, "")
		if err != nil {
			t.Fatal(err)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}
		groupID := group.ID
		groupSHA1 := group.TextSHA1

		msg2, group, topic, err := service.Submit(ctx, user, msgText, "")
		if err != nil {
			t.Fatal(err)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}
		if group.ID != groupID {
			t.Fatal("unexpected group ID", groupID, group.ID)
		}
		if group.TextSHA1 != groupSHA1 {
			t.Fatal("unexpected group SHA1", groupSHA1, group.TextSHA1)
		}

		msgs, err := repo.MessagesV2.ListByGroup(ctx, groupID)
		if err != nil {
			t.Fatal(err)
		}
		if len(msgs) != 2 {
			t.Fatal("unexpected length", len(msgs))
		}
		if msgs[0].ID != msg.ID {
			t.Fatal("unexpected msgs[0].ID")
		}
		if msgs[1].ID != msg2.ID {
			t.Fatal("unexpected msgs[1].ID")
		}
	})
}
