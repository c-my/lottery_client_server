package services

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/repositories"
)

// AwardService handles some of the CRUID operations of the award datamodel
type AwardService interface {
	GetAll() []datamodels.Award
}

// NewAwardService returns a new AwardService object
func NewAwardService(repository repositories.AwardRepository) AwardService {
	return &awardService{repo: repository}
}

type awardService struct {
	repo repositories.AwardRepository
}

func (s *awardService) GetAll() []datamodels.Award {
	return s.repo.SelectAll()
}
