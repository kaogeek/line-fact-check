package di

import (
	"context"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

// ContainerTest is a container for testing.
// In addition to the usual cleanup functions provided by [Container],
// it also provides cleanup function to clear all data from the database
// BEFORE and AFTER each test.
type ContainerTest Container

func NewTest(
	conf config.Config,
	conn postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	handler handler.Handler,
	server server.Server,
) (
	ContainerTest,
	func(),
) {
	clearData(conn, "init")
	cleanup := func() {
		clearData(conn, "teardown")
	}
	return ContainerTest(New(
		conf,
		conn,
		querier,
		repo,
		handler,
		server,
	)), cleanup
}

func clearData(conn postgres.DBTX, stage string) {
	slog.Warn("Clearing all data from database", "stage", stage)

	ctx := context.Background()
	_, err := conn.Exec(ctx, "DELETE FROM topics")
	if err != nil {
		slog.Error("Failed to delete topics", "error", err)
		panic(err)
	}
	_, err = conn.Exec(ctx, "DELETE FROM messages")
	if err != nil {
		slog.Error("Failed to delete messages", "error", err)
		panic(err)
	}
	_, err = conn.Exec(ctx, "DELETE FROM user_messages")
	if err != nil {
		slog.Error("Failed to delete user_messages", "error", err)
		panic(err)
	}

	slog.Warn("Cleared all data from database", "stage", stage)
}
