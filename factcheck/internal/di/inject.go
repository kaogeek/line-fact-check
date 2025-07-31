//go:build wireinject
// +build wireinject

package di

import "github.com/google/wire"

func InitializeContainerTest() (Container, func(), error) {
	wire.Build(ProviderSetTest)
	return Container{}, nil, nil
}
