// Package di is where we define Wire DI.
package di

import (
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Container struct {
	Conf            config.Config
	PostgresConn    postgres.DBTX
	PostgresQuerier postgres.Querier
	Repository      repo.Repository
	Service         core.Service
	Handler         handler.Handler
	Server          server.Server
}

func New(
	conf config.Config,
	conn postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	service core.Service,
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
