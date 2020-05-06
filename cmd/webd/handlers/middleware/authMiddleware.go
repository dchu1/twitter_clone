package middleware

import (
	"context"
	"net/http"
	"net/url"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
	"google.golang.org/grpc/status"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionId")
		if err != nil {
			handlers.APIResponse(w, r, http.StatusUnauthorized, "Invalid session", make(map[string]string))
			return
		}
		token, _ := url.QueryUnescape(cookie.Value)
		user, err := handlers.AuthClient.GetUserId(r.Context(), &authpb.AuthToken{Token: token})
		if err != nil {
			e, _ := status.FromError(err)
			if e.Message() == "invalid token" {
				handlers.APIResponse(w, r, http.StatusUnauthorized, "Authentication unsuccesful", make(map[string]string))
			} else {
				handlers.APIResponse(w, r, http.StatusInternalServerError, "Database not responding", make(map[string]string))
			}
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
