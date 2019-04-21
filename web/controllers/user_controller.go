package controllers

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/repositories"
	"github.com/c-my/lottery_client_server/services"
	"github.com/c-my/lottery_client_server/web/logger"
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
		logger.Warning.Println("get users failed:", err)
	}
	return res
}

func (c *UserController) Append(u datamodels.User) {
	c.Service.Add(u)
}

func (c *UserController) RandomlyGet() datamodels.User {
	return c.Service.GetRandomly()
}

func (c *UserController) RandomlyGetAll() (users []datamodels.User) {
	return c.Service.GetAllRandomly()
}
