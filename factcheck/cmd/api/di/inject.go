//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
)

// InitializeHandler returns our HTTP API server.
func InitializeServer() (server.Server, func(), error) {
	wire.Build(ProviderSet)
	return nil, nil, nil
}

// InitializeContainer returns all components of interest,
// perfect for integration test or debugging
func InitializeContainer() (Container, func(), error) {
	wire.Build(ProviderSet)
	return Container{}, nil, nil
}

func InitializeContainerTest() (Container, func(), error) {
	wire.Build(ProviderSetTest)
	return Container{}, nil, nil
}
