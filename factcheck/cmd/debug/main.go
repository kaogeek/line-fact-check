package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/alexflint/go-arg"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/di"
)

type cli struct {
	Submit      *cmdSubmit      `arg:"subcommand:submit"` // Submit submits new message as user
	CreateTopic *cmdCreateTopic `arg:"subcommand:create-topic"`
	AssignTopic *cmdAssignTopic `arg:"subcommand:assign-topic"` // AssignTopic assigns a msg group to topicID
	Dump        *cmdDump        `arg:"subcommand:dump"`         // Dump dumps stuff in DB
}

type cmdSubmit struct {
	Text    string `arg:"positional"`
	TopicID string `arg:"positional"`
}

type cmdAssignTopic struct {
	GroupID string `arg:"positional"`
	TopicID string `arg:"positional"`
}

type cmdCreateTopic struct {
	Name        string `arg:"positional"`
	Description string `arg:"positional"`
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

	ctx := context.Background()
	conf := container.Config
	switch {
	case c.Submit != nil:
		text := c.Submit.Text
		topicID := c.Submit.TopicID
		if text == "" {
			text = "dummy debug text"
		}
		err := c.Submit.submit(ctx, conf.HTTP, text, topicID)
		if err != nil {
			panic(err)
		}

	case c.CreateTopic != nil:
		c.CreateTopic.createTopic(ctx, conf.HTTP)

	case c.AssignTopic != nil:
		c.AssignTopic.assignTopic(ctx, conf.HTTP)

	case c.Dump != nil:
		topics, err := container.Repository.Topics.List(ctx, 0, 0)
		if err != nil {
			panic(err)
		}
		slog.Info("dumping topics", "data", topics)
	}
}

func callHTTP(
	ctx context.Context,
	method string,
	url string,
	body *bytes.Buffer,
) error {
	slog.InfoContext(ctx, "callHTTP: start", "url", url)
	slog.InfoContext(ctx, "callHTTP: body to be sent", "body", body.String())
	req, err := http.NewRequestWithContext(ctx, method, url, body)
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
	slog.InfoContext(ctx, "callHTTP: got response",
		"url", url,
		"status_code", resp.StatusCode,
		"body", respBody.String(),
	)
	return nil
}
