package controllers

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/repositories"
	"github.com/c-my/lottery_client_server/services"
	"github.com/c-my/lottery_client_server/web/logger"
)

var ActivityControl = ActivityController{
	Service: services.NewActivityService(repositories.NewActivityRepository()),
}

type ActivityController struct {
	Service services.ActivityService
}

func (c *ActivityController) GetAll() []byte {
	res, err := json.Marshal(c.Service.GetAll())
	if err != nil {
		logger.Warning.Println("get activities failed:", err)
	}
	return res
}
