//go:build wireinject
// +build wireinject

package di

import "github.com/google/wire"

func InitializeContainerTest() (ContainerTest, func(), error) {
	wire.Build(ProviderSetTest)
	return ContainerTest{}, nil, nil
}
