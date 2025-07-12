//go:build integration_test
// +build integration_test

package handler_test

import (
	"log/slog"
	"os"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
)

type TestSuite struct {
	container di.Container
}

func init() {
	// Set slog level to DEBUG and log with json formatter
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
}
