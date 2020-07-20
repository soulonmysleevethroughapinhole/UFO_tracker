package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/soulonmysleevethroughapinhole/UFO_tracker/domain/comments"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/services"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/http_utils"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

var (
	CommController commControllerInterface = &commController{}
)

type commControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	//Get(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type commController struct{}

func (c *commController) Create(w http.ResponseWriter, r *http.Request) {
	//auth
	username := "TruthSeeker7"

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.RespondError(w, respErr)
		return
	}
	defer r.Body.Close()

	var commRequest comments.Comment
	if err := json.Unmarshal(reqBody, &commRequest); err != nil {
		respErr := rest_errors.NewBadRequestError("invalid item json body")
		http_utils.RespondError(w, respErr)
		return
	}
	commRequest.Username = username

	res, createErr := services.CommService.CreateComment(commRequest)
	if createErr != nil {
		http_utils.RespondError(w, createErr)
		return
	}

	http_utils.RespondJson(w, http.StatusCreated, res)
}

func (c *commController) GetAll(w http.ResponseWriter, r *http.Request) {
	comms, getErr := services.CommService.ReadCommentsThread(1)
	if getErr != nil {
		http_utils.RespondError(w, getErr)
		return
	}
	http_utils.RespondJson(w, http.StatusOK, comms)
}
