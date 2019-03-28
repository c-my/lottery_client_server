package services

import "github.com/c-my/lottery_client_server/datamodels"

type BullenCommentService interface {
	GetAll() []datamodels.BulletComment
}
