package di

import (
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
	"github.com/kaogeek/line-fact-check/factcheck/internal/core"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

// ContainerTest is a container for testing.
// In addition to the usual cleanup functions provided by [Container],
// it also provides cleanup function to clear all data from the database
// BEFORE and AFTER each test.
//
// They are defined separate just in case we need something
type ContainerTest Container

func NewTestV2(
	conf config.Config,
	conn postgres.DBTX,
	querier postgres.Querier,
	repo repo.Repository,
	service core.Service,
	handler handler.Handler,
	server server.Server,
) (
	ContainerTest,
	func(),
) {
	return ContainerTest(New(
			conf,
			conn,
			querier,
			repo,
			service,
			handler,
			server,
		)), func() {
			slog.Debug("containerTest cleanup")
		}
}
