package middleware

import (
	"fmt"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
)

//var sess *session.Session

func SetAuthSession(s *session.Session) {
	//sess = s
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.GlobalSessions.SessionStart(w, r)
		if sess == nil || sess.Get("authenticated") == nil || !sess.Get("authenticated").(bool) {
			http.Error(w, "No Valid Session", http.StatusUnauthorized)
			return
		}
		fmt.Printf("I am: %s\n", session.Sess.Get("username"))
		next.ServeHTTP(w, r)
	})
}
