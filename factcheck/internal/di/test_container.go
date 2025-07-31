package di

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

func NewTest(
	conf config.Config,
	db postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	service core.Service,
) (
	Container,
	func(),
) {
	clearData(db, "init")
	cleanup := func() {
		clearData(db, "teardown")
	}
	return New(
		conf,
		db,
		querier,
		repo,
		service,
	), cleanup
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
			slog.ErrorContext(ctx, "failed to delete table", "i", i, "table", t)
			panic(err)
		}
	}
	slog.WarnContext(ctx, "Cleared all data from database", "stage", stage)
}

func delete(ctx context.Context, conn postgres.DBTX, table string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table))
	return err
}
