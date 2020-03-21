package app

import (
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/cmd/webd/models"
)

func (appList *models.App) FollowUser(followingUserID uint64, UserIDToFollow uint64) {

	//Add userID to be followed in the following list of user who wants to follow
	followingUserIDObject := appList.users[followingUserID]
	following := followingUserIDObject.following
	newfollowing := append(following, UserIDToFollow)
	followingUserIDObject.following = newfollowing

	//Add userID who is following in the followers list of the user being followed
	UserIDToFollowObject := appList.users[UserIDToFollow]
	followers := UserIDToFollowObject.followers
	newfollowers := append(followers, followersUserId)
	UserIDToFollowObject.followers = newfollowers

}

func (appList *models.App) UnFollowUser(followingUserID uint64, UserIDToUnfollow uint64) {

	//Remove userID to be unfollowed from the following list of the user initiating unfollow request
	followingUserIDObject := appList.users[followingUserID]
	following := followingUserIDObject.following
	newfollowing := GetUpdatedList(following, UserIDToUnfollow)
	if newfollowing {
		followingUserIDObject.following = newfollowing
	}

	//Remove userID who is initiating the unfollow request from the followers list of the user being unfollowed
	UserIDToUnfollowObject := appList.users[UserIDToUnfollow]
	followers := UserIDToUnfollowObject.followers
	newfollowers := GetUpdatedList(followers, followingUserID)
	if newfollowers {
		UserIDToUnfollowObject.followers = newfollowers
	}
}

func (*models.User) GetUpdatedList(appList []*models.User, UserID uint64) {
	var updatedList []*models.User
	for userIDIndex := range appList {
		if appList[userIDIndex].id == UserID {
			updatedList = append(appList[:userIDIndex], appList[userIDIndex+1:])
		}

		return updatedList
	}
}

func (appList *models.App) CreatePost(userID uint64, message string) {
	currTime := time.Now()
	newPost = &Post{appList.postID, currTime, message, userID}
	appList.posts[appList.postID] = newPost
	appList.postID++

}

func (appList *models.App) AddUser(userName string) {
	newUser := &User{appList.userId, userName, make([]*User, 10), make([]*models.User, 10), make([]*models.Post, 10)}
	appList.users[appList.userId] = newUser
	appList.userID++
}
