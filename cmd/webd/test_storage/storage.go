package test_storage

import (
	"fmt"
	"strings"
	"time"
)

type TestStorage struct {
	users  map[uint64]*User
	posts  map[uint64]*Post
	userId uint64
	postId uint64
}

type Post struct {
	id        uint64    // This is a unique id. Type might be different depending on how we generate unique ids.
	timestamp time.Time // time this post was made
	message   string    // the text of the post
	author    []*User
}

type User struct {
	id        uint64  // This is a unique id. Type might be different depending on how we generate unique ids.
	name      string  // The user's name
	following []*User // list of ids this user is following
	followers []*User // list of ids following this user
	posts     []*Post // list of posts by this user
}

func (a *TestStorage) AddUser(name string) {
	newUser := &User{a.userId, name, make([]*User, 10), make([]*User, 10), make([]*Post, 10)}
	a.users[a.userId] = newUser
	a.userId++
	fmt.Printf("Added user %s. Users: %v\n", name, a.users)
}

func Make() *TestStorage {
	return &TestStorage{make(map[uint64]*User), make(map[uint64]*Post), 0, 0}
}

func (a *TestStorage) PrintUsers() string {
	var b strings.Builder
	for key, value := range a.users {
		b.WriteString(fmt.Sprintf("%s[%d], ", value.name, key))
	}
	return b.String()

}
