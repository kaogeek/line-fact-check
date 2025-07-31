package di

import (
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Container struct {
	Config          config.Config
	PostgresConn    postgres.DBTX
	PostgresQuerier postgres.Querier
	Repository      repo.Repository
	Service         core.Service
}

func New(
	conf config.Config,
	db postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	service core.Service,
) Container {
	return Container{
		Config:          conf,
		PostgresConn:    db,
		PostgresQuerier: querier,
		Repository:      repo,
		Service:         service,
	}
}
