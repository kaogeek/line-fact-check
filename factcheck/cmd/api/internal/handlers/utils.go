package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func create[T any](
	w http.ResponseWriter,
	r *http.Request,
	repo interface {
		Create(context.Context, T) (T, error)
	},
) {
	data, err := decode[T](r)
	if err != nil {
		errBadRequest(w, err.Error())
		return
	}
	created, err := repo.Create(r.Context(), data)
	if err != nil {
		errInternalError(w, err.Error())
		return
	}
	replyJson(w, created, http.StatusCreated)
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
	replyJson(w, l, http.StatusOK)
}

func decode[T any](r *http.Request) (T, error) {
	var zero T
	b, err := body(r)
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

func body(r *http.Request) ([]byte, error) {
	var b []byte
	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		return nil, fmt.Errorf("failed to copy body request into buffer: %w", err)
	}
	return b, nil
}

// replyJson calls replyJsonError, and on non-nil error, writes 500 response
func replyJson(w http.ResponseWriter, data any, status int) {
	err := replyJsonError(w, data, status)
	if err != nil {
		errInternalError(w, err.Error())
	}
}

// replyJsonError marshals data into JSON string before writing response.
// If marshaling failed, the response is left untouched and the error is returned.
func replyJsonError(w http.ResponseWriter, data any, status int) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	contentTypeJSON(w.Header())
	w.WriteHeader(status)
	w.Write(j)
	return nil
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
