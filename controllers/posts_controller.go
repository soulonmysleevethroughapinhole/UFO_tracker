package controllers

import "net/http"

var (
	PostsController postsControllerInterface = &postsController{}
)

type postsControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)

	Search(w http.ResponseWriter, r *http.Request)
}

type postsController struct{}

func (c *postsController) Create(w http.ResponseWriter, r *http.Request) {}

func (c *postsController) Get(w http.ResponseWriter, r *http.Request) {}

func (c *postsController) Update(w http.ResponseWriter, r *http.Request) {}

func (c *postsController) Delete(w http.ResponseWriter, r *http.Request) {}

func (c *postsController) Search(w http.ResponseWriter, r *http.Request) {}

