// Package di is where we define Wire DI.
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

type Container struct {
	Conf            config.Config
	PostgresConn    postgres.DBTX
	PostgresQuerier postgres.Querier
	Repository      repo.Repository
	Handler         handler.Handler
	Server          server.Server
}

type ContainerTest Container

func New(
	conf config.Config,
	conn postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	handler handler.Handler,
	server server.Server,
) Container {
	return Container{
		Conf:            conf,
		PostgresConn:    conn,
		PostgresQuerier: querier,
		Repository:      repo,
		Handler:         handler,
		Server:          server,
	}
}

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
		panic(err)
	}
	_, err = conn.Exec(ctx, "DELETE FROM messages")
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx, "DELETE FROM user_messages")
	if err != nil {
		panic(err)
	}

	slog.Warn("Cleared all data from database", "stage", stage)
}
