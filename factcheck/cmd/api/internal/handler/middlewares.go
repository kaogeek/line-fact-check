package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/kaogeek/line-fact-check/factcheck"
)

type (
	CtxKey string
)

const (
	CtxKeyUserID   CtxKey = "FACTCHECK_USERID"
	CtxKeyUserType CtxKey = "FACTCHECK_USERTYPE"
	CtxKeyUserInfo CtxKey = "FACTCHECK_USERINFO"
)

func decodeAdmin(ctx context.Context, userID string) context.Context {
	ctx = context.WithValue(ctx, CtxKeyUserID, userID)
	ctx = context.WithValue(ctx, CtxKeyUserType, factcheck.TypeUserMessageAdmin)
	ctx = context.WithValue(ctx, CtxKeyUserInfo, factcheck.UserInfo{
		UserType: factcheck.TypeUserMessageAdmin,
		UserID:   userID,
	})
	return ctx
}

// MiddlewareAuth handles only authentication
func MiddlewareAuth(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(decodeAdmin(r.Context(), "mock-factcheck-admin"))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func MiddlewareAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType := r.Context().Value(CtxKeyUserType)
		if userType == nil {
			w.WriteHeader(http.StatusUnauthorized)
			write(r.Context(), w, []byte("unauthorized: missing data"))
			return
		}
		userType, ok := userType.(factcheck.TypeUser)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			write(r.Context(), w, []byte("unauthorized: missing data"))
			return
		}
		if userType != factcheck.TypeUserMessageAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			write(r.Context(), w, []byte("unauthorized: bad data"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func write(ctx context.Context, w http.ResponseWriter, b []byte) {
	_, err := w.Write(b)
	if err != nil {
		slog.ErrorContext(ctx, "error writing response",
			"err", err,
			"bytes", string(b),
		)
		return
	}
}
