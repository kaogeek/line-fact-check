// Package di is where we define Wire DI.
package di

import (
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/handler"
	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/internal/server"
	"github.com/kaogeek/line-fact-check/factcheck/internal/di"
)

type Container struct {
	di.Container
	Handler handler.Handler
	Server  server.Server
}

func New(
	container di.Container,
	handler handler.Handler,
	server server.Server,
) Container {
	return Container{
		Container: container,
		Handler:   handler,
		Server:    server,
	}
}

func NewTest(
	container di.Container,
	handler handler.Handler,
	server server.Server,
) Container {
	return Container{
		Container: container,
		Handler:   handler,
		Server:    server,
	}
}
