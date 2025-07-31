package di

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

// ContainerTest is a container for testing.
// In addition to the usual cleanup functions provided by [Container],
// it also provides cleanup function to clear all data from the database
// BEFORE and AFTER each test.
type ContainerTest Container

func NewTest(
	conf config.Config,
	db postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	service core.Service,
) (
	ContainerTest,
	func(),
) {
	clearData(db, "init")
	cleanup := func() {
		clearData(db, "teardown")
	}
	return ContainerTest(New(
		conf,
		db,
		querier,
		repo,
		service,
	)), cleanup
}

func clearData(conn postgres.DBTX, stage string) {
	tables := [4]string{
		"topics",
		"messages_v2",
		"message_groups",
		"answers",
	}
	ctx := context.Background()
	slog.WarnContext(ctx, "Clearing all data from database", "stage", stage)
	for i, t := range tables {
		err := delete(ctx, conn, t)
		if err != nil {
			slog.Error("failed to delete table", "i", i, "table", t)
			panic(err)
		}
	}
	slog.WarnContext(ctx, "Cleared all data from database", "stage", stage)
}

func delete(ctx context.Context, conn postgres.DBTX, table string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table))
	return err
}
