//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
)

// InitializeHandler returns our API handler.
func InitializeHandler() (handler.Handler, func(), error) {
	wire.Build(ProviderSet)
	return nil, nil, nil
}

// InitializeContainer returns all components of interest,
// perfect for integration test or debugging
func InitializeContainer() (Container, func(), error) {
	wire.Build(ProviderSet)
	return Container{}, nil, nil
}
