package handlers

import (
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	sess := session.GlobalSessions.SessionQuery(w, r)
	if sess != nil {
		sess.Set("authenticated", false)
		session.GlobalSessions.SessionDestroy(w, r)
	}
	APIResponse(w, r, 200, "Logout successful", make(map[string]string))
}