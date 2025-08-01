package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexflint/go-arg"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
)

type cli struct {
	// Submit submits new message as user
	Submit *cmdSubmit `arg:"subcommand:submit"`

	// Dump dumps stuff in DB
	Dump *cmdDump `arg:"subcommand:dump"`
}

type cmdSubmit struct {
	Text    string `arg:"positional"`
	TopicID string `arg:"positional"`
}

type cmdDump struct {
	Tables []string `arg:"positional"`
}

//nolint:noctx
func main() {
	c := cli{}
	err := arg.Parse(&c)
	if err != nil {
		panic(err)
	}
	container, cleanup, err := di.InitializeContainer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	quit := make(chan os.Signal, 1) // Buffered so it won't block on 2x Ctrl-C
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	ctx := context.Background()

	go func() {
		switch {
		case c.Submit != nil:
			err := c.Submit.submit(ctx, container.Config.HTTP, c.Submit.Text, c.Submit.TopicID)
			if err != nil {
				panic(err)
			}
		case c.Dump != nil:
			topics, err := container.Repository.Topics.List(ctx, 0, 0)
			if err != nil {
				panic(err)
			}
			slog.Info("dumping topics", "data", topics)
		}
	}()

	<-quit
	slog.Info("[main] exitting")
}
