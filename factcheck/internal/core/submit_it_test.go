//go:build integration_test
// +build integration_test

package core_test

import (
	"testing"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di/ittest"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

func TestSubmit(t *testing.T) {
	msgText := "TestSubmit.Text"
	user := factcheck.UserInfo{
		UserType: factcheck.TypeUserMessageLINEChat,
		UserID:   "TestSubmit.UserID",
	}

	t.Run("error - empty text", func(t *testing.T) {
		container, cleanup, err := ittest.InitializeContainerTest(t)
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
		container, cleanup, err := ittest.InitializeContainerTest(t)
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
		container, cleanup, err := ittest.InitializeContainerTest(t)
		if err != nil {
			panic(err)
		}
		defer cleanup()
		ctx := t.Context()
		service := container.Service
		repo := container.Repository

		text0 := "TestSubmit.1.text0"
		text1 := "TestSubmit.1.text1"

		msg0, group, topic, err := service.Submit(ctx, user, text0, "")
		if err != nil {
			t.Fatal(err)
		}
		if msg0.GroupID != group.ID {
			t.Fatalf("unexpected message group ID '%s', expecting '%s'", msg0.GroupID, group.ID)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}

		// Save group 1 values
		groupID := group.ID
		groupSHA1 := group.TextSHA1

		// msg1 has the same text0 as with msg0,
		// so they should fall under the same message group
		msg1, group, topic, err := service.Submit(ctx, user, text0, "")
		if err != nil {
			t.Fatal(err)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}
		if msg1.GroupID != group.ID {
			t.Fatalf("unexpected message group ID '%s', expecting '%s'", msg0.GroupID, groupID)
		}
		if group.ID != groupID {
			t.Fatal("unexpected group ID", groupID, group.ID)
		}
		if group.TextSHA1 != groupSHA1 {
			t.Fatal("unexpected group SHA1", groupSHA1, group.TextSHA1)
		}

		msg2, group, topic, err := service.Submit(ctx, user, text1, "")
		if err != nil {
			t.Fatal(err)
		}
		if topic != nil {
			t.Fatal("unexpected topic", topic)
		}
		if msg2.GroupID != group.ID {
			t.Fatalf("unexpected message group ID '%s', expecting '%s'", msg0.GroupID, groupID)
		}
		if group.ID == groupID {
			t.Fatal("expected group IDs to be different", groupID)
		}

		list, err := repo.MessagesV2.ListByGroup(ctx, groupID)
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 2 {
			t.Fatal("unexpected length", len(list))
		}
		if list[0].ID != msg0.ID {
			t.Fatal("unexpected msgs[0].ID")
		}
		if list[1].ID != msg1.ID {
			t.Fatal("unexpected msgs[1].ID")
		}
	})
}
