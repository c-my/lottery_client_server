package controllers

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/services"
)

// AwardController is our award controller
type AwardController struct {
	Service services.AwardService
}

// Get returns list of awards
func (c *AwardController) Get() (result []datamodels.Award) {
	return c.Service.GetAll()
}
