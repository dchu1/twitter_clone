package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/middleware"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/config"
	"github.com/rs/cors"
	"github.com/spf13/viper"
)

func main() {
	// Read config
	config.NewConfig(".")

	// Register grpc clients
	handlers.RegisterClients()

	mux := http.NewServeMux()
	// mux.HandleFunc("/login", handlers.Login)
	// mux.HandleFunc("/signup", handlers.Signup)
	// mux.HandleFunc("/logout", handlers.Logout)
	// mux.HandleFunc("/feed", handlers.Feed)
	// mux.HandleFunc("/follow/create", handlers.FollowCreateHandler)
	// mux.HandleFunc("follow/destroy", handlers.FollowDestroyHandler)
	// mux.HandleFunc("/user", handlers.UserHandler)
	// mux.HandleFunc("/post", handlers.PostHandler)
	// mux.HandleFunc("/user/following", handlers.UserFollowingHandler)
	// mux.HandleFunc("/user/notfollowing", handlers.UserNotFollowingHandler)

	mux.Handle("/login", middleware.MiddlewareInjector(http.HandlerFunc(handlers.Login), middleware.ContextMiddleware))
	mux.Handle("/signup", middleware.MiddlewareInjector(http.HandlerFunc(handlers.Signup), middleware.ContextMiddleware))
	mux.Handle("/logout", middleware.MiddlewareInjector(http.HandlerFunc(handlers.Logout), middleware.ContextMiddleware))
	mux.Handle("/feed", middleware.MiddlewareInjector(http.HandlerFunc(handlers.Feed), middleware.ContextMiddleware, middleware.AuthMiddleware))
	mux.Handle("/follow/create", middleware.MiddlewareInjector(http.HandlerFunc(handlers.FollowCreateHandler), middleware.ContextMiddleware, middleware.AuthMiddleware))
	mux.Handle("/follow/destroy", middleware.MiddlewareInjector(http.HandlerFunc(handlers.FollowDestroyHandler), middleware.ContextMiddleware, middleware.AuthMiddleware))
	mux.Handle("/user", middleware.MiddlewareInjector(http.HandlerFunc(handlers.UserHandler), middleware.ContextMiddleware, middleware.AuthMiddleware))
	mux.Handle("/post", middleware.MiddlewareInjector(http.HandlerFunc(handlers.PostHandler), middleware.ContextMiddleware, middleware.AuthMiddleware))
	mux.Handle("/user/following", middleware.MiddlewareInjector(http.HandlerFunc(handlers.UserFollowingHandler), middleware.ContextMiddleware, middleware.AuthMiddleware))
	mux.Handle("/user/notfollowing", middleware.MiddlewareInjector(http.HandlerFunc(handlers.UserNotFollowingHandler), middleware.ContextMiddleware, middleware.AuthMiddleware))

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
	fmt.Println("Server running on port", viper.GetStringSlice("webserver.ports")[0])
	err := http.ListenAndServe(":"+viper.GetStringSlice("webserver.ports")[0], handler) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
