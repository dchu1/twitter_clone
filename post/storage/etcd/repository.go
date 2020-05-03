package etcd

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"log"
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
func GetPostRepository() post.PostRepository {
	return &postRepository{Client}
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
			log.Fatal(err)
		}

		_, err := postRepo.storage.Put(ctx, postPrefix+strconv.FormatUint(post.PostID, 10), buf.String())
		if err != nil {
			log.Fatal(err)
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
			case <-result:
				postRepo.DeletePost(context.Background(), p.PostID)
				return
			case <-errorchan:
				return
			}
		}()

		return 0, ctx.Err()
	}
}

func (postRepo *postRepository) getPostId(ctx context.Context) (uint64, error) {

	result := make(chan uint64, 1)
	var err error
	getId := func(stm concurrency.STM) error {
		// what happens if get fails? It just never returns, so how do I account for that?
		resp := stm.Get(postIdGenKey)

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
		stm.Put(postIdGenKey, strconv.FormatUint(id+1, 10))
		return nil
	}
	_, err = concurrency.NewSTM(postRepo.storage, getId)
	if err != nil {
		return 0, err
	}
	return <-result, nil
}

// GetPosts retrieves an array of post from the post map
func (postRepo *postRepository) GetPost(ctx context.Context, postID uint64) (*pb.Post, error) {
	result := make(chan *pb.Post, 1)
	errorchan := make(chan error, 1)

	go func() {
		var post pb.Post
		resp, err := postRepo.storage.Get(ctx, postPrefix+strconv.FormatUint(postID, 10))
		if err != nil {
			log.Fatal(err)
		}
		dec := gob.NewDecoder(bytes.NewReader(resp.Kvs[0].Value))
		if err := dec.Decode(&post); err != nil {
			log.Fatalf("could not decode message (%v)", err)
			errorchan <- errors.New("post not found")
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
		resp, err := postRepo.storage.Get(ctx, postPrefix+strconv.FormatUint(extent[0], 10), clientv3.WithRange(postPrefix+strconv.FormatUint(extent[1]+1, 10)))
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

func (postRepo *postRepository) UpdatePost(ctx context.Context, p pb.Post) error {
	return errors.New("Feature not implemented")
}

func (postRepo *postRepository) DeletePost(ctx context.Context, postID uint64) error {
	errorchan := make(chan error, 1)
	buffer := make(chan *pb.Post, 1)

	go func() {

		postEntry, err := postRepo.GetPost(ctx, postID)
		if err != nil {
			errorchan <- errors.New("post not exist")
			return
		}
		_, err = postRepo.storage.Delete(ctx, postPrefix+strconv.FormatUint(postID, 10))
		if err != nil {
			log.Fatal(err)
			errorchan <- err
		} else {
			buffer <- postEntry
		}

	}()

	select {
	case err := <-errorchan:
		return err
	case <-ctx.Done():
		// if ctx done, need to continue to listen to know whether to add postEntry back into db
		go func() {
			select {
			case err := <-errorchan:
				// if result != nil, an error occurred and so don't need to add back into db
				if err != nil {
					return
				}
				postEntry := <-buffer
				var buf bytes.Buffer
				if err := gob.NewEncoder(&buf).Encode(postEntry); err != nil {
					log.Fatal(err)
				}

				_, err = postRepo.storage.Put(ctx, postPrefix+strconv.FormatUint(postEntry.PostID, 10), buf.String())
				if err != nil {
					log.Fatal(err)
				}
				return
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
			log.Fatal(err)
			errorchan <- err
			return
		}
		postArr := make([]*pb.Post, 0, len(userIDs)*100)

		for _, v := range resp.Kvs {
			var post pb.Post
			dec := gob.NewDecoder(bytes.NewReader(v.Value))
			if err := dec.Decode(&post); err != nil {
				log.Fatalf("could not decode message (%v)", err)
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
