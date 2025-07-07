package di

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handlers"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

var ProviderSet = wire.NewSet(
	ProviderSetBase,
	repo.New,
	handlers.New,
	New,
)

var ProviderSetTest = wire.NewSet(
	ProviderSetBaseTest,
	repo.New,
	handlers.New,
	New,
)

var ProviderSetBase = wire.NewSet(
	wire.Bind(new(postgres.DBTX), new(*pgx.Conn)),
	wire.Bind(new(postgres.Querier), new(*postgres.Queries)),
	config.New,
	postgres.New,
	postgres.NewConn,
)

var ProviderSetBaseTest = wire.NewSet(
	wire.Bind(new(postgres.DBTX), new(*pgx.Conn)),
	wire.Bind(new(postgres.Querier), new(*postgres.Queries)),
	config.NewTest,
	postgres.New,
	postgres.NewConn,
)
