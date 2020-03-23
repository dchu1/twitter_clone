package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
	_ "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/storage"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/config"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/middleware"
	"github.com/rs/cors"
)

func main() {
	// Read config
	cfg := config.GetConfig(".")
	session.InitializeGlobalSessions()
	handlers.SetState(app.MakeApp())

	mux := http.NewServeMux()
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/signup", handlers.Signup)
	mux.HandleFunc("/logout", handlers.Logout)
	mux.Handle("/feed", middleware.MiddlewareInjector(http.HandlerFunc(handlers.Feed), middleware.AuthMiddleware))
	mux.Handle("/follow/create", middleware.MiddlewareInjector(http.HandlerFunc(handlers.FollowCreateHandler), middleware.AuthMiddleware))
	mux.Handle("/follow/destroy", middleware.MiddlewareInjector(http.HandlerFunc(handlers.FollowDestroyHandler), middleware.AuthMiddleware))
	mux.Handle("/user", middleware.MiddlewareInjector(http.HandlerFunc(handlers.UserHandler), middleware.AuthMiddleware))
	mux.Handle("/post", middleware.MiddlewareInjector(http.HandlerFunc(handlers.PostHandler), middleware.AuthMiddleware))

	origins := []string{"http://localhost:4200"}
	headers := []string{"Content-Type", "X-Requested-With", "Range"}
	exposeHeader := []string{"Accept-Ranges", "Content-Encoding", "Content-Length", "Content-Range", "Set-Cookie"}
	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowCredentials: true,
		AllowedHeaders:   headers,
		ExposedHeaders:   exposeHeader,
	})

	handler := cors.Default().Handler(mux)
	handler = c.Handler(handler)
	fmt.Println("Server running on port", cfg.Server.Port)
	err := http.ListenAndServe(":"+cfg.Server.Port, handler) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
