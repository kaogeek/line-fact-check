package di

import (
	"net/http"

	"github.com/google/wire"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
)

var ProviderSet = wire.NewSet(
	wire.Bind(new(server.Server), new(*http.Server)),
	config.New,
	di.ProviderSetDatabase,
	di.ProviderSetRepo,
	di.ProviderSetCore,
	di.New,
	handler.New,
	server.New,
	New,
)

var ProviderSetTestV2 = wire.NewSet(
	wire.Bind(new(server.Server), new(*http.Server)),
	config.NewTest,
	di.ProviderSetDatabase,
	di.ProviderSetRepo,
	di.ProviderSetCore,
	di.NewTest,
	handler.New,
	server.New,
	NewTestV2,
)
