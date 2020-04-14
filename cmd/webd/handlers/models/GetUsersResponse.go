package handlermodels

import userpb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"

type GetUsersResponse struct {
	Users []*userpb.AccountInformation
}
