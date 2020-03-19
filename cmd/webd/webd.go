package main

import (
	"log"
	"net/http"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/config"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/handlers"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/test_storage"
)

func main() {
	// Read config
	cfg := config.GetConfig("./cmd/webd/")
	storage := test_storage.Make()
	rootHandler := handlers.MakeRootHandler(storage)
	err := http.ListenAndServe(":"+cfg.Server.Port, rootHandler) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
