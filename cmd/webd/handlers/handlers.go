package handlers

import (
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"
)

var application *app.App

func SetState(app *app.App) {
	application = app
}
