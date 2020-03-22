package handlermodels

type FollowRequest struct {
	SourceUserId uint64 `json:"sourceUserId,omitempty"`
	TargetUserId uint64 `json:"targetUserId,omitempty"`
}
