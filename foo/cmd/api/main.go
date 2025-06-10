package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	slog.Info("listening on port 8080")
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from foo-api, %s! You were doing %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)
	}))
}
