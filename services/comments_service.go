package services

import (
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/domain/comments"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

var (
	CommService commServiceInterface = &commService{}
)

type commServiceInterface interface {
	CreateComment(comment comments.Comment) (*comments.Comment, rest_errors.RestErr)
	ReadComment(ID int64) (*comments.Comment, rest_errors.RestErr)
	ReadCommentsThread(ID int64, level string) ([]comments.Comment, rest_errors.RestErr)
}

type commService struct{}

func (s *commService) CreateComment(comm comments.Comment) (*comments.Comment, rest_errors.RestErr) {
	//validate
	if err := comm.Create(); err != nil {
		return nil, err
	}

	return &comm, nil
}

func (s *commService) ReadComment(ID int64) (*comments.Comment, rest_errors.RestErr) {
	comm := comments.Comment{}

	if err := comm.Read(ID); err != nil {
		return nil, err
	}
	return &comm, nil
}

func (s *commService) ReadCommentsThread(ID int64, level string) ([]comments.Comment, rest_errors.RestErr) {
	dao := &comments.Comment{}

	return dao.ReadAll(ID, level)
}
