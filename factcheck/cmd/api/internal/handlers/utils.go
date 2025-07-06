package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type optionsCreate[T any] struct {
	check  func(context.Context, T) error
	modify func(context.Context, T) T
}

type optionCreate[T any] func(*optionsCreate[T])

func createCheck[T any](f func(context.Context, T) error) optionCreate[T] {
	return func(c *optionsCreate[T]) { c.check = f }
}
func createModify[T any](f func(context.Context, T) T) optionCreate[T] {
	return func(c *optionsCreate[T]) { c.modify = f }
}

func create[T any](
	w http.ResponseWriter,
	r *http.Request,
	repo interface {
		Create(context.Context, T) (T, error)
	},
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
	created, err := repo.Create(r.Context(), data)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, created, http.StatusCreated)
}

func list[T any](
	w http.ResponseWriter,
	r *http.Request,
	repo interface {
		List(context.Context) ([]T, error)
	},
) {
	l, err := repo.List(r.Context())
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	sendJSON(w, l, http.StatusOK)
}

func getByID[T any](
	w http.ResponseWriter,
	r *http.Request,
	repo interface {
		GetByID(context.Context, string) (T, error)
	},
) {
	id := paramID(r)
	if id == "" {
		errBadRequest(w, "empty id")
		return
	}
	data, err := repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errNotFound(w, id)
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
	repo interface {
		Delete(context.Context, string) error
	},
) {
	id := paramID(r)
	if id == "" {
		errBadRequest(w, "empty id")
		return
	}
	err := repo.Delete(r.Context(), id)
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
	var zero T
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return zero, err
	}
	var t T
	err = json.Unmarshal(b, &t)
	if err != nil {
		return zero, err
	}
	return t, nil
}

func sendText(w http.ResponseWriter, text string, status int) {
	contentTypeText(w.Header())
	w.WriteHeader(status)
	w.Write([]byte(text))
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
		return err
	}
	contentTypeJSON(w.Header())
	w.WriteHeader(status)
	w.Write(j)
	return nil
}

func errNotFound(w http.ResponseWriter, id string) {
	contentTypeText(w.Header())
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "not found: %s", id)
}

func errInternalError(w http.ResponseWriter, err string) {
	contentTypeText(w.Header())
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "server error: %s", err)
}

func errBadRequest(w http.ResponseWriter, err string) {
	contentTypeText(w.Header())
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "bad request: %s", err)
}

func contentTypeJSON(h http.Header) {
	h.Add("Content-Type", "application/json; charset=utf-8")
}

func contentTypeText(h http.Header) {
	h.Add("Content-Type", "text/plain; charset=utf-8")
}
