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
			return
		}
		errNotFound(w, fmt.Sprintf("%s not found: %s", resourceType, filter))
		return
	}
	errInternalError(w, err.Error())
}

func errNotFound(w http.ResponseWriter, id string) {
	w.WriteHeader(http.StatusNotFound)
	contentTypeText(w.Header())
	fmt.Fprintf(w, "not found: %s", id)
}

func errInternalError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	contentTypeText(w.Header())
	fmt.Fprintf(w, "server error: %s", err)
}

func errBadRequest(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	contentTypeText(w.Header())
	fmt.Fprintf(w, "bad request: %s", err)
}

func contentTypeJSON(h http.Header) {
	h.Add("Content-Type", "application/json; charset=utf-8")
}

func contentTypeText(h http.Header) {
	h.Add("Content-Type", "text/plain; charset=utf-8")
}
