package handlers

import (
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/auth/session"
)

var sessionManager *session.Manager
var a *app.App
var sess *session.Session

func SetState(app *app.App) {
	a = app
}
