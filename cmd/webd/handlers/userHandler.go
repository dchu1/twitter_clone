package handlers

import (
	"fmt"
	"net/http"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // FOR TESTING
		u := a.GetUsers()
		for _, v := range u {
			w.Write([]byte(fmt.Sprintf("%s\n", v)))
		}
	default:
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
	}
}
