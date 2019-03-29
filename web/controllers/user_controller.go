package controllers

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/repositories"
	"github.com/c-my/lottery_client_server/services"
	"github.com/labstack/gommon/log"
)

// UserController is our user controller

var UserControl = UserController{
	Service: services.NewUserService(repositories.NewUserRepository()),
}

type UserController struct {
	Service services.UserService
}

// Get returns list of users
func (c *UserController) Get() []byte {
	res, err := json.Marshal(c.Service.GetAll())
	if err != nil {
		log.Fatal("get users failed: ", err)
	}
	return res
}
