package etcd

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"log"
	"strconv"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

const userPrefix string = "User/"
const followingPrefix string = "Following"
const userIdGenKey string = "NextUserId"

type userRepository struct {
	storage *clientv3.Client
}

// GetUserRepository returns a UserRepository that uses package level storage
func GetUserRepository() user.UserRepository {
	return &userRepository{Client}
}

// NewUserRepository reutnrs a UserRepository that uses the given storage
func NewUserRepository(storage *clientv3.Client) user.UserRepository {
	return &userRepository{storage}
}

// CreateUser adds a user to the appropriate data structures
func (userRepo *userRepository) CreateUser(ctx context.Context, info *userpb.AccountInformation) (uint64, error) {
	result := make(chan uint64, 1)
	errorchan := make(chan error, 1)

	// Ignore commented out code
	// go func() {
	// 	id, err := userRepo.getUserId(ctx)
	// 	if err != nil {
	// 		errorchan <- err
	// 		return
	// 	}
	// 	info.UserId = id
	// 	// user.AccountInformation = info
	// 	// user.Followers = make(map[uint64]uint64)
	// 	// user.Following = make(map[uint64]uint64)

	// 	var buf bytes.Buffer
	// 	var buf2 bytes.Buffer
	// 	if err := gob.NewEncoder(&buf).Encode(*info); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if err := gob.NewEncoder(&buf2).Encode(*info); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	createUser := func(stm concurrency.STM) error {
	// 		if _, err := userRepo.storage.Put(ctx, userPrefix+strconv.FormatUint(info.UserId, 10), buf.String()); err != nil {
	// 			return err
	// 		}
	// 		if _, err = userRepo.storage.Put(ctx, userPrefix+strconv.FormatUint(info.UserId, 10)+followingPrefix, buf2.String()); err != nil {
	// 			// may need to also delete? or is that handled by etcd?
	// 			return err
	// 		}
	// 		return nil
	// 	}
	// 	if _, err := concurrency.NewSTM(userRepo.storage, createUser); err != nil {
	// 		errorchan <- err
	// 		return
	// 	}
	// 	result <- info.UserId
	// }()

	go func() {
		//info.UserId = 0
		user := new(pb.User)
		info.UserId, _ = userRepo.getUserId(ctx)
		user.AccountInformation = info
		user.Followers = make(map[uint64]uint64)
		user.Following = make(map[uint64]uint64)

		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(user); err != nil {
			log.Fatal(err)
		}

		_, err := userRepo.storage.Put(ctx, userPrefix+strconv.FormatUint(info.UserId, 10), buf.String())
		if err != nil {
			log.Fatal(err)
		}
		result <- info.UserId
	}()

	select {
	case userID := <-result:
		return userID, nil
	case err := <-errorchan:
		//Sending 0 as an invalid postID
		return 0, err
	case <-ctx.Done():
		// if ctx.Done(), we need to make sure that if the user has or will be created, it is deleted,
		// so start a new go routine to monitor the result and error channels
		go func() {
			select {
			case userID := <-result:
				userRepo.DeleteUser(context.Background(), userID)
				return
			case <-errorchan:
				return
			}
		}()
		return 0, ctx.Err()
	}
}

// GetUser creates a copy of the specified user.
func (userRepo *userRepository) GetUser(ctx context.Context, userID uint64) (*pb.User, error) {
	result := make(chan *pb.User, 1)
	errorchan := make(chan error, 1)

	go func() {
		var user pb.User
		resp, err := userRepo.storage.Get(ctx, userPrefix+strconv.FormatUint(userID, 10))
		if err != nil {
			log.Fatal(err)
		}
		dec := gob.NewDecoder(bytes.NewReader(resp.Kvs[0].Value))
		if err := dec.Decode(&user); err != nil {
			log.Fatalf("could not decode message (%v)", err)
		}
		result <- &user
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

// GetUsers creates a copy of the specified users.
func (userRepo *userRepository) GetUsers(ctx context.Context, userIDs []uint64) ([]*pb.User, error) {
	result := make(chan []*pb.User, 1)
	errorchan := make(chan error, 1)

	go func() {
		extent, err := findRange(userIDs)
		if err != nil {
			errorchan <- err
			return
		}
		resp, err := userRepo.storage.Get(ctx, userPrefix+strconv.FormatUint(extent[0], 10), clientv3.WithRange(userPrefix+strconv.FormatUint(extent[1]+1, 10)))
		cp := make([]*pb.User, 0, len(userIDs))
		for _, v := range resp.Kvs {
			var user pb.User
			dec := gob.NewDecoder(bytes.NewReader(v.Value))
			if err := dec.Decode(&user); err != nil {
				log.Fatalf("could not decode message (%v)", err)
			}
			cp = append(cp, &user)
		}
		result <- cp
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

// GetAllUsers returns all users
func (userRepo *userRepository) GetAllUsers(ctx context.Context) ([]*pb.User, error) {
	result := make(chan []*pb.User, 1)
	errorchan := make(chan error, 1)

	go func() {
		resp, err := userRepo.storage.Get(ctx, "User/", clientv3.WithPrefix())
		if err != nil {
			log.Fatal(err)
		}
		cp := make([]*pb.User, 0, 100)
		for _, v := range resp.Kvs {
			var user pb.User
			dec := gob.NewDecoder(bytes.NewReader(v.Value))
			if err := dec.Decode(&user); err != nil {
				log.Fatalf("could not decode message (%v)", err)
			}
			cp = append(cp, &user)
		}
		result <- cp
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

	// result := make(chan error, 1)
	// go func() {
	// 	if followingUserID == UserIDToFollow {
	// 		result <- errors.New("duplicate user ids")
	// 	} else {

	// 		//Add userID to be followed in the following list of user who wants to follow
	// 		followingUserIDObject, err := userRepo.storage.getUserEntry(followingUserID)
	// 		if err != nil {
	// 			result <- err
	// 		} else {
	// 			followingUserIDObject.followingRWMu.Lock()
	// 			followingUserIDObject.user.Following[UserIDToFollow] = UserIDToFollow
	// 			followingUserIDObject.followingRWMu.Unlock()

	// 			//Add userID who is following in the followers list of the user being followed
	// 			UserIDToFollowObject, err := userRepo.storage.getUserEntry(UserIDToFollow)
	// 			if err != nil {
	// 				result <- err
	// 			} else {
	// 				UserIDToFollowObject.followersRWMu.Lock()
	// 				UserIDToFollowObject.user.Followers[followingUserID] = followingUserID
	// 				UserIDToFollowObject.followersRWMu.Unlock()

	// 				result <- nil
	// 			}
	// 		}
	// 	}
	// }()

	// select {
	// case res := <-result:
	// 	return res
	// case <-ctx.Done():
	// 	// listen to the result channel in case the operation was successful, then unfollow
	// 	go func() {
	// 		res := <-result
	// 		if res == nil {
	// 			userRepo.UnFollowUser(context.Background(), followingUserID, UserIDToFollow)
	// 		}
	// 	}()
	// 	return ctx.Err()
	// }
	return nil
}

// UnFollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user has unfollowed another user
func (userRepo *userRepository) UnFollowUser(ctx context.Context, followingUserID uint64, UserIDToUnfollow uint64) error {

	// result := make(chan error, 1)

	// go func() {
	// 	if followingUserID == UserIDToUnfollow {
	// 		result <- errors.New("duplicate user ids")
	// 	} else {

	// 		//Remove userID to be unfollowed from the following list of the user initiating unfollow request
	// 		followingUserIDObject, err := userRepo.storage.getUserEntry(followingUserID)
	// 		if err != nil {
	// 			result <- err
	// 		} else {
	// 			followingUserIDObject.followingRWMu.Lock()
	// 			newfollowing := followingUserIDObject.user.Following
	// 			delete(newfollowing, UserIDToUnfollow)
	// 			followingUserIDObject.user.Following = newfollowing
	// 			followingUserIDObject.followingRWMu.Unlock()

	// 			//Remove userID who is initiating the unfollow request from the followers list of the user being unfollowed
	// 			UserIDToUnfollowObject, err := userRepo.storage.getUserEntry(UserIDToUnfollow)
	// 			if err != nil {
	// 				result <- err
	// 			} else {
	// 				UserIDToUnfollowObject.followersRWMu.Lock()
	// 				newfollowers := UserIDToUnfollowObject.user.Followers
	// 				delete(newfollowers, followingUserID)
	// 				UserIDToUnfollowObject.user.Followers = newfollowers
	// 				UserIDToUnfollowObject.followersRWMu.Unlock()

	// 				result <- nil
	// 			}
	// 		}
	// 	}
	// }()

	// select {
	// case res := <-result:
	// 	return res
	// case <-ctx.Done():
	// 	// listen to the result channel in case the operation was successful, then follow
	// 	go func() {
	// 		res := <-result
	// 		if res == nil {
	// 			userRepo.FollowUser(context.Background(), followingUserID, UserIDToUnfollow)
	// 		}
	// 	}()
	// 	return ctx.Err()
	// }
	return nil
}

// GetUserByUsername returns a user object by their username
func (userRepo *userRepository) GetUserByUsername(ctx context.Context, email string) (*pb.User, error) {
	// result := make(chan *pb.User, 1)
	// errorchan := make(chan error, 1)

	// go func() {
	// 	userRepo.storage.usersRWMu.RLock()

	// 	exists := false

	// 	for _, v := range userRepo.storage.users {
	// 		if v.user.AccountInformation.Email == email {
	// 			result <- v.user
	// 			exists = true
	// 		}
	// 	}
	// 	if !exists {
	// 		errorchan <- errors.New("user not found")
	// 	}
	// 	userRepo.storage.usersRWMu.RUnlock()
	// }()

	// select {
	// case user := <-result:
	// 	return user, nil
	// case err := <-errorchan:
	// 	return nil, err
	// case <-ctx.Done():
	// 	return nil, ctx.Err()
	// }
	return nil, nil
}

// GetFollowing returns an array of users that the given user is following
func (userRepo *userRepository) GetFollowing(ctx context.Context, userId uint64) ([]*pb.User, error) {
	// result := make(chan []*pb.User, 1)
	// errorchan := make(chan error, 1)

	// go func() {
	// 	// Get the user object from the users map
	// 	userEntry, err := userRepo.storage.getUserEntry(userId)
	// 	if err != nil {
	// 		errorchan <- err
	// 	} else {
	// 		userEntry.followingRWMu.RLock()
	// 		defer userEntry.followingRWMu.RUnlock()
	// 		databaseError := false
	// 		tempArray := make([]*pb.User, 0, 100)
	// 		for k := range userEntry.user.Following {
	// 			followingEntry, err := userRepo.storage.getUserEntry(k)
	// 			if err != nil {
	// 				// if we have an error here, it means our following data structure has an entry inconsistent
	// 				// with our user structure
	// 				databaseError = true
	// 				errorchan <- errors.New("database corruption")
	// 				panic("database corruption")
	// 			}
	// 			tempArray = append(tempArray, followingEntry.user)
	// 		}
	// 		if !databaseError {
	// 			result <- tempArray
	// 		}
	// 	}
	// }()

	// select {
	// case user := <-result:
	// 	return user, nil
	// case err := <-errorchan:
	// 	return nil, err
	// case <-ctx.Done():
	// 	return nil, ctx.Err()
	// }
	return nil, nil
}

// GetNotFollowing returns an array of users that the given user is not following
func (userRepo *userRepository) GetNotFollowing(ctx context.Context, userId uint64) ([]*pb.User, error) {

	// result := make(chan []*pb.User, 1)
	// errorchan := make(chan error, 1)

	// go func() {
	// 	// Get the user object from the users map
	// 	userEntry, err := userRepo.storage.getUserEntry(userId)
	// 	if err != nil {
	// 		errorchan <- err
	// 	} else {
	// 		userEntry.followingRWMu.RLock()
	// 		defer userEntry.followingRWMu.RUnlock()

	// 		tempArray := make([]*pb.User, 0, 100)

	// 		// Iterate through entire user list
	// 		userRepo.storage.usersRWMu.RLock()
	// 		defer userRepo.storage.usersRWMu.RUnlock()
	// 		for k, v := range userRepo.storage.users {
	// 			// check if user k exists in the user's following list. If not, add it to our
	// 			// temp array
	// 			_, exists := userEntry.user.Following[k]
	// 			if !exists && k != userId {
	// 				tempArray = append(tempArray, v.user)
	// 			}
	// 		}

	// 		result <- tempArray
	// 	}
	// }()

	// select {
	// case user := <-result:
	// 	return user, nil
	// case err := <-errorchan:
	// 	return nil, err
	// case <-ctx.Done():
	// 	return nil, ctx.Err()
	// }
	return nil, nil
}

// DeleteUser removes a user
func (userRepo *userRepository) DeleteUser(ctx context.Context, userID uint64) error {
	// result := make(chan error, 1)
	// buffer := make(chan *userpb.User, 1)

	// go func() {
	// 	resp, err := userRepo.storage.Get(ctx, userPrefix+strconv.FormatUint(userID, 10))
	// 	if err != nil {
	// 		result <- err
	// 		return
	// 	}
	// 	_, err = userRepo.storage.Delete(ctx, userPrefix+strconv.FormatUint(userID, 10))
	// 	if err != nil {
	// 		result <- err
	// 		return
	// 	}
	// 	buffer <- resp.Kvs[0].Value
	// 	result <- nil
	// }()

	// select {
	// case ret := <-result:
	// 	return ret
	// case <-ctx.Done():
	// 	// if ctx done, need to continue to listen to know whether to add userEntry back into db
	// 	go func() {
	// 		select {
	// 		case userEntry := <-buffer:
	// 			userRepo.storage.usersRWMu.Lock()
	// 			defer userRepo.storage.usersRWMu.Unlock()
	// 			userRepo.storage.users[userID] = userEntry
	// 			return
	// 		case <-result:
	// 			// if result != nil, an error occurred and so don't need to add back into db
	// 			if result != nil {
	// 				return
	// 			}
	// 		}

	// 	}()
	// 	return ctx.Err()
	// }
	return nil
}
func (userRepo *userRepository) UpdateUserAccountInfo(ctx context.Context, info *userpb.AccountInformation) error {
	return errors.New("Feature not implemented")
}

func (userRepo *userRepository) getUserId(ctx context.Context) (uint64, error) {
	// i don't know why i need this result channel...should be able to get
	// the response of the get call from the txn response...
	result := make(chan uint64, 1)
	var err error
	getId := func(stm concurrency.STM) error {
		// what happens if get fails? It just never returns, so how do I account for that?
		resp := stm.Get(userIdGenKey)

		// if resp = "", we need to initialize first
		if resp == "" {
			resp = "1"
		}

		id, err := strconv.ParseUint(resp, 10, 64)
		if err != nil {
			result <- uint64(0)
			return err
		}
		result <- id
		stm.Put(userIdGenKey, strconv.FormatUint(id+1, 10))
		return nil
	}
	_, err = concurrency.NewSTM(userRepo.storage, getId)
	if err != nil {
		return 0, err
	}
	return <-result, nil

	// select {
	// case ret := <-result:
	// 	if ret == uint64(0) {
	// 		err = errors.New("could not get user id")
	// 	}
	// 	return ret, err
	// case <-ctx.Done():
	// 	return uint64(0), ctx.Err()
	// }
}

func findRange(array []uint64) ([2]uint64, error) {
	ret := [2]uint64{1, 1}
	for i := 0; i < len(array); i++ {
		if ret[0] > array[i] {
			ret[0] = array[i]
		}
		if ret[1] < array[i] {
			ret[1] = array[i]
		}
	}
	return ret, nil
}

func encodeUser(user *pb.User) (string, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(user); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func decodeUser(userBytes []byte) (*pb.User, error) {
	var user pb.User
	dec := gob.NewDecoder(bytes.NewReader(userBytes))
	if err := dec.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
