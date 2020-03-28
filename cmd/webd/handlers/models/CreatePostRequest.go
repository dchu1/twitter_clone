package handlermodels

type CreatePostRequest struct {
	// UserId  uint64 `json:"userId,omitempty"`
	Message string `json:"message,omitempty"`
}
