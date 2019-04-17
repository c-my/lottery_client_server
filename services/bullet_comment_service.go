package services

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/repositories"
)

type BulletCommentService interface {
	GetAll() []datamodels.BulletComment
	Add(comment datamodels.BulletComment)
}

type bulletCommentService struct {
	repo repositories.BulletCommentRepository
}

func (s *bulletCommentService) GetAll() []datamodels.BulletComment {
	return s.repo.SelectAll()
}

func (s *bulletCommentService) Add(danmu datamodels.BulletComment) {
	s.repo.Append(danmu)
}

func NewBulletCommentServece(repository repositories.BulletCommentRepository) BulletCommentService {
	return &bulletCommentService{repo: repository}
}
