package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type optionsCreate[T any] struct {
	check  func(context.Context, T) error // check checks if T is valid
	modify func(context.Context, T) T     // modify returns T to be created
}

type optionCreate[T any] func(*optionsCreate[T])

// createCheck allows f to inspect unmarshaled request body.
func createCheck[T any](f func(context.Context, T) error) optionCreate[T] {
	return func(c *optionsCreate[T]) { c.check = f }
}

// createModify allows f to inject values into data to be created.
// Common use case is assigning a new UUID and created_at.
func createModify[T any](f func(context.Context, T) T) optionCreate[T] {
	return func(c *optionsCreate[T]) { c.modify = f }
}

// create defines centralized behavior for creating an entry in the database.
// It allows fine-grained control via [optionsCreate].
//
// If you are implementing a more complex use case, e.g. request body differs from T,
// write your own handler to keep this function simple and stupid.
func create[T any](
	w http.ResponseWriter,
	r *http.Request,
	createFn func(context.Context, T) (T, error),
	opts ...optionCreate[T],
) {
	options := optionsCreate[T]{}
	for i := range opts {
		opts[i](&options)
	}
	data, err := decode[T](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	if options.check != nil {
		err := options.check(r.Context(), data)
		if err != nil {
			errBadRequest(w, err.Error())
			return
		}
	}
	if options.modify != nil {
		data = options.modify(r.Context(), data)
	}
	created, err := createFn(r.Context(), data)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, created, http.StatusCreated)
}

func list[T any](
	w http.ResponseWriter,
	r *http.Request,
	listFn func(context.Context) ([]T, error),
) {
	l, err := listFn(r.Context())
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, l, http.StatusOK)
}

// getBy uses getFn to get a T based on filter F.
func getBy[T any, F any](
	w http.ResponseWriter,
	r *http.Request,
	filter F,
	getFn func(ctx context.Context, filter F) (T, error),
) {
	data, err := getFn(r.Context(), filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errNotFound(w, fmt.Sprintf("not found for filter %+v: %s", filter, err.Error()))
			return
		}
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, data, http.StatusOK)
}

func deleteByID[T any](
	w http.ResponseWriter,
	r *http.Request,
	deleteFn func(context.Context, string) error,
) {
	id := paramID(r)
	if id == "" {
		errBadRequest(w, "empty id")
		return
	}
	err := deleteFn(r.Context(), id)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendText(w, "ok", http.StatusOK)
}

func paramID(r *http.Request) string {
	return chi.URLParam(r, "id")
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
