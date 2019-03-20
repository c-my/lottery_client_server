package controllers

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_iris/repositories"
)

type UserController struct {
	service repositories.UserRepository
}

func (c *UserController) Get(results []datamodels.User) {
	return c.service.SellectAll()
}
