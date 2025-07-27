package handler

import (
	"context"
	"fmt"
	"net/http"
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
	sendJSON(w, http.StatusCreated, created)
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
	sendJSON(w, http.StatusOK, l)
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
		handleNotFound(w, err, "resource", fmt.Sprintf("%+v", filter))
		return
	}
	sendJSON(w, http.StatusOK, data)
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
