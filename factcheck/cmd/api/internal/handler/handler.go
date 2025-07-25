package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/kaogeek/line-fact-check/factcheck/internal/repo"
)

type Handler interface {
	CreateTopic(http.ResponseWriter, *http.Request)
	ListAllTopics(http.ResponseWriter, *http.Request)
	ListTopicsHome(http.ResponseWriter, *http.Request)
	CountTopicsHome(http.ResponseWriter, *http.Request)
	GetTopicByID(http.ResponseWriter, *http.Request)
	DeleteTopicByID(http.ResponseWriter, *http.Request)
	UpdateTopicStatus(http.ResponseWriter, *http.Request)
	UpdateTopicDescription(http.ResponseWriter, *http.Request)
	UpdateTopicName(http.ResponseWriter, *http.Request)

	CreateMessage(http.ResponseWriter, *http.Request)
	ListMessagesByTopicID(http.ResponseWriter, *http.Request)
	DeleteMessageByID(http.ResponseWriter, *http.Request)

	NewUserMessage(http.ResponseWriter, *http.Request)
}

type handler struct {
	repository repo.Repository
	topics     repo.Topics
	messagesv2 repo.MessagesV2
	groups     repo.MessagesV2Groups

	// TO BE DEPRECATED

	messages     repo.Messages
	userMessages repo.UserMessages
}

func New(repo repo.Repository) Handler {
	return &handler{
		repository: repo,
		topics:     repo.Topics,
		messagesv2: repo.MessagesV2,
		groups:     repo.MessagesV2Groups,

		messages:     repo.Messages,
		userMessages: repo.UserMessages,
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

func paramID(r *http.Request) string {
	return chi.URLParam(r, "id")
}

func limitOffSet(r *http.Request) (int, int, error) {
	query := r.URL.Query().Get
	queryLimit := query("limit")
	queryOffset := query("offset")
	limit, offset := 0, 0
	var err error
	if queryLimit != "" {
		limit, err = strconv.Atoi(queryLimit)
		if err != nil {
			return 0, 0, fmt.Errorf("bad query limit: '%s'", queryLimit)
		}
	}
	if queryOffset != "" {
		offset, err = strconv.Atoi(queryOffset)
		if err != nil {
			return 0, 0, fmt.Errorf("bad query offset: '%s'", queryOffset)
		}
	}
	return limit, offset, nil
}

func decode[T any](r *http.Request) (T, error) {
	var t T
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		var zero T
		return zero, err
	}
	return t, nil
}

func sendText(w http.ResponseWriter, text string, status int) {
	w.WriteHeader(status)
	contentTypeText(w.Header())
	_, err := w.Write([]byte(text))
	if err != nil {
		slog.Error("error writing to response", "error", err)
	}
}

// sendJSON calls replyJsonError, and on non-nil error, writes 500 response
func sendJSON(w http.ResponseWriter, data any, status int) {
	err := replyJSON(w, data, status)
	if err != nil {
		errInternalError(w, err.Error())
	}
}

// replyJSON marshals data into JSON string before writing response.
// If marshaling failed, the response is left untouched and the error is returned.
func replyJSON(w http.ResponseWriter, data any, status int) error {
	j, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal json error: %w", err)
	}

	w.WriteHeader(status)
	contentTypeJSON(w.Header())
	_, err = w.Write(j)
	if err != nil {
		slog.Error("error writing to response", "error", err)
		return fmt.Errorf("write to response error: %w", err)
	}
	return nil
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
