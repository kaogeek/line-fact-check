package ittest

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/google/wire"

	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

// ProviderSet provides all of internal objects, including ContainerTest
// It uses background context for cleanup operations
var ProviderSet = wire.NewSet(
	config.NewTest,
	di.ProviderSetDatabase,
	di.ProviderSetRepo,
	di.ProviderSetCore,
	New,
)

// ProviderSetTest provides all of internal objects, including ContainerTest
// It uses [testing.T] to provide context for cleanup operations
var ProviderSetTest = wire.NewSet(
	config.NewTest,
	di.ProviderSetDatabase,
	di.ProviderSetRepo,
	di.ProviderSetCore,
	NewTest,
)

func New(
	conf config.Config,
	conn postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	service core.Service,
) (
	di.Container,
	func(),
) {
	clearData(conn, "init")
	cleanup := func() {
		clearData(conn, "teardown")
	}
	return di.Container{
		Config:          conf,
		PostgresConn:    conn,
		PostgresQuerier: querier,
		Repository:      repo,
		Service:         service,
	}, cleanup
}

func NewTest(
	t *testing.T,
	conf config.Config,
	conn postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	service core.Service,
) (
	di.Container,
	func(),
) {
	clearDataWithT(t, conn, "init")
	cleanup := func() {
		clearDataWithT(t, conn, "teardown")
	}
	return di.Container{
		Config:          conf,
		PostgresConn:    conn,
		PostgresQuerier: querier,
		Repository:      repo,
		Service:         service,
	}, cleanup
}

func tables() [4]string {
	return [4]string{
		"topics",
		"messages_v2",
		"message_groups",
		"answers",
	}
}

func clearData(conn postgres.DBTX, stage string) {
	ctx := context.Background()
	slog.WarnContext(ctx, "Clearing all data from database", "stage", stage)
	for i, table := range tables() {
		err := delete(ctx, conn, table)
		if err != nil {
			slog.ErrorContext(ctx, "failed to delete table", "i", i, "table", table)
			panic(err)
		}
	}
	slog.WarnContext(ctx, "Cleared all data from database", "stage", stage)
}

func clearDataWithT(t *testing.T, conn postgres.DBTX, stage string) {
	ctx := t.Context()
	t.Logf("Clearing all data from database: stage=%s", stage)
	for i, table := range tables() {
		err := delete(ctx, conn, table)
		if err != nil {
			slog.ErrorContext(ctx, "failed to delete table", "i", i, "table", table)
			t.Fatal("error deleting table", table, err)
		}
	}
	t.Logf("Cleared all data from database: stage=%s", stage)
}

func delete(ctx context.Context, conn postgres.DBTX, table string) error {
	_, err := conn.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table))
	return err
}
