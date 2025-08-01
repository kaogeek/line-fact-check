package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck/internal/config"
)

func (c *cmdCreateTopic) createTopic(
	ctx context.Context,
	conf config.HTTP,
) error {
	url := hostURL(conf.ListenAddr)
	url += "/topics"
	slog.Info("got new url", "config_http", conf, "url", url)
	data := map[string]any{
		"name":        c.Name,
		"description": c.Description,
	}
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(data)
	if err != nil {
		panic(err)
	}
	return callHTTP(ctx, http.MethodPost, url, body)
}

func (c *cmdAssignTopic) assignTopic(
	ctx context.Context,
	conf config.HTTP,
) error {
	url := hostURL(conf.ListenAddr)
	url += "/admin/message-groups/assign/" + c.GroupID
	slog.Info("got new url", "config_http", conf, "url", url)
	data := map[string]any{
		"topic_id": c.TopicID,
	}
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(data)
	if err != nil {
		panic(err)
	}
	return callHTTP(ctx, http.MethodPut, url, body)
}
