package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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
	return callHTTP(ctx, http.MethodPost, url, body)
}
