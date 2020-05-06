package etcd

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user"
	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/user/userpb"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

const userPrefix string = "/User/"
const followingPrefix string = "/Following"
const emailPrefix string = "/Email"
const accountInfoPrefix string = "/AccountInfo"
const userIdGenKey string = "/NextUserId"

type userRepository struct {
	storage *clientv3.Client
}

// GetUserRepository returns a UserRepository that uses package level storage
func GetUserRepository() user.UserRepository {
	return &userRepository{Client}
}

// NewUserRepository retunrs a UserRepository that uses the given storage
func NewUserRepository(storage *clientv3.Client) user.UserRepository {
	return &userRepository{storage}
}

// CreateUser adds a user to the appropriate data structures
func (userRepo *userRepository) CreateUser(ctx context.Context, user *userpb.User) (uint64, error) {
	result := make(chan uint64, 1)
	errorchan := make(chan error, 1)

	// go func() {
	// 	id, err := userRepo.getUserId(ctx)
	// 	if err != nil {
	// 		errorchan <- err
	// 		return
	// 	}
	// 	info.UserId = id

	// 	var buf bytes.Buffer
	// 	if err := gob.NewEncoder(&buf).Encode(*info); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	// We'll check if the user already exists, then do a write if it does not
	// 	createUser := func(stm concurrency.STM) error {
	// 		if _, err := userRepo.storage.Get(ctx, userPrefix+strconv.FormatUint(info.UserId, 10), buf.String()); err != nil {
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
		// userId, err := userRepo.getUserId(ctx)
		// if err != nil {
		// 	errorchan <- err
		// 	return
		// }
		// user.AccountInformation.UserId = userId

		userEncoded, err := encodeUser(user)
		if err != nil {
			errorchan <- err
			return
		}
		_, err = userRepo.storage.Put(ctx, userPrefix+strconv.FormatUint(user.AccountInformation.UserId, 10), userEncoded)
		if err != nil {
			errorchan <- err
			return
		}
		result <- user.AccountInformation.UserId
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
		resp, err := userRepo.storage.Get(ctx, userPrefix+strconv.FormatUint(userID, 10))
		if err != nil {
			errorchan <- err
			return
		}
		if resp.Kvs == nil {
			errorchan <- errors.New("user not found")
			return
		}
		user, err := decodeUser(resp.Kvs[0].Value)
		if err != nil {
			errorchan <- err
			return
		}
		result <- user
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
		resp, err := userRepo.storage.Get(ctx, userPrefix+extent[0], clientv3.WithRange(userPrefix+extent[1]+"\x00"))
		cp := make([]*pb.User, 0, len(userIDs))
		userIdLookup := make(map[uint64]struct{})
		for _, v := range userIDs {
			userIdLookup[v] = struct{}{}
		}
		for _, v := range resp.Kvs {
			user, err := decodeUser(v.Value)
			if err != nil {
				errorchan <- err
				return
			}
			if _, exists := userIdLookup[user.AccountInformation.UserId]; exists {
				cp = append(cp, user)
			}
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
		resp, err := userRepo.storage.Get(ctx, userPrefix, clientv3.WithPrefix())
		if err != nil {
			errorchan <- err
			return
		}
		cp := make([]*pb.User, 0, 100)
		for _, v := range resp.Kvs {
			user, err := decodeUser(v.Value)
			if err != nil {
				errorchan <- err
				return
			}
			cp = append(cp, user)
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
func (userRepo *userRepository) FollowUser(ctx context.Context, followerId uint64, followedId uint64) error {
	result := make(chan error, 1)
	go func() {
		followUser := func(stm concurrency.STM) error {
			// Get Values
			followerK := userPrefix + strconv.FormatUint(followerId, 10)
			followedK := userPrefix + strconv.FormatUint(followedId, 10)
			followerV, followedV := stm.Get(followerK), stm.Get(followedK)

			if followerV == "" {
				return fmt.Errorf("following user %d not found", followerId)
			}
			if followedK == "" {
				return fmt.Errorf("followed user %d not found", followedId)
			}
			// Decode
			follower, err := decodeUser([]byte(followerV))
			if err != nil {
				return err
			}
			followed, err := decodeUser([]byte(followedV))
			if err != nil {
				return err
			}

			// Modify
			follower.Following[followedId] = followedId
			followed.Followers[followerId] = followerId

			// Encode
			followerEncoded, err := encodeUser(follower)
			if err != nil {
				return err
			}
			followedEncoded, err := encodeUser(followed)
			if err != nil {
				return err
			}

			// Put
			stm.Put(followerK, followerEncoded)
			stm.Put(followedK, followedEncoded)
			return nil
		}
		if _, err := concurrency.NewSTM(userRepo.storage, followUser); err != nil {
			for i := 0; i < 20; i++ {
				fmt.Println("Retrying")
				if _, err := concurrency.NewSTM(userRepo.storage, followUser); err == nil {
					break
				}
			}
			result <- err
			return
		}
		result <- nil
	}()

	select {
	case res := <-result:
		return res
	case <-ctx.Done():
		// listen to the result channel in case the operation was successful, then unfollow
		go func() {
			res := <-result
			if res == nil {
				userRepo.UnFollowUser(context.Background(), followerId, followedId)
			}
		}()
		return ctx.Err()
	}
}

// UnFollowUser updates the following user's following map, and the followed user's followers map
// to reflect that a user has unfollowed another user
func (userRepo *userRepository) UnFollowUser(ctx context.Context, followerId uint64, followedId uint64) error {
	result := make(chan error, 1)

	go func() {
		unfollowUser := func(stm concurrency.STM) error {
			// Get Values
			followerK := userPrefix + strconv.FormatUint(followerId, 10)
			followedK := userPrefix + strconv.FormatUint(followedId, 10)
			followerV, followedV := stm.Get(followerK), stm.Get(followedK)

			// Decode
			follower, err := decodeUser([]byte(followerV))
			if err != nil {
				return err
			}
			followed, err := decodeUser([]byte(followedV))
			if err != nil {
				return err
			}

			// Modify
			delete(follower.Following, followedId)
			delete(followed.Followers, followerId)

			// Encode
			followerEncoded, err := encodeUser(follower)
			if err != nil {
				return err
			}
			followedEncoded, err := encodeUser(followed)
			if err != nil {
				return err
			}

			// Put
			stm.Put(followerK, followerEncoded)
			stm.Put(followedK, followedEncoded)

			return nil
		}
		if _, err := concurrency.NewSTM(userRepo.storage, unfollowUser); err != nil {
			result <- err
			return
		}
		result <- nil
	}()

	select {
	case res := <-result:
		return res
	case <-ctx.Done():
		// listen to the result channel in case the operation was successful, then follow
		go func() {
			res := <-result
			if res == nil {
				userRepo.FollowUser(context.Background(), followerId, followedId)
			}
		}()
		return ctx.Err()
	}
}

// GetUserByUsername returns a user object by their username
func (userRepo *userRepository) GetUserByUsername(ctx context.Context, email string) (*pb.User, error) {
	result := make(chan *pb.User, 1)
	errorchan := make(chan error, 1)

	go func() {
		var idx int
		found := false
		users, err := userRepo.GetAllUsers(ctx)
		if err != nil {
			errorchan <- err
			return
		}
		for i, v := range users {
			if v.AccountInformation.Email == email {
				found = true
				idx = i
				break
			}
		}
		if found {
			result <- users[idx]
		} else {
			errorchan <- errors.New("wrong credentials")
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
	result := make(chan []*pb.User, 1)
	errorchan := make(chan error, 1)

	go func() {
		user, err := userRepo.GetUser(ctx, userId)
		if err != nil {
			errorchan <- err
			return
		}
		if len(user.Following) == 0 {
			result <- make([]*pb.User, 0)
			return
		}
		userIdArr := make([]uint64, 0, len(user.Following))
		for k := range user.Following {
			userIdArr = append(userIdArr, k)
		}
		users, err := userRepo.GetUsers(ctx, userIdArr)
		if err != nil {
			errorchan <- err
			return
		}
		result <- users
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
	result := make(chan []*pb.User, 1)
	errorchan := make(chan error, 1)

	go func() {
		user, err := userRepo.GetUser(ctx, userId)
		if err != nil {
			errorchan <- err
			return
		}
		userIdArr := make([]uint64, 0, len(user.Following))
		for k := range user.Following {
			userIdArr = append(userIdArr, k)
		}
		users, err := userRepo.GetAllUsers(ctx)
		if err != nil {
			errorchan <- err
			return
		}
		// remove Users not in user's following list
		filteredUsers := users[:0]
		for _, v := range users {
			if _, exists := user.Following[v.AccountInformation.UserId]; !exists && v.AccountInformation.UserId != userId {
				filteredUsers = append(filteredUsers, v)
			}
		}

		result <- filteredUsers
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

// DeleteUser removes a user
func (userRepo *userRepository) DeleteUser(ctx context.Context, userID uint64) error {
	result := make(chan error, 1)
	buffer := make(chan *userpb.User, 1)

	go func() {
		// Fetch the user to buffer it
		resp, err := userRepo.storage.Get(ctx, userPrefix+strconv.FormatUint(userID, 10))
		if err != nil {
			result <- err
			return
		}

		// Check to make sure user exists. Not sure that this code works properly...
		if resp.Kvs[0].Value[0] == 0 {
			result <- errors.New("user not found")
			return
		}

		// Delete the user
		_, err = userRepo.storage.Delete(ctx, userPrefix+strconv.FormatUint(userID, 10))
		if err != nil {
			result <- err
			return
		}
		user, err := decodeUser(resp.Kvs[0].Value)
		if err != nil {
			result <- err
		}
		buffer <- user
		result <- nil
	}()

	select {
	case ret := <-result:
		return ret
	case <-ctx.Done():
		// if ctx done, need to continue to listen to know whether to add userEntry back into db
		go func() {
			select {
			case user := <-buffer:
				userRepo.CreateUser(context.Background(), user)
				return
			case <-result:
				// if result != nil, an error occurred and so don't need to add back into db
				if result != nil {
					return
				}
			}
		}()
		return ctx.Err()
	}
}
func (userRepo *userRepository) UpdateUserAccountInfo(ctx context.Context, info *userpb.AccountInformation) error {
	return errors.New("Feature not implemented")
}

func (userRepo *userRepository) NextUserId() (uint64, error) {
	// I don't know why the Txn Response is empty... I would prefer to read the value
	// from the transaction response rather than directly updating redId
	var err error
	var retId uint64
	getId := func(stm concurrency.STM) error {
		// what happens if get fails? It just never returns, so how do I account for that?
		resp := stm.Get(userIdGenKey)
		// if resp = "", we need to initialize first
		if resp == "" {
			resp = "1"
		}
		id, err := strconv.ParseUint(resp, 10, 64)
		if err != nil {
			return err
		}
		retId = id
		stm.Put(userIdGenKey, strconv.FormatUint(id+1, 10))
		return nil
	}
	_, err = concurrency.NewSTM(userRepo.storage, getId)
	return retId, err

}

func convertToStrings(arr []uint64) ([]string, error) {
	retArr := make([]string, len(arr))
	for i, v := range arr {
		retArr[i] = strconv.FormatUint(v, 10)
	}
	return retArr, nil
}

func findRange(array []uint64) ([2]string, error) {
	// ret := [2]uint64{math.MaxUint64, 0}
	// for i := 0; i < len(array); i++ {
	// 	if ret[0] > array[i] {
	// 		ret[0] = array[i]
	// 	}
	// 	if ret[1] < array[i] {
	// 		ret[1] = array[i]
	// 	}
	// }
	// return ret, nil
	var ret [2]string
	arr, err := convertToStrings(array)
	if err != nil {
		return ret, err
	}
	sort.Strings(arr)
	ret[0] = arr[0]
	ret[1] = arr[len(arr)-1]
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
