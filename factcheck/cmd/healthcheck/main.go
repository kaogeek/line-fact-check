package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	addr := os.Getenv("FACTCHECKAPI_LISTEN_ADDRESS")
	timeoutEnv := os.Getenv("FACTCHECKAPI_TIMEOUTMS_WRITE")
	timeoutMs, err := strconv.Atoi(timeoutEnv)
	if err != nil {
		timeoutMs = 2000
	}
	c := http.Client{
		Timeout: time.Millisecond * time.Duration(timeoutMs),
	}

	url := fmt.Sprintf("http://0.0.0.0%s/health", addr)
	slog.Info("healthcheck", "url", url, "FACTCHECKAPI_LISTEN_ADDRESS", addr, "timeout_ms", timeoutMs)

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
