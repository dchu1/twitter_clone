package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers"
	authpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/auth/authentication"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionId")
		if err != nil {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
		token, _ := url.QueryUnescape(cookie.Value)
		user, err := handlers.AuthClient.GetUserId(r.Context(), &authpb.AuthToken{Token: token})
		// }
		ctx := context.WithValue(r.Context(), "user", user)
		fmt.Println(user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
