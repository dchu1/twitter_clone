package handlers

import (
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"
)

var a *app.App

func SetState(app *app.App) {
	a = app
}
