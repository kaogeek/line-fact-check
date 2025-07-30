package di

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

// ProviderSetDatabase provides all database-related dependencies
var ProviderSetDatabase = wire.NewSet(
	wire.Bind(new(postgres.DBTX), new(*pgxpool.Pool)),
	wire.Bind(new(postgres.Querier), new(*postgres.Queries)),
	postgres.New,
	postgres.NewConn,
)

// ProviderSetRepo provides repository layer
var ProviderSetRepo = wire.NewSet(
	repo.New,
)

// ProviderSetCore provides business logic layer
var ProviderSetCore = wire.NewSet(
	wire.Bind(new(core.Service), new(core.ServiceFactcheck)),
	core.New,
)
