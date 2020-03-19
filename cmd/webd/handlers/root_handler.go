package handlers

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/test_storage"
)

type RootHandler struct {
	st                  *test_storage.TestStorage
	userHandler         *UserHandler
	loginHandler        *LoginHandler
	registrationHandler *RegistrationHandler
	postHandler         *PostHandler
	followHandler       *FollowHandler
}

func MakeRootHandler(st *test_storage.TestStorage) *RootHandler {
	return &RootHandler{
		st,
		&UserHandler{st},
		&LoginHandler{st},
		&RegistrationHandler{st},
		&PostHandler{st},
		&FollowHandler{st},
	}
}

func (r *RootHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	fmt.Println("RootHandler received a request")

	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)
	if head == "users" {
		r.userHandler.ServeHTTP(res, req)
		return
	} else if head == "posts" {
		r.postHandler.ServeHTTP(res, req)
		return
	} else if head == "follow" {
		r.followHandler.ServeHTTP(res, req)
		return
	} else if head == "login" {
		r.loginHandler.ServeHTTP(res, req)
		return
	} else if head == "registration" {
		r.registrationHandler.ServeHTTP(res, req)
		return
	} else if head == "" {
		fmt.Fprintf(res, "Landing Page")
		return
	}
	http.Error(res, "404 Not Found", http.StatusNotFound)
}

// ShiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash. Taken from
// https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
