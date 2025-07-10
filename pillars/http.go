package pillars

import (
	"fmt"
	"net/http"
)

func HandlerEcho(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello from %s, %s! You were doing %s %s\n", name, r.RemoteAddr, r.Method, r.RequestURI)
	}
}

func HandlerOk(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	}
}
