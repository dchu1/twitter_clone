package handlermodels

import "time"

type FeedResponse struct {
	Posts []Post `json:"posts,omitempty"`
}

type Post struct {
	Id        uint64    `json:"postId"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message,omitempty"`
	Author    Author    `json:"author,omitempty"`
}

type Author struct {
	UserID    uint64 `json:"userId"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
}
