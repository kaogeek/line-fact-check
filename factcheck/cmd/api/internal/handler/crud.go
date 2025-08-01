package handler

import (
	"context"
	"fmt"
	"net/http"
)

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
	sendJSON(r.Context(), w, http.StatusOK, l)
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
	sendJSON(r.Context(), w, http.StatusOK, data)
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
	sendText(r.Context(), w, "ok", http.StatusOK)
}
