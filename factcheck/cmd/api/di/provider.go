package di

import (
	"net/http"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/internal/service"
)

var ProviderSet = wire.NewSet(
	wire.Bind(new(server.Server), new(*http.Server)),
	wire.Bind(new(service.Service), new(service.ServiceFactcheck)),
	ProviderSetBase,
	repo.New,
	service.New,
	handler.New,
	server.New,
	New,
)

var ProviderSetTest = wire.NewSet(
	wire.Bind(new(server.Server), new(*http.Server)),
	wire.Bind(new(service.Service), new(service.ServiceFactcheck)),
	ProviderSetBaseTest,
	repo.New,
	service.New,
	handler.New,
	server.New,
	NewTest,
)

var ProviderSetBase = wire.NewSet(
	wire.Bind(new(postgres.DBTX), new(*pgxpool.Pool)),
	wire.Bind(new(postgres.Querier), new(*postgres.Queries)),
	config.New,
	postgres.New,
	postgres.NewConn,
)

var ProviderSetBaseTest = wire.NewSet(
	wire.Bind(new(postgres.DBTX), new(*pgxpool.Pool)),
	wire.Bind(new(postgres.Querier), new(*postgres.Queries)),
	config.NewTest,
	postgres.New,
	postgres.NewConn,
)
