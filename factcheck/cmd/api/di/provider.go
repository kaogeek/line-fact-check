package di

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handlers"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

var ProviderSetMain = wire.NewSet(
	ProviderSetDatabase,
	repo.New,
	handlers.New,
	New,
)

var ProviderSetDatabase = wire.NewSet(
	wire.Bind(new(postgres.DBTX), new(*pgx.Conn)),
	wire.Bind(new(postgres.Querier), new(*postgres.Queries)),
	config.New,
	postgres.New,
	postgres.NewConn,
)
