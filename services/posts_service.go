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
	return nil, nil
}

func (s *postService) ReadPost(ID int64) (*posts.Post, rest_errors.RestErr) {
	return nil, nil
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
