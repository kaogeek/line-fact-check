package di

import (
	"github.com/google/wire"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handlers"
	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
	"github.com/kaogeek/line-fact-check/factcheck/models/postgres"
)

var ProviderSetAPI = wire.NewSet(
	postgres.New,
	repo.New,
	handlers.New,
)
