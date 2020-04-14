package middleware

import (
	"context"
	"net/http"
	"time"
)

// ContextMiddleware attaches a timeout to the context of the request
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ctx    context.Context
			cancel context.CancelFunc
		)
		//timeout, err := time.ParseDuration(Cfg.Application.ContextTimeout)
		timeout, err := time.ParseDuration("60s")
		if err == nil {
			// The request has a timeout, so create a context that is
			// canceled automatically when the timeout expires.
			ctx, cancel = context.WithTimeout(r.Context(), timeout)
		} else {
			ctx, cancel = context.WithCancel(r.Context())
		}
		defer cancel() // Cancel ctx as soon as handleSearch returns.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
