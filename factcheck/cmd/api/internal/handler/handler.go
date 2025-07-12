package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Handler interface {
	CreateTopic(http.ResponseWriter, *http.Request)
	ListTopics(http.ResponseWriter, *http.Request)
	GetTopicByID(http.ResponseWriter, *http.Request)
	DeleteTopicByID(http.ResponseWriter, *http.Request)
	UpdateTopicStatus(http.ResponseWriter, *http.Request)
	UpdateTopicDescription(http.ResponseWriter, *http.Request)
	UpdateTopicName(http.ResponseWriter, *http.Request)

	CreateMessage(http.ResponseWriter, *http.Request)
	ListMessagesByTopicID(http.ResponseWriter, *http.Request)
	DeleteMessageByID(http.ResponseWriter, *http.Request)
}

type handler struct {
	topics   repo.Topics
	messages repo.Messages
}

func New(repo repo.Repository) Handler {
	return &handler{
		topics:   repo.Topic,
		messages: repo.Message,
	}
}

// handleNotFound standardizes not-found error handling in handlers
func handleNotFound(w http.ResponseWriter, err error, resourceType string, filter string) {
	if repo.IsNotFound(err) {
		var notFoundErr *repo.ErrNotFound
		if errors.As(err, &notFoundErr) {
			errNotFound(w, notFoundErr.Error())
		} else {
			errNotFound(w, fmt.Sprintf("%s not found: %s", resourceType, filter))
		}
		return
	}
	errInternalError(w, err.Error())
}
