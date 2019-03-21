package controllers

import (
	"github.com/c-my/lottery_client_server/services"

	"github.com/c-my/lottery_client_server/datamodels"
)

// UserController is our user controller
type UserController struct {
	Service services.UserService
}

// Get returns list of users
func (c *UserController) Get() (results []datamodels.User) {
	return c.Service.GetAll()
}
