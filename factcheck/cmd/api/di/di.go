// Package di is where we define Wire DI.
package di

import (
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handlers"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

type Container struct {
	PostgresConn    postgres.DBTX
	PostgresQuerier postgres.Querier
	Handler         handlers.Handler
}

func New(
	conn postgres.DBTX,
	querier postgres.Querier,
	handler handlers.Handler,
) Container {
	return Container{
		PostgresConn:    conn,
		PostgresQuerier: querier,
		Handler:         handler,
	}
}
