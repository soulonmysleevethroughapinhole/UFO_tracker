package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/soulonmysleevethroughapinhole/UFO_tracker/domain/tags"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/services"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/http_utils"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

var (
	TagController tagControllerInterface = &tagController{}
)

type tagControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetTagsPost(w http.ResponseWriter, r *http.Request)
}

type tagController struct{}

func (c *tagController) Create(w http.ResponseWriter, r *http.Request) {
	username := "TruthSeeker7"

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.RespondError(w, respErr)
		return
	}
	defer r.Body.Close()

	var tagRequest tags.TagInd
	if err := json.Unmarshal(reqBody, &tagRequest); err != nil {
		respErr := rest_errors.NewBadRequestError("inbalid item json body")
		http_utils.RespondError(w, respErr)
		return
	}
	tagRequest.Username = username

	log.Println(tagRequest.TagContent)

	res, createErr := services.TagService.CreateTag(tagRequest)
	if createErr != nil {
		http_utils.RespondError(w, createErr)
		return
	}
	http_utils.RespondJson(w, http.StatusCreated, res)
}

func (c *tagController) GetTagsPost(w http.ResponseWriter, r *http.Request) {
	username := "TruthSeeker7"

	tags, getErr := services.TagService.ReadTagsPost(1, username)
	if getErr != nil {
		http_utils.RespondError(w, getErr)
		return
	}
	http_utils.RespondJson(w, http.StatusOK, tags)
}
