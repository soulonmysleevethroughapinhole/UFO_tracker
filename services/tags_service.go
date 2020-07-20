package services

import (
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/domain/tags"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

var (
	TagService tagServiceInterface = &tagService{}
)

type tagServiceInterface interface {
	CreateTag(tags.TagInd) (*tags.TagInd, rest_errors.RestErr)
	ReadTagsPost(ID int64, username string) ([]tags.Tag, rest_errors.RestErr)
}

type tagService struct{}

func (s *tagService) CreateTag(tag tags.TagInd) (*tags.TagInd, rest_errors.RestErr) {
	//validate

	if err := tag.Create(); err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *tagService) ReadTagsPost(ID int64, username string) ([]tags.Tag, rest_errors.RestErr) {
	dao := &tags.Tag{}
	return dao.ReadTagsPost(ID, username)
}
