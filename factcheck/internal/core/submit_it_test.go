//go:build integration_test
// +build integration_test

package core_test

import (
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestSubmit(t *testing.T) {
	msgText := "TestSubmit.Text"
	user := factcheck.UserInfo{
		UserType: factcheck.TypeUserMessageLINEChat,
		UserID:   "TestSubmit.UserID",
	}

	t.Run("error - empty text", func(t *testing.T) {
		container, cleanup, err := di.InitializeContainerTest()
		if err != nil {
			panic(err)
		}
		defer cleanup()

		emptyText := ""
		_, _, _, err = container.Service.Submit(t.Context(), user, emptyText, "")
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("error - no such topic", func(t *testing.T) {
		container, cleanup, err := di.InitializeContainerTest()
		if err != nil {
			panic(err)
		}
		defer cleanup()

		noSuchTopicID := utils.NewID().String()
		_, _, _, err = container.Service.Submit(t.Context(), user, msgText, noSuchTopicID)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("normal - all new", func(t *testing.T) {
		container, cleanup, err := di.InitializeContainerTest()
		if err != nil {
			panic(err)
		}
		defer cleanup()
		ctx := t.Context()
		service := container.Service
		repo := container.Repository

		msg1, group, topic, err := service.Submit(ctx, user, msgText, "")
		if err != nil {
			t.Fatal(err)
		}
		if msg1.GroupID != group.ID {
			t.Fatalf("unexpected message group ID '%s', expecting '%s'", msg1.GroupID, group.ID)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}

		// Save group 1 values
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
		if msg2.GroupID != group.ID {
			t.Fatalf("unexpected message group ID '%s', expecting '%s'", msg1.GroupID, groupID)
		}
		if group.ID != groupID {
			t.Fatal("unexpected group ID", groupID, group.ID)
		}
		if group.TextSHA1 != groupSHA1 {
			t.Fatal("unexpected group SHA1", groupSHA1, group.TextSHA1)
		}

		// Different text, different hash, different group
		msgText2 := "TestSubmit.1.msgText2"
		msg3, group, topic, err := service.Submit(ctx, user, msgText2, "")
		if err != nil {
			t.Fatal(err)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}
		if msg3.GroupID != group.ID {
			t.Fatalf("unexpected message group ID '%s', expecting '%s'", msg1.GroupID, groupID)
		}
		if group.ID == groupID {
			t.Fatal("expected group IDs to be different", groupID)
		}

		msgs, err := repo.MessagesV2.ListByGroup(ctx, groupID)
		if err != nil {
			t.Fatal(err)
		}
		if len(msgs) != 2 {
			t.Fatal("unexpected length", len(msgs))
		}
		if msgs[0].ID != msg1.ID {
			t.Fatal("unexpected msgs[0].ID")
		}
		if msgs[1].ID != msg2.ID {
			t.Fatal("unexpected msgs[1].ID")
		}
	})
}
