package middleware

import (
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := session.GlobalSessions.SessionQuery(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if sess == nil || sess.Get("authenticated") == nil || !sess.Get("authenticated").(bool) {
			http.Error(w, "No Valid Session", http.StatusUnauthorized)
			return
		}
		// create a context and add the session to it
		ctx := session.NewContext(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
