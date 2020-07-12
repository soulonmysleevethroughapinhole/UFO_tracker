package controllers

import "net/http"

var (
	UsersController usersControllerInterface = &usersController{}
)

type usersControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	//Update(w http.ResponseWriter, r *http.Request)
	//Delete(w http.ResponseWriter, r *http.Request)
}

type usersController struct{}

func (c *usersController) Create(w http.ResponseWriter, r *http.Request) {

}

func (c *usersController) Get(w http.ResponseWriter, r *http.Request) {

}
