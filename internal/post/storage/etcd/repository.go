package etcd

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post"
	pb "github.com/Distributed-Systems-CSGY9223/yjs310-shs572-dfc296-final-project/internal/post/postpb"
	"github.com/golang/protobuf/ptypes"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
)

const postPrefix string = "Post/"
const postIdGenKey string = "NextPostId"

type postRepository struct {
	storage *clientv3.Client
}

// GetPostRepository returns a UserRepository that uses package level storage
func GetPostRepository(client *clientv3.Client) post.PostRepository {
	return &postRepository{client}
}

// NewPostRepository reutnrs a UserRepository that uses the given storage
func NewPostRepository(storage *clientv3.Client) post.PostRepository {
	return &postRepository{storage}
}

// CreatePost inserts a post into our post map
func (postRepo *postRepository) CreatePost(ctx context.Context, p *pb.Post) (uint64, error) {
	result := make(chan uint64, 1)
	errorchan := make(chan error, 1)
	go func() {

		post := new(pb.Post)
		post.PostID, _ = postRepo.getPostId(ctx)
		post.Timestamp, _ = ptypes.TimestampProto(time.Now())
		post.Message = p.Message
		post.UserId = p.UserId

		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(post); err != nil {
			errorchan <- err
			return
		}

		_, err := postRepo.storage.Put(ctx, postPrefix+strconv.FormatUint(post.PostID, 10), buf.String())
		if err != nil {
			errorchan <- err
			return
		}
		result <- post.PostID
	}()

	select {
	case postID := <-result:
		return postID, nil
	case err := <-errorchan:
		//Sending 0 as an invalid postID
		return 0, err
	case <-ctx.Done():
		go func() {
			select {
			case postID := <-result:
				// if ctx.Done(), we need to make sure that if the post has or will be created, it is deleted,
				// so start a new go routine to monitor the result and error channels
				postRepo.DeletePost(context.Background(), postID)
				return
			case <-errorchan:
				return
			}
		}()

		return 0, ctx.Err()
	}
}

func (postRepo *postRepository) getPostId(ctx context.Context) (uint64, error) {

	var err error
	var retId uint64
	getId := func(stm concurrency.STM) error {
		// what happens if get fails? It just never returns, so how do I account for that?
		resp := stm.Get(postIdGenKey)
		// if resp = "", we need to initialize first
		if resp == "" {
			resp = "1"
		}
		id, err := strconv.ParseUint(resp, 10, 64)
		if err != nil {
			return err
		}
		retId = id
		stm.Put(postIdGenKey, strconv.FormatUint(id+1, 10))
		return nil
	}
	_, err = concurrency.NewSTM(postRepo.storage, getId)
	return retId, err
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPost(ctx context.Context, postID uint64) (*pb.Post, error) {
	result := make(chan *pb.Post, 1)
	errorchan := make(chan error, 1)

	go func() {
		var post pb.Post
		resp, err := postRepo.storage.Get(ctx, postPrefix+strconv.FormatUint(postID, 10))
		if err != nil {
			errorchan <- err
			return
		}
		dec := gob.NewDecoder(bytes.NewReader(resp.Kvs[0].Value))
		if err := dec.Decode(&post); err != nil {
			errorchan <- errors.New("Could not decode message")
		} else {
			result <- &post
		}
	}()

	select {
	case post := <-result:
		return post, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPosts(ctx context.Context, postIDs []uint64) ([]*pb.Post, error) {
	result := make(chan []*pb.Post, 1)
	errorchan := make(chan error, 1)

	go func() {
		extent, err := findRange(postIDs)
		if err != nil {
			errorchan <- err
			return
		}
		resp, err := postRepo.storage.Get(ctx, postPrefix+extent[0], clientv3.WithRange(postPrefix+extent[1]+"\x00"))
		if err != nil {
			errorchan <- err
			return
		}
		cp := make([]*pb.Post, 0, len(postIDs))
		for _, v := range resp.Kvs {
			var post pb.Post
			dec := gob.NewDecoder(bytes.NewReader(v.Value))
			if err := dec.Decode(&post); err != nil {
				log.Fatalf("could not decode message (%v)", err)
			}
			cp = append(cp, &post)
		}
		result <- cp
	}()

	select {
	case posts := <-result:
		return posts, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func convertToStrings(arr []uint64) ([]string, error) {
	retArr := make([]string, len(arr))
	for i, v := range arr {
		retArr[i] = strconv.FormatUint(v, 10)
	}
	return retArr, nil
}

func findRange(array []uint64) ([2]string, error) {
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

func (postRepo *postRepository) UpdatePost(ctx context.Context, p pb.Post) error {
	return errors.New("Feature not implemented")
}

func (postRepo *postRepository) DeletePost(ctx context.Context, postID uint64) error {
	errorchan := make(chan error, 1)
	buffer := make(chan *pb.Post, 1)

	go func() {

		// Fetch the user to buffer it
		resp, err := postRepo.storage.Get(ctx, postPrefix+strconv.FormatUint(postID, 10))
		if err != nil {
			errorchan <- err
			return
		}
		if resp.Kvs[0].Value[0] == 0 {
			errorchan <- errors.New("user not found")
			return
		}

		//Delete the post
		_, err = postRepo.storage.Delete(ctx, postPrefix+strconv.FormatUint(postID, 10))
		if err != nil {
			errorchan <- err
			return
		}
		var post pb.Post
		dec := gob.NewDecoder(bytes.NewReader(resp.Kvs[0].Value))
		if err := dec.Decode(&post); err != nil {
			errorchan <- err
			return
		}
		buffer <- &post

	}()

	select {
	case err := <-errorchan:
		return err
	case <-ctx.Done():
		// if ctx done, need to continue to listen to know whether to add postEntry back into db
		go func() {
			select {
			case post := <-buffer:
				postRepo.CreatePost(context.Background(), post)
				return
			case err := <-errorchan:
				// if err != nil, an error occurred and so don't need to add back into db
				if err != nil {
					return
				}
			}

		}()
		return ctx.Err()
	}
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPostsByAuthor(ctx context.Context, userIDs []uint64) ([]*pb.Post, error) {
	result := make(chan []*pb.Post, 1)
	errorchan := make(chan error, 1)

	go func() {

		resp, err := postRepo.storage.Get(ctx, postPrefix, clientv3.WithPrefix())
		if err != nil {
			errorchan <- err
			return
		}
		postArr := make([]*pb.Post, 0, len(userIDs)*100)

		for _, v := range resp.Kvs {
			var post pb.Post
			dec := gob.NewDecoder(bytes.NewReader(v.Value))
			if err := dec.Decode(&post); err != nil {
				errorchan <- errors.New("could not decode message")
				return
			}
			for _, u := range userIDs {
				if post.UserId == u {
					postArr = append(postArr, &post)
					break
				}
			}

		}
		result <- postArr
	}()

	select {
	case posts := <-result:
		return posts, nil
	case err := <-errorchan:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
