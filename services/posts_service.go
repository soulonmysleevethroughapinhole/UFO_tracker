package services

import (
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/domain/posts"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

var (
	PostService postServiceInterface = &postService{}
)

type postServiceInterface interface {
	CreatePost(post posts.Post) (*posts.Post, rest_errors.RestErr)
	ReadPost(ID int64) (*posts.Post, rest_errors.RestErr)
	ReadAllPosts() ([]posts.Post, rest_errors.RestErr)
}

type postService struct{}

func (s *postService) CreatePost(post posts.Post) (*posts.Post, rest_errors.RestErr) {
	// if err := post.Validate(); err != nil {
	// 	return nil, err
	// }

	if err := post.Create(); err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *postService) ReadPost(ID int64) (*posts.Post, rest_errors.RestErr) {
	post := posts.Post{ID: ID}

	if err := post.Read(); err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *postService) UpdatePost(post posts.Post) (*posts.Post, rest_errors.RestErr) {
	return nil, nil
}

func (s *postService) DeletePost(post posts.Post) (*posts.Post, rest_errors.RestErr) {
	return nil, nil
}

func (s *postService) ReadAllPosts() ([]posts.Post, rest_errors.RestErr) {
	dao := &posts.Post{}

	return dao.ReadAll()
}
