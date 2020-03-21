package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	session "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
	_ "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/storage"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/config"
	handler "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/models"
	"github.com/rs/cors"
)

var globalSessions *session.Manager
var sess session.Session

func login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println("Hi")
	sess = globalSessions.SessionStart(w, r)

	json.Unmarshal([]byte(body), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//code to check if user exists
	fmt.Println("Hello")
	sess.Set("username", user.Email)
	sess.Set("authenticated", true)
	handler.APIResponse(w, r, 200, "Login successful", make(map[string]string)) // send data to client side
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	body, err := ioutil.ReadAll(r.Body)

	json.Unmarshal([]byte(body), &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	handler.APIResponse(w, r, http.StatusCreated, "Signup successful", make(map[string]string))
}

func logout(w http.ResponseWriter, r *http.Request) {
	sess.Set("authenticated", false)
	globalSessions.SessionDestroy(w, r)
	handler.APIResponse(w, r, 200, "Logout successful", make(map[string]string))
}

func newsfeed(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := sess.Get("authenticated").(bool)
	if isAuthenticated {
		handler.APIResponse(w, r, 200, "Authorised", make(map[string]string))
	} else {
		handler.APIResponse(w, r, http.StatusUnauthorized, "Unauthorised user", make(map[string]string))
	}

}

// Then, initialize the session manager
func init() {
	globalSessions, _ = session.NewManager("memory", "client_sessionid", 3600)
}

func main() {
	// Read config
	cfg := config.GetConfig(".")
	// storage := test_storage.Make()
	// rootHandler := handlers.MakeRootHandler(storage)
	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/signup", signup)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/home", newsfeed)

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
