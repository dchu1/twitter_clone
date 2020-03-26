package handlermodels

import "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/app"

type GetUserFollowingResponse struct {
	Users []app.User
}
