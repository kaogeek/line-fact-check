// Package di is where we define Wire DI.
package di

import (
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Container struct {
	PostgresConn    postgres.DBTX
	PostgresQuerier postgres.Querier
	Repository      repo.Repository
	Handler         handler.Handler
}

func New(
	conn postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	handler handler.Handler,
) Container {
	return Container{
		PostgresConn:    conn,
		PostgresQuerier: querier,
		Repository:      repo,
		Handler:         handler,
	}
}
