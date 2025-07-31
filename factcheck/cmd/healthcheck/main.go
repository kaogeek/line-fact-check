package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	addr := conf.HTTP.ListenAddr
	url := fmt.Sprintf("http://0.0.0.0%s/health", addr) // TODO: port or addr config?
	timeoutMsRead := time.Millisecond * time.Duration(conf.HTTP.TimeoutMsRead)
	timeoutMsWrite := time.Millisecond * time.Duration(conf.HTTP.TimeoutMsWrite)
	timeout := timeoutMsRead + timeoutMsWrite
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c := http.Client{
		Timeout: timeoutMsRead + timeoutMsWrite,
	}

	slog.InfoContext(ctx, "healthcheck",
		"addr", addr,
		"url", url,
		"timeout_ms", timeout,
		"timeout_ms_read", timeoutMsRead,
		"timeout_ms_write", timeoutMsWrite,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			slog.ErrorContext(ctx, "error closing response body",
				"error", err,
				"addr", addr,
				"url", url,
				"timeout_ms", timeout,
				"timeout_ms_read", timeoutMsRead,
				"timeout_ms_write", timeoutMsWrite,
			)
		}
	}()
	if resp.StatusCode == http.StatusOK {
		return
	}
	slog.ErrorContext(ctx, "got wrong code",
		"actual", resp.StatusCode,
		"expected", http.StatusOK,
		"addr", addr,
		"url", url,
		"timeout_ms", timeout,
		"timeout_ms_read", timeoutMsRead,
		"timeout_ms_write", timeoutMsWrite,
	)
}
