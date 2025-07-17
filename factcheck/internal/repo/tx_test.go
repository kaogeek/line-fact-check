//go:build integration_test
// +build integration_test

package repo_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

// TestTransactionIsolationLevels tests different isolation levels with race conditions
func TestTransactionIsolationLevels(t *testing.T) {
	t.Run("ReadCommitted", func(t *testing.T) {
		t.Run("CommitAfterTx2", func(t *testing.T) {
			app, cleanup, err := di.InitializeContainerTest()
			if err != nil {
				t.Fatalf("Failed to initialize test container: %v", err)
			}
			defer cleanup()

			ctx := context.Background()

			// Create test data with unique ID
			now := utils.TimeNow().Round(0)
			utils.TimeFreeze(now)
			defer utils.TimeUnfreeze()

			sharedTopic := factcheck.Topic{
				ID:           "550e8400-e29b-41d4-a716-446655440001",
				Name:         "Shared Topic - ReadCommitted",
				Description:  "This topic will be updated by competing transactions",
				Status:       factcheck.StatusTopicPending,
				Result:       "",
				ResultStatus: factcheck.StatusTopicResultNone,
				CreatedAt:    now,
				UpdatedAt:    nil,
			}

			createdTopic, err := app.Repository.Topics.Create(ctx, sharedTopic)
			if err != nil {
				t.Fatalf("Failed to create shared topic: %v", err)
			}
			testReadCommittedIsolation_CommitAfterTx2(t, &app.Repository, createdTopic.ID)
		})

		t.Run("CommitBeforeTx2", func(t *testing.T) {
			app, cleanup, err := di.InitializeContainerTest()
			if err != nil {
				t.Fatalf("Failed to initialize test container: %v", err)
			}
			defer cleanup()
			ctx := context.Background()

			// Create test data with unique ID
			now := utils.TimeNow().Round(0)
			utils.TimeFreeze(now)
			defer utils.TimeUnfreeze()

			sharedTopic := factcheck.Topic{
				ID:           "550e8400-e29b-41d4-a716-446655440001",
				Name:         "Shared Topic - ReadCommitted",
				Description:  "This topic will be updated by competing transactions",
				Status:       factcheck.StatusTopicPending,
				Result:       "",
				ResultStatus: factcheck.StatusTopicResultNone,
				CreatedAt:    now,
				UpdatedAt:    nil,
			}

			createdTopic, err := app.Repository.Topics.Create(ctx, sharedTopic)
			if err != nil {
				t.Fatalf("Failed to create shared topic: %v", err)
			}
			testReadCommittedIsolation_CommitBeforeTx2(t, &app.Repository, createdTopic.ID)
		})
	})

	t.Run("RepeatableRead_ShouldPreventDirtyReads", func(t *testing.T) {
		app, cleanup, err := di.InitializeContainerTest()
		if err != nil {
			t.Fatalf("Failed to initialize test container: %v", err)
		}
		defer cleanup()
		ctx := context.Background()

		// Create test data with unique ID
		now := utils.TimeNow().Round(0)
		utils.TimeFreeze(now)
		defer utils.TimeUnfreeze()

		sharedTopic := factcheck.Topic{
			ID:           "550e8400-e29b-41d4-a716-446655440002",
			Name:         "Shared Topic - RepeatableRead",
			Description:  "This topic will be updated by competing transactions",
			Status:       factcheck.StatusTopicPending,
			Result:       "",
			ResultStatus: factcheck.StatusTopicResultNone,
			CreatedAt:    now,
			UpdatedAt:    nil,
		}

		createdTopic, err := app.Repository.Topics.Create(ctx, sharedTopic)
		if err != nil {
			t.Fatalf("Failed to create shared topic: %v", err)
		}
		testRepeatableReadIsolation(t, &app.Repository, createdTopic.ID)
	})

	t.Run("Serializable_ShouldPreventPhantomReads", func(t *testing.T) {
		app, cleanup, err := di.InitializeContainerTest()
		if err != nil {
			t.Fatalf("Failed to initialize test container: %v", err)
		}
		defer cleanup()

		ctx := context.Background()

		// Create test data with unique ID
		now := utils.TimeNow().Round(0)
		utils.TimeFreeze(now)
		defer utils.TimeUnfreeze()

		sharedTopic := factcheck.Topic{
			ID:           "550e8400-e29b-41d4-a716-446655440003",
			Name:         "Shared Topic - Serializable",
			Description:  "This topic will be updated by competing transactions",
			Status:       factcheck.StatusTopicPending,
			Result:       "",
			ResultStatus: factcheck.StatusTopicResultNone,
			CreatedAt:    now,
			UpdatedAt:    nil,
		}

		createdTopic, err := app.Repository.Topics.Create(ctx, sharedTopic)
		if err != nil {
			t.Fatalf("Failed to create shared topic: %v", err)
		}

		testSerializableIsolation(t, &app.Repository, createdTopic.ID)
	})

	t.Run("ConcurrentUpdates_ShouldHandleConflicts", func(t *testing.T) {
		app, cleanup, err := di.InitializeContainerTest()
		if err != nil {
			t.Fatalf("Failed to initialize test container: %v", err)
		}
		defer cleanup()

		ctx := context.Background()

		// Create test data with unique ID
		now := utils.TimeNow().Round(0)
		utils.TimeFreeze(now)
		defer utils.TimeUnfreeze()

		sharedTopic := factcheck.Topic{
			ID:           "550e8400-e29b-41d4-a716-446655440004",
			Name:         "Shared Topic - Concurrent",
			Description:  "This topic will be updated by competing transactions",
			Status:       factcheck.StatusTopicPending,
			Result:       "",
			ResultStatus: factcheck.StatusTopicResultNone,
			CreatedAt:    now,
			UpdatedAt:    nil,
		}

		createdTopic, err := app.Repository.Topics.Create(ctx, sharedTopic)
		if err != nil {
			t.Fatalf("Failed to create shared topic: %v", err)
		}

		testConcurrentUpdates(t, &app.Repository, createdTopic.ID)
	})
}

// testReadCommittedIsolation_CommitAfterTx2 tests that ReadCommitted allows dirty reads
// https://www.postgresql.org/docs/16/transaction-iso.html#XACT-READ-COMMITTED
func testReadCommittedIsolation_CommitAfterTx2(t *testing.T, r *repo.Repository, topicID string) {
	ctx := context.Background()
	var wg sync.WaitGroup
	var err1, err2 error
	ch := make(chan struct{})
	newDescription := "Updated by TX1"
	var topic factcheck.Topic

	// Start transaction 1 (updater)
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx1, err := r.BeginTx(ctx, repo.ReadCommitted)
		if err != nil {
			err1 = err
			return
		}
		close(ch)
		defer tx1.Rollback(ctx)
		// Update the topic
		_, err = r.Topics.UpdateDescription(ctx, topicID, newDescription, repo.WithTx(tx1))
		if err != nil {
			err1 = err
			return
		}
		err = tx1.Commit(t.Context())
		if err != nil {
			err1 = err
			return
		}
	}()

	// Start transaction 2 (reader)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ch
		var err error
		tx2, err := r.BeginTx(ctx, repo.ReadCommitted)
		if err != nil {
			err2 = err
			return
		}
		defer tx2.Rollback(ctx)
		// Try to read the topic - should see the uncommitted change
		topic, err = r.Topics.GetByID(ctx, topicID, repo.WithTx(tx2))
		if err != nil {
			err2 = err
			return
		}
		// In ReadCommitted, we should see the uncommitted change
		if topic.Description != newDescription {
			t.Errorf("Expected to see uncommitted change in ReadCommitted, got: %s", topic.Description)
		}
	}()

	wg.Wait()
	if err1 != nil {
		t.Errorf("Transaction 1 failed: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Transaction 2 failed: %v", err2)
	}
}

func testReadCommittedIsolation_CommitBeforeTx2(t *testing.T, r *repo.Repository, topicID string) {
	ctx := context.Background()
	var wg sync.WaitGroup
	var tx1, tx2 repo.Tx
	var err1, err2 error
	ch := make(chan struct{})
	newDescription := "Updated by TX1"

	// Start transaction 1 (updater)
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx1, err1 = r.BeginTx(ctx, repo.ReadCommitted)
		if err1 != nil {
			return
		}
		defer tx1.Rollback(ctx)
		// Update the topic
		_, err := r.Topics.UpdateDescription(ctx, topicID, newDescription, repo.WithTx(tx1))
		if err != nil {
			err1 = err
			return
		}
		err = tx1.Commit(t.Context())
		if err != nil {
			err1 = err
			return
		}
		close(ch)
	}()

	// Start transaction 2 (reader)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ch
		tx2, err2 = r.BeginTx(ctx, repo.ReadCommitted)
		if err2 != nil {
			return
		}
		defer tx2.Rollback(ctx)
		// Try to read the topic - should see the uncommitted change
		topic, err2 := r.Topics.GetByID(ctx, topicID, repo.WithTx(tx2))
		if err2 != nil {
			return
		}
		// In ReadCommitted, we should see the uncommitted change
		if topic.Description != newDescription {
			t.Errorf("Expected to see uncommitted change in ReadCommitted, got: %s", topic.Description)
		}
	}()

	wg.Wait()
	if err1 != nil {
		t.Errorf("Transaction 1 failed: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Transaction 2 failed: %v", err2)
	}
}

// testRepeatableReadIsolation tests that RepeatableRead prevents dirty reads
func testRepeatableReadIsolation(t *testing.T, r *repo.Repository, topicID string) {
	ctx := context.Background()
	var wg sync.WaitGroup
	var err1, err2 error

	// Channels for proper synchronization
	tx1Started := make(chan struct{})
	tx1Updated := make(chan struct{})
	tx2Started := make(chan struct{})

	// Reset the topic description first
	originalDescription := "Original description"
	newDescription := "Updated by TX1"
	_, err := r.Topics.UpdateDescription(ctx, topicID, originalDescription)
	if err != nil {
		t.Fatalf("Failed to reset topic description: %v", err)
	}

	// Start transaction 1 (updater)
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx1, err := r.BeginTx(ctx, repo.RepeatableRead)
		if err != nil {
			err1 = err
			return
		}
		defer tx1.Rollback(ctx)

		// Signal that TX1 has started
		close(tx1Started)

		// Wait for TX2 to start
		<-tx2Started

		// Update the topic
		_, err = r.Topics.UpdateDescription(ctx, topicID, newDescription, repo.WithTx(tx1))
		if err != nil {
			err1 = err
			return
		}

		// Signal that TX1 has updated
		close(tx1Updated)
	}()

	// Start transaction 2 (reader)
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Wait for TX1 to start first
		<-tx1Started

		tx2, err := r.BeginTx(ctx, repo.RepeatableRead)
		if err != nil {
			err2 = err
			return
		}
		defer tx2.Rollback(ctx)

		// Signal that TX2 has started
		close(tx2Started)

		// Wait for TX1 to update
		<-tx1Updated

		// Now read the topic - should NOT see the uncommitted change
		topic, err := r.Topics.GetByID(ctx, topicID, repo.WithTx(tx2))
		if err != nil {
			err2 = err
			return
		}

		// In RepeatableRead, we should NOT see the uncommitted change
		if topic.Description != originalDescription {
			t.Errorf("Expected to NOT see uncommitted change in RepeatableRead, got: %s", topic.Description)
		}
	}()

	wg.Wait()
	if err1 != nil {
		t.Errorf("Transaction 1 failed: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Transaction 2 failed: %v", err2)
	}
}

// testSerializableIsolation tests that Serializable prevents phantom reads
func testSerializableIsolation(t *testing.T, r *repo.Repository, topicID string) {
	ctx := context.Background()
	var wg sync.WaitGroup
	var err1, err2 error

	// Channels for proper synchronization
	tx1Started := make(chan struct{})
	tx1FirstRead := make(chan struct{})
	tx2Inserted := make(chan struct{})

	// Start transaction 1 (reader that will be affected by phantom reads)
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx1, err := r.BeginTx(ctx, repo.Serializable)
		if err != nil {
			err1 = err
			return
		}
		defer tx1.Rollback(ctx)

		// Signal that TX1 has started
		close(tx1Started)

		// First read - count topics
		topics1, err := r.Topics.List(ctx, 0, 0, repo.WithTx(tx1))
		if err != nil {
			err1 = err
			return
		}
		count1 := len(topics1.Data)

		// Signal that TX1 has completed first read
		close(tx1FirstRead)

		// Wait for TX2 to insert a new topic
		<-tx2Inserted

		// Second read - should see the same count in Serializable
		topics2, err := r.Topics.List(ctx, 0, 0, repo.WithTx(tx1))
		if err != nil {
			err1 = err
			return
		}
		count2 := len(topics2.Data)

		// In Serializable, both reads should return the same count
		if count1 != count2 {
			t.Errorf("Expected same count in Serializable isolation, got: %d vs %d", count1, count2)
		}
	}()

	// Start transaction 2 (inserter)
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Wait for TX1 to start first
		<-tx1Started
		// Wait for TX1 to complete first read
		<-tx1FirstRead
		tx2, err := r.BeginTx(ctx, repo.Serializable)
		if err != nil {
			err2 = err
			return
		}
		defer tx2.Rollback(ctx)

		// Create a new topic with unique ID
		newTopic := factcheck.Topic{
			ID:           fmt.Sprintf("550e8400-e29b-41d4-a716-44665544%04d", time.Now().UnixNano()%10000),
			Name:         "Phantom Topic",
			Description:  "This topic should not be visible to TX1",
			Status:       factcheck.StatusTopicPending,
			Result:       "",
			ResultStatus: factcheck.StatusTopicResultNone,
			CreatedAt:    utils.TimeNow(),
			UpdatedAt:    nil,
		}
		_, err = r.Topics.Create(ctx, newTopic, repo.WithTx(tx2))
		if err != nil {
			err2 = err
			return
		}
		// Commit the transaction
		err = tx2.Commit(ctx)
		if err != nil {
			err2 = err
			return
		}
		// Signal that TX2 has inserted and committed
		close(tx2Inserted)
	}()

	wg.Wait()
	if err1 != nil {
		t.Errorf("Transaction 1 failed: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Transaction 2 failed: %v", err2)
	}
}

// testConcurrentUpdates tests concurrent updates to the same resource
func testConcurrentUpdates(t *testing.T, r *repo.Repository, topicID string) {
	ctx := context.Background()
	var wg sync.WaitGroup
	var success1, success2 bool

	tx1Started := make(chan struct{})
	tx1Read := make(chan struct{})
	tx2Started := make(chan struct{})
	tx2Read := make(chan struct{})
	bothRead := make(chan struct{})

	// Reset the topic description first
	_, err := r.Topics.UpdateDescription(ctx, topicID, "Original description")
	if err != nil {
		t.Fatalf("Failed to reset topic description: %v", err)
	}

	// Start transaction 1
	wg.Add(1)
	go func() {
		defer wg.Done()
		tx1, err := r.BeginTx(ctx, repo.Serializable)
		if err != nil {
			return
		}
		defer func() {
			if !success1 {
				tx1.Rollback(ctx)
			}
		}()

		// Signal that TX1 has started
		close(tx1Started)

		// Wait for TX2 to start
		<-tx2Started

		// Read the topic
		_, err = r.Topics.GetByID(ctx, topicID, repo.WithTx(tx1))
		if err != nil {
			return
		}

		// Signal that TX1 has read
		close(tx1Read)

		// Wait for both transactions to have read
		<-bothRead

		// Update the topic
		_, err = r.Topics.UpdateDescription(ctx, topicID, "Updated by TX1", repo.WithTx(tx1))
		if err != nil {
			return
		}

		// Try to commit
		err = tx1.Commit(ctx)
		if err == nil {
			success1 = true
		}
	}()

	// Start transaction 2
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Wait for TX1 to start first
		<-tx1Started

		tx2, err := r.BeginTx(ctx, repo.Serializable)
		if err != nil {
			return
		}
		defer func() {
			if !success2 {
				tx2.Rollback(ctx)
			}
		}()

		// Signal that TX2 has started
		close(tx2Started)

		// Wait for TX1 to read
		<-tx1Read

		// Read the topic
		_, err = r.Topics.GetByID(ctx, topicID, repo.WithTx(tx2))
		if err != nil {
			return
		}

		// Signal that TX2 has read
		close(tx2Read)

		// Wait for both transactions to have read
		<-bothRead

		// Update the topic
		_, err = r.Topics.UpdateDescription(ctx, topicID, "Updated by TX2", repo.WithTx(tx2))
		if err != nil {
			return
		}

		// Try to commit
		err = tx2.Commit(ctx)
		if err == nil {
			success2 = true
		}
	}()

	// Wait for both transactions to read, then signal them to proceed
	go func() {
		<-tx1Read
		<-tx2Read
		close(bothRead)
	}()

	wg.Wait()

	// In Serializable isolation, one transaction should succeed and one should fail
	// due to serialization failure
	if success1 && success2 {
		t.Error("Both transactions succeeded, expected one to fail due to serialization conflict")
	}

	if !success1 && !success2 {
		t.Error("Both transactions failed, expected one to succeed")
	}

	// Check final state
	finalTopic, err := r.Topics.GetByID(ctx, topicID)
	if err != nil {
		t.Fatalf("Failed to get final topic state: %v", err)
	}

	expectedDescriptions := []string{"Updated by TX1", "Updated by TX2"}
	found := false
	for _, expected := range expectedDescriptions {
		if finalTopic.Description == expected {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Final topic description '%s' not in expected values: %v", finalTopic.Description, expectedDescriptions)
	}
}

// TestTransactionRollback tests that rollback works correctly
func TestTransactionRollback(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()

	ctx := context.Background()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create a topic
	originalTopic := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440005",
		Name:         "Rollback Test Topic",
		Description:  "Original description",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	createdTopic, err := app.Repository.Topics.Create(ctx, originalTopic)
	if err != nil {
		t.Fatalf("Failed to create topic: %v", err)
	}

	// Start a transaction
	tx, err := app.Repository.Begin(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Update the topic within the transaction
	_, err = app.Repository.Topics.UpdateDescription(ctx, createdTopic.ID, "Modified description", repo.WithTx(tx))
	if err != nil {
		t.Fatalf("Failed to update topic in transaction: %v", err)
	}

	// Verify the change is visible within the transaction
	topicInTx, err := app.Repository.Topics.GetByID(ctx, createdTopic.ID, repo.WithTx(tx))
	if err != nil {
		t.Fatalf("Failed to get topic in transaction: %v", err)
	}

	if topicInTx.Description != "Modified description" {
		t.Errorf("Expected modified description in transaction, got: %s", topicInTx.Description)
	}

	// Rollback the transaction
	err = tx.Rollback(ctx)
	if err != nil {
		t.Fatalf("Failed to rollback transaction: %v", err)
	}

	// Verify the change is NOT visible outside the transaction
	topicAfterRollback, err := app.Repository.Topics.GetByID(ctx, createdTopic.ID)
	if err != nil {
		t.Fatalf("Failed to get topic after rollback: %v", err)
	}

	if topicAfterRollback.Description != "Original description" {
		t.Errorf("Expected original description after rollback, got: %s", topicAfterRollback.Description)
	}
}

// TestTransactionCommit tests that commit works correctly
func TestTransactionCommit(t *testing.T) {
	app, cleanup, err := di.InitializeContainerTest()
	if err != nil {
		t.Fatalf("Failed to initialize test container: %v", err)
	}
	defer cleanup()

	ctx := context.Background()

	// Create test data
	now := utils.TimeNow().Round(0)
	utils.TimeFreeze(now)
	defer utils.TimeUnfreeze()

	// Create a topic
	originalTopic := factcheck.Topic{
		ID:           "550e8400-e29b-41d4-a716-446655440006",
		Name:         "Commit Test Topic",
		Description:  "Original description",
		Status:       factcheck.StatusTopicPending,
		Result:       "",
		ResultStatus: factcheck.StatusTopicResultNone,
		CreatedAt:    now,
		UpdatedAt:    nil,
	}

	createdTopic, err := app.Repository.Topics.Create(ctx, originalTopic)
	if err != nil {
		t.Fatalf("Failed to create topic: %v", err)
	}

	// Start a transaction
	tx, err := app.Repository.Begin(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	// Update the topic within the transaction
	_, err = app.Repository.Topics.UpdateDescription(ctx, createdTopic.ID, "Committed description", repo.WithTx(tx))
	if err != nil {
		t.Fatalf("Failed to update topic in transaction: %v", err)
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	// Verify the change is visible outside the transaction
	topicAfterCommit, err := app.Repository.Topics.GetByID(ctx, createdTopic.ID)
	if err != nil {
		t.Fatalf("Failed to get topic after commit: %v", err)
	}

	if topicAfterCommit.Description != "Committed description" {
		t.Errorf("Expected committed description after commit, got: %s", topicAfterCommit.Description)
	}
}
