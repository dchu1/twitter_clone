package middleware

import (
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
)

var sess *session.Session

func SetAuthSession(s *session.Session) {
	sess = s
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if *sess == nil || !(*sess).Get("authenticated").(bool) {
			http.Error(w, "No Valid Session", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
