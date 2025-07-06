//go:build wireinject
// +build wireinject

package di

import "github.com/google/wire"

func InitializeContainer() (Container, func(), error) {
	wire.Build(ProviderSetMain)
	return Container{}, nil, nil
}
