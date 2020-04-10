package memstorage

import (
	"context"
	"errors"
	"sync"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
)

type userRepository struct {
	storage *memoryStorage
}

type userEntry struct {
	followingRWMu sync.RWMutex // protects following map
	followersRWMu sync.RWMutex // protects followers map
	postsRWMu     sync.RWMutex // protects posts map
	user          *user.User
}

func NewUserRepository(storage *memoryStorage) user.UserRepository {
	return &userRepository{storage}
}

// CreateUser adds a user to the appropriate data structures
func (userRepository *userRepository) CreateUser(ctx context.Context, info user.AccountInformation) (uint64, error) {
	// Check whether user already exists
	userObj, _ := userRepository.GetUserByUsername(ctx, info.Email)
	if userObj != nil {
		return 0, errors.New("duplicate email")
	}

	info.UserID = userRepository.storage.generateUserId()
	newUserEntry := new(userEntry)
	newUserEntry.user = user.NewUser(info)

	userRepository.storage.usersRWMu.Lock()
	userRepository.storage.users[info.UserID] = newUserEntry
	userRepository.storage.usersRWMu.Unlock()

	return info.UserID, nil
}

// GetUser creates a deep copy of the specified users.
func (userRepo *userRepository) GetUser(ctx context.Context, userID uint64) (*user.User, error) {
	userRepo.storage.usersRWMu.RLock()
	defer userRepo.storage.usersRWMu.RUnlock()
	userEntry, exists := userRepo.storage.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return userEntry.user.Clone(), nil
}

// GetUsers creates a deep copy of the specified users.
func (userRepo *userRepository) GetUsers(ctx context.Context, userIDs []uint64) ([]*user.User, error) {
	userRepo.storage.usersRWMu.RLock()
	defer userRepo.storage.usersRWMu.RUnlock()
	cp := make([]*user.User, 0, len(userIDs))
	for _, v := range userIDs {
		cp = append(cp, userRepo.storage.users[v].user.Clone())
	}
	return cp, nil
}

// FollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user is following another user
func (userRepo *userRepository) FollowUser(ctx context.Context, followingUserID uint64, UserIDToFollow uint64) error {
	if followingUserID == UserIDToFollow {
		return errors.New("duplicate user ids")
	}

	//Add userID to be followed in the following list of user who wants to follow
	followingUserIDObject, err := userRepo.storage.getUserEntry(followingUserID)
	if err != nil {
		return err
	}
	followingUserIDObject.followingRWMu.Lock()
	newfollowing := followingUserIDObject.user.Following
	newfollowing[UserIDToFollow] = struct{}{}
	followingUserIDObject.user.Following = newfollowing
	followingUserIDObject.followingRWMu.Unlock()

	//Add userID who is following in the followers list of the user being followed
	UserIDToFollowObject, err := userRepo.storage.getUserEntry(UserIDToFollow)
	if err != nil {
		return err
	}
	UserIDToFollowObject.followersRWMu.Lock()
	newfollowers := UserIDToFollowObject.user.Followers
	newfollowers[followingUserID] = struct{}{}
	UserIDToFollowObject.user.Followers = newfollowers
	UserIDToFollowObject.followersRWMu.Unlock()

	return nil
}

// UnFollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user has unfollowed another user
func (userRepo *userRepository) UnFollowUser(ctx context.Context, followingUserID uint64, UserIDToUnfollow uint64) error {
	if followingUserID == UserIDToUnfollow {
		return errors.New("duplicate user ids")
	}

	//Remove userID to be unfollowed from the following list of the user initiating unfollow request
	followingUserIDObject, err := userRepo.storage.getUserEntry(followingUserID)
	if err != nil {
		return err
	}
	followingUserIDObject.followingRWMu.Lock()
	newfollowing := followingUserIDObject.user.Following
	delete(newfollowing, UserIDToUnfollow)
	followingUserIDObject.user.Following = newfollowing
	followingUserIDObject.followingRWMu.Unlock()

	//Remove userID who is initiating the unfollow request from the followers list of the user being unfollowed
	UserIDToUnfollowObject, err := userRepo.storage.getUserEntry(UserIDToUnfollow)
	if err != nil {
		return err
	}
	UserIDToUnfollowObject.followersRWMu.Lock()
	newfollowers := UserIDToUnfollowObject.user.Followers
	delete(newfollowers, followingUserID)
	UserIDToUnfollowObject.user.Followers = newfollowers
	UserIDToUnfollowObject.followersRWMu.Unlock()

	return nil
}

// // ValidateCredentials checks whether the given username and password match
// // those stored in the credentials map
// func (appList *App) ValidateCredentials(username string, password string) bool {
// 	appList.credentialsRWMu.RLock()
// 	defer appList.credentialsRWMu.RUnlock()
// 	return appList.credentials[username] == password
// }

// GetUserByUsername returns a user object by their username
func (userRepo *userRepository) GetUserByUsername(ctx context.Context, email string) (*user.User, error) {
	userRepo.storage.usersRWMu.RLock()
	defer userRepo.storage.usersRWMu.RUnlock()
	for _, v := range userRepo.storage.users {
		if v.user.AccountInformation.Email == email {
			return v.user.Clone(), nil
		}
	}
	return nil, errors.New("user not found")
}

// GetFollowing returns an array of users that the given user is following
func (userRepo *userRepository) GetFollowing(ctx context.Context, userId uint64) ([]*user.User, error) {
	// Get the user object from the users map
	userEntry, err := userRepo.storage.getUserEntry(userId)
	if err != nil {
		return nil, err
	}
	userEntry.followingRWMu.RLock()
	defer userEntry.followingRWMu.RUnlock()

	tempArray := make([]*user.User, 0, 100)
	for k := range userEntry.user.Following {
		followingEntry, err := userRepo.storage.getUserEntry(k)
		if err != nil {
			// if we have an error here, it means our following data structure has an entry inconsistent
			// with our user structure
			panic("database corruption")
		}
		tempArray = append(tempArray, followingEntry.user.Clone())
	}
	return tempArray, nil
}

// GetNotFollowing returns an array of users that the given user is not following
func (userRepo *userRepository) GetNotFollowing(ctx context.Context, userId uint64) ([]*user.User, error) {
	// Get the user object from the users map
	userEntry, err := userRepo.storage.getUserEntry(userId)
	if err != nil {
		return nil, err
	}
	userEntry.followingRWMu.RLock()
	defer userEntry.followingRWMu.RUnlock()

	tempArray := make([]*user.User, 0, 100)

	// Iterate through entire user list
	userRepo.storage.usersRWMu.RLock()
	defer userRepo.storage.usersRWMu.RUnlock()
	for k, v := range userRepo.storage.users {
		// check if user k exists in the user's following list. If not, add it to our
		// temp array
		_, exists := userEntry.user.Following[k]
		if !exists && k != userId {
			tempArray = append(tempArray, v.user.Clone())
		}
	}
	return tempArray, nil
}

func (userRepo *userRepository) AddPost(ctx context.Context, userID uint64, postID uint64) error {
	userEntry, err := userRepo.storage.getUserEntry(userID)
	if err != nil {
		return err
	}
	userEntry.postsRWMu.Lock()
	defer userEntry.postsRWMu.Unlock()
	userEntry.user.Posts = append(userEntry.user.Posts, postID)
	return nil
}

func (userRepo *userRepository) DeleteUser(ctx context.Context, userID uint64) error {
	return errors.New("Feature not implemented")
}
func (userRepo *userRepository) UpdateUserAccountInfo(ctx context.Context, info user.AccountInformation) error {
	return errors.New("Feature not implemented")
}
