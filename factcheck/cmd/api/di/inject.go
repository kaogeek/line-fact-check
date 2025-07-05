//go:build wireinject
// +build wireinject

package di

import "github.com/google/wire"

func InitializeContainer() (Container, error) {
	wire.Build(ProviderSetMain)
	return Container{}, nil
}
