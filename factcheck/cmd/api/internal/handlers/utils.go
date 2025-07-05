package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func body(r *http.Request) ([]byte, error) {
	var b []byte
	_, err := io.ReadFull(r.Body, b)
	if err != nil {
		return nil, fmt.Errorf("failed to copy body request into buffer: %w", err)
	}
	return b, nil
}

// replyJson marshals data into JSON string before writing response.
// If marshaling failed, the response is left untouched and the error is returned.
func replyJson(w http.ResponseWriter, data any, status int) error {
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
