package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/domain/posts"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/services"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/http_utils"
	"github.com/soulonmysleevethroughapinhole/UFO_tracker/utils/rest_errors"
)

var (
	PostsController postsControllerInterface = &postsController{}
)

type postsControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)

	GetAll(w http.ResponseWriter, r *http.Request)
}

type postsController struct{}

func (c *postsController) Create(w http.ResponseWriter, r *http.Request) {
	//auth
	username := "TruthSeeker7"

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.RespondError(w, respErr)
		return
	}
	defer r.Body.Close()

	var postRequest posts.Post
	if err := json.Unmarshal(reqBody, &postRequest); err != nil {
		respErr := rest_errors.NewBadRequestError("invalid item json body")
		http_utils.RespondError(w, respErr)
		return
	}
	postRequest.Username = username

	res, createErr := services.PostService.CreatePost(postRequest)
	if createErr != nil {
		http_utils.RespondError(w, createErr)
		return
	}

	http_utils.RespondJson(w, http.StatusCreated, res)
}

func (c *postsController) Get(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respErr := rest_errors.NewBadRequestError("invalid request body")
		http_utils.RespondError(w, respErr)
		return
	}

	post, getErr := services.PostService.ReadPost(postID)
	if getErr != nil {
		http_utils.RespondError(w, getErr)
		return
	}
	http_utils.RespondJson(w, http.StatusOK, post)
}

func (c *postsController) Update(w http.ResponseWriter, r *http.Request) {}

func (c *postsController) Delete(w http.ResponseWriter, r *http.Request) {}

func (c *postsController) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, getErr := services.PostService.ReadAllPosts()
	if getErr != nil {
		http_utils.RespondError(w, getErr)
		return
	}

	http_utils.RespondJson(w, http.StatusOK, posts)
}
