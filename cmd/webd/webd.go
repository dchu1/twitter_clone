package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"
	session "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
	_ "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/storage"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/config"
	handler "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers"
	handlermodels "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers/models"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/models"
	"github.com/rs/cors"
)

var globalSessions *session.Manager
var sess session.Session
var a *app.App

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
	reqMessage := handlermodels.CreateUserRequest{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(b, &reqMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = a.AddUser(reqMessage.Firstname, reqMessage.Lastname, reqMessage.Email, reqMessage.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		//handler.APIResponse(w, r, 200, "Authorised", make(map[string]string))
	} else {
		handler.APIResponse(w, r, http.StatusUnauthorized, "Unauthorised user", make(map[string]string))
	}
	switch r.Method {
	case "POST":
		reqMessage := handlermodels.FeedRequest{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(b, &reqMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		feed, err := a.GetFeed(reqMessage.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, v := range feed {
			fmt.Println(v.Message)
			w.Write([]byte(fmt.Sprintln(v.Message)))
		}
		//handler.APIResponse(w, r, http.StatusOK, "Added successful;y", make(map[string]string))
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	/*isAuthenticated := sess.Get("authenticated").(bool)
	if isAuthenticated {
		handler.APIResponse(w, r, 200, "Authorised", make(map[string]string))
	} else {
		handler.APIResponse(w, r, http.StatusUnauthorized, "Unauthorised user", make(map[string]string))
	}*/

	switch r.Method {
	case "GET": // FOR TESTING
		u := a.GetUsers()
		for _, v := range u {
			w.Write([]byte(fmt.Sprintf("%s\n", v)))
		}
	case "POST":
		reqMessage := handlermodels.CreateUserRequest{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(b, &reqMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = a.AddUser(reqMessage.Firstname, reqMessage.Lastname, reqMessage.Email, reqMessage.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		handler.APIResponse(w, r, http.StatusOK, "Added successful;y", make(map[string]string))
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := sess.Get("authenticated").(bool)
	if isAuthenticated {
		//handler.APIResponse(w, r, 200, "Authorised", make(map[string]string))
	} else {
		handler.APIResponse(w, r, http.StatusUnauthorized, "Unauthorised user", make(map[string]string))
		return
	}
	switch r.Method {
	case "POST":
		reqMessage := handlermodels.CreatePostRequest{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(b, &reqMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = a.CreatePost(reqMessage.UserId, reqMessage.Message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		handler.APIResponse(w, r, http.StatusOK, "Added successful;y", make(map[string]string))
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}

}

func followCreateHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := sess.Get("authenticated").(bool)
	if isAuthenticated {
		//handler.APIResponse(w, r, 200, "Authorised", make(map[string]string))
	} else {
		handler.APIResponse(w, r, http.StatusUnauthorized, "Unauthorised user", make(map[string]string))
		return
	}
	switch r.Method {
	case "POST":
		reqMessage := handlermodels.FollowRequest{}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(b, &reqMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		a.FollowUser(reqMessage.SourceUserId, reqMessage.TargetUserId)
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
	}
}

func followDestroyHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated := sess.Get("authenticated").(bool)
	if isAuthenticated {
		handler.APIResponse(w, r, 200, "Authorised", make(map[string]string))
	} else {
		handler.APIResponse(w, r, http.StatusUnauthorized, "Unauthorised user", make(map[string]string))
	}
	switch r.Method {
	case "POST":
		//a.UnFollowUser(nil, nil)
	default:
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
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
	mux.HandleFunc("/follow/create", newsfeed)
	mux.HandleFunc("/follow/destroy", newsfeed)
	mux.HandleFunc("/user", userHandler)
	mux.HandleFunc("/post", postHandler)

	a = app.MakeApp()
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
