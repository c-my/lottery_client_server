package controllers

import (
	"encoding/json"
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/repositories"
	"github.com/c-my/lottery_client_server/services"
	"github.com/c-my/lottery_client_server/web/logger"
)

var DanmuControl = BulletCommentController{
	Servece: services.NewBulletCommentServece(repositories.NewBulletCommentRepository()),
}

type BulletCommentController struct {
	Servece services.BulletCommentService
}

func (c *BulletCommentController) GetAll() []byte {
	res, err := json.Marshal(c.Servece.GetAll())
	if err != nil {
		logger.Warning.Println("get danmus failed:", err)
		return []byte{}
	}
	return res
}

func (c *BulletCommentController) Append(danmu datamodels.BulletComment) {
	c.Servece.Add(danmu)
}
