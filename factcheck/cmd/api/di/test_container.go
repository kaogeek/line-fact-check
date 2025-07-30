package di

import (
	"context"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
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
	service core.Service,
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
		service,
		handler,
		server,
	)), cleanup
}

func clearData(conn postgres.DBTX, stage string) {
	ctx := context.Background()
	slog.WarnContext(ctx, "Clearing all data from database", "stage", stage)

	_, err := conn.Exec(ctx, "DELETE FROM topics")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete topics", "error", err)
		panic(err)
	}
	_, err = conn.Exec(ctx, "DELETE FROM messages")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete messages", "error", err)
		panic(err)
	}
	_, err = conn.Exec(ctx, "DELETE FROM user_messages")
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete user_messages", "error", err)
		panic(err)
	}

	slog.WarnContext(ctx, "Cleared all data from database", "stage", stage)
}
