package handlers

import (
	"context"
	"net/http"
	"net/url"
	"time"

	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
)

// Logout is the handler for /logout
func Logout(w http.ResponseWriter, r *http.Request) {
	//Get user id of the session
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cookie, err := r.Cookie("sessionId")
	if err != nil || cookie.Value != "" {
		token, _ := url.QueryUnescape(cookie.Value)
		_, err = AuthClient.RemoveAuthToken(ctx, &authpb.AuthToken{Token: token})
		expiration := time.Now()
		cookie := http.Cookie{Name: "sessionId", Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
	APIResponse(w, r, 200, "Logout successful", make(map[string]string))
}
