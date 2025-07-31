package di

import (
	"net/http"

	"github.com/google/wire"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
)

// ProviderSet provides everything cmd/api needs
var ProviderSet = wire.NewSet(
	wire.Bind(new(server.Server), new(*http.Server)),
	di.ProviderSet,
	handler.New,
	server.New,
	wire.Struct(new(Container), "*"),
)

// ProviderSetTest provides everything [ProviderSet] does,
// but with its own cleanup functions from internal di
var ProviderSetTest = wire.NewSet(
	wire.Bind(new(server.Server), new(*http.Server)),
	di.ProviderSetTest,
	handler.New,
	server.New,
	wire.Struct(new(Container), "*"),
)
