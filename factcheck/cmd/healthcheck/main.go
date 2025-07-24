package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
)

func main() {
	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	c := http.Client{
		Timeout: time.Millisecond * time.Duration(conf.HTTP.TimeoutMsRead),
	}
	addr := conf.HTTP.ListenAddr
	url := fmt.Sprintf("http://0.0.0.0%s/health", addr) // TODO: port or addr config?

	slog.Info("healthcheck",
		"addr", addr,
		"url", url,
		"timeout_ms_read", conf.HTTP.TimeoutMsRead,
		"timeout_ms_write", conf.HTTP.TimeoutMsWrite,
	)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}
}
