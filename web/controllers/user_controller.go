package controllers

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_iris/repositories"
)

// UserController is our user controller
type UserController struct {
	service repositories.UserRepository
}

// Get returns list of users
func (c *UserController) Get(results []datamodels.User) {
	return c.service.SellectAll()
}
