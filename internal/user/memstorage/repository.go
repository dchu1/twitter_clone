package memstorage

import (
	"context"
	"errors"
	"sync"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
)

type userRepository struct {
	storage *userStorage
}

type userEntry struct {
	followingRWMu sync.RWMutex // protects following map
	followersRWMu sync.RWMutex // protects followers map
	user          *pb.User
}

func GetUserRepository() user.UserRepository {
	return &userRepository{UserStorage}
}

func NewUserRepository(storage *userStorage) user.UserRepository {
	return &userRepository{storage}
}

// CreateUser adds a user to the appropriate data structures
func (userRepo *userRepository) CreateUser(ctx context.Context, info *userpb.AccountInformation) (uint64, error) {
	result := make(chan uint64)
	errorchan := make(chan error)

	go func() {

		info.UserId = userRepo.storage.generateUserId()
		newUserEntry := new(userEntry)
		newUserEntry.user = new(pb.User)
		newUserEntry.user.AccountInformation = info
		newUserEntry.user.Followers = make(map[uint64]uint64)
		newUserEntry.user.Following = make(map[uint64]uint64)

		userRepo.storage.usersRWMu.Lock()
		userRepo.storage.users[info.UserId] = newUserEntry
		userRepo.storage.usersRWMu.Unlock()
		result <- info.UserId
		errorchan <- nil

	}()

	select {
	case userID := <-result:
		return userID, nil
	case err := <-errorchan:
		//Sending 0 as an invalid postID
		return 0, err
	case <-ctx.Done():
		delete(userRepo.storage.users, info.UserId)
		return 0, ctx.Err()
	}
}

// GetUser creates a deep copy of the specified users.
func (userRepo *userRepository) GetUser(ctx context.Context, userID uint64) (*pb.User, error) {
	result := make(chan *pb.User)
	errorchan := make(chan error)

	go func() {

		userRepo.storage.usersRWMu.RLock()
		defer userRepo.storage.usersRWMu.RUnlock()
		userEntry, exists := userRepo.storage.users[userID]
		if !exists {
			result <- nil
			errorchan <- errors.New("user not found")
		} else {
			result <- userEntry.user
			errorchan <- nil
		}
	}()

	select {
	case user := <-result:
		return user, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetUsers creates a deep copy of the specified users.
func (userRepo *userRepository) GetUsers(ctx context.Context, userIDs []uint64) ([]*pb.User, error) {
	result := make(chan []*pb.User)
	errorchan := make(chan error)

	go func() {
		userRepo.storage.usersRWMu.RLock()
		defer userRepo.storage.usersRWMu.RUnlock()
		cp := make([]*pb.User, 0, len(userIDs))
		for _, v := range userIDs {
			cp = append(cp, userRepo.storage.users[v].user)
		}
		result <- cp
		errorchan <- nil
	}()

	select {
	case users := <-result:
		return users, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (userRepo *userRepository) GetAllUsers(ctx context.Context) ([]*pb.User, error) {
	result := make(chan []*pb.User)
	errorchan := make(chan error)

	go func() {
		userRepo.storage.usersRWMu.RLock()
		defer userRepo.storage.usersRWMu.RUnlock()
		tempArr := make([]*pb.User, 0, len(userRepo.storage.users))
		for _, u := range userRepo.storage.users {
			tempArr = append(tempArr, u.user)
		}
		result <- tempArr
		errorchan <- nil
	}()

	select {
	case users := <-result:
		return users, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// FollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user is following another user
func (userRepo *userRepository) FollowUser(ctx context.Context, followingUserID uint64, UserIDToFollow uint64) error {

	result := make(chan error)
	go func() {
		if followingUserID == UserIDToFollow {
			result <- errors.New("duplicate user ids")
		} else {

			//Add userID to be followed in the following list of user who wants to follow
			followingUserIDObject, err := userRepo.storage.getUserEntry(followingUserID)
			if err != nil {
				result <- err
			} else {
				followingUserIDObject.followingRWMu.Lock()
				followingUserIDObject.user.Following[UserIDToFollow] = UserIDToFollow
				followingUserIDObject.followingRWMu.Unlock()

				//Add userID who is following in the followers list of the user being followed
				UserIDToFollowObject, err := userRepo.storage.getUserEntry(UserIDToFollow)
				if err != nil {
					result <- err
				} else {
					UserIDToFollowObject.followersRWMu.Lock()
					UserIDToFollowObject.user.Followers[followingUserID] = followingUserID
					UserIDToFollowObject.followersRWMu.Unlock()

					result <- nil
				}
			}
		}
	}()

	select {
	case res := <-result:
		return res
	case <-ctx.Done():
		return ctx.Err()
	}
}

// UnFollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user has unfollowed another user
func (userRepo *userRepository) UnFollowUser(ctx context.Context, followingUserID uint64, UserIDToUnfollow uint64) error {

	result := make(chan error)

	go func() {
		if followingUserID == UserIDToUnfollow {
			result <- errors.New("duplicate user ids")
		} else {

			//Remove userID to be unfollowed from the following list of the user initiating unfollow request
			followingUserIDObject, err := userRepo.storage.getUserEntry(followingUserID)
			if err != nil {
				result <- err
			} else {
				followingUserIDObject.followingRWMu.Lock()
				newfollowing := followingUserIDObject.user.Following
				delete(newfollowing, UserIDToUnfollow)
				followingUserIDObject.user.Following = newfollowing
				followingUserIDObject.followingRWMu.Unlock()

				//Remove userID who is initiating the unfollow request from the followers list of the user being unfollowed
				UserIDToUnfollowObject, err := userRepo.storage.getUserEntry(UserIDToUnfollow)
				if err != nil {
					result <- err
				} else {
					UserIDToUnfollowObject.followersRWMu.Lock()
					newfollowers := UserIDToUnfollowObject.user.Followers
					delete(newfollowers, followingUserID)
					UserIDToUnfollowObject.user.Followers = newfollowers
					UserIDToUnfollowObject.followersRWMu.Unlock()

					result <- nil
				}
			}
		}
	}()

	select{
	case res := <-result:
		return res
	case <-ctx.Done():
		return ctx.Err()
}

// GetUserByUsername returns a user object by their username
func (userRepo *userRepository) GetUserByUsername(ctx context.Context, email string) (*pb.User, error) {
	result := make(chan *pb.User)
	errorchan := make(chan error)

	go func() {
		userRepo.storage.usersRWMu.RLock()
		defer userRepo.storage.usersRWMu.RUnlock()
		exists := false

		for _, v := range userRepo.storage.users {
			if v.user.AccountInformation.Email == email {
				result <- v.user
				errorchan <- nil
				exists = true
			}
		}
		if !exists {
			result <- nil
			errorchan <- errors.New("user not found")
		}
	}()

	select {
	case user := <-result:
		return user, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetFollowing returns an array of users that the given user is following
func (userRepo *userRepository) GetFollowing(ctx context.Context, userId uint64) ([]*pb.User, error) {
	result := make(chan []*pb.User)
	errorchan := make(chan error)

	go func() {
		// Get the user object from the users map
		userEntry, err := userRepo.storage.getUserEntry(userId)
		if err != nil {
			result <- nil
			errorchan <- err
		} else {
			userEntry.followingRWMu.RLock()
			defer userEntry.followingRWMu.RUnlock()
			databaseError := false
			tempArray := make([]*pb.User, 0, 100)
			for k := range userEntry.user.Following {
				followingEntry, err := userRepo.storage.getUserEntry(k)
				if err != nil {
					// if we have an error here, it means our following data structure has an entry inconsistent
					// with our user structure
					databaseError = true
					result <- nil
					errorchan <- errors.New("database corruption")
					panic("database corruption")
				}
				tempArray = append(tempArray, followingEntry.user)
			}
			if !databaseError {
				result <- tempArray
				errorchan <- nil
			}
		}
	}()

	select {
	case user := <-result:
		return user, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

}

// GetNotFollowing returns an array of users that the given user is not following
func (userRepo *userRepository) GetNotFollowing(ctx context.Context, userId uint64) ([]*pb.User, error) {

	result := make(chan []*pb.User)
	errorchan := make(chan error)

	go func() {
		// Get the user object from the users map
		userEntry, err := userRepo.storage.getUserEntry(userId)
		if err != nil {
			result <- nil
			errorchan <- err
		} else {
			userEntry.followingRWMu.RLock()
			defer userEntry.followingRWMu.RUnlock()

			tempArray := make([]*pb.User, 0, 100)

			// Iterate through entire user list
			userRepo.storage.usersRWMu.RLock()
			defer userRepo.storage.usersRWMu.RUnlock()
			for k, v := range userRepo.storage.users {
				// check if user k exists in the user's following list. If not, add it to our
				// temp array
				_, exists := userEntry.user.Following[k]
				if !exists && k != userId {
					tempArray = append(tempArray, v.user)
				}
			}

			result <- tempArray
			errorchan <- nil
		}
	}()

	select {
	case user := <-result:
		return user, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}

}

func (userRepo *userRepository) DeleteUser(ctx context.Context, userID uint64) error {
	return errors.New("Feature not implemented")
}
func (userRepo *userRepository) UpdateUserAccountInfo(ctx context.Context, info *userpb.AccountInformation) error {
	return errors.New("Feature not implemented")
}
