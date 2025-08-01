// Package di provides common provider sets for factcheck programs,
// like database connection, data repository, core service object, etc.
package di

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

// ProviderSet provides all of internal objects
var ProviderSet = wire.NewSet(
	config.New,
	ProviderSetDatabase,
	ProviderSetRepo,
	ProviderSetCore,
	wire.Struct(new(Container), "*"),
)

// ProviderSetTest provides all of internal objects, including ContainerTest
var ProviderSetTest = wire.NewSet(
	config.NewTest,
	ProviderSetDatabase,
	ProviderSetRepo,
	ProviderSetCore,
	NewTest,
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
