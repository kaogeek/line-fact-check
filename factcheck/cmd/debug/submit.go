package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
)

func hostURL(u string) string {
	if strings.HasPrefix(u, ":") {
		u = fmt.Sprintf("localhost%s", u)
	}
	return fmt.Sprintf("http://%s", u)
}

//nolint:noctx
func (c *cmdSubmit) submit(
	ctx context.Context,
	conf config.HTTP,
	text string,
	topicID string,
) error {
	url := hostURL(conf.ListenAddr)
	url += "/messages/"
	slog.Info("got new url", "config_http", conf, "url", url)
	data := struct {
		Text    string `json:"text"`
		TopicID string `json:"topic_id"`
	}{
		Text:    text,
		TopicID: topicID,
	}
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(data)
	if err != nil {
		panic(err)
	}
	slog.Info("body to be sent", "body", body.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return err
	}
	client := http.Client{
		Timeout: time.Second * 2,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBody := bytes.NewBuffer(nil)
	_, err = io.Copy(respBody, resp.Body)
	if err != nil {
		panic(err)
	}

	slog.Info("got response",
		"url", url,
		"status_code", resp.StatusCode,
		"body", respBody.String(),
	)
	return nil
}
