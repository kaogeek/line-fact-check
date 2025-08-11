//go:build wireinject
// +build wireinject

package ittest

import (
	"testing"

	"github.com/google/wire"

	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
)

func InitializeContainer() (di.Container, func(), error) {
	wire.Build(ProviderSet)
	return di.Container{}, nil, nil
}

func InitializeContainerTest(t *testing.T) (di.Container, func(), error) {
	wire.Build(ProviderSetTest)
	return di.Container{}, nil, nil
}
