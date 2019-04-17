package services

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/repositories"
)

type ActivityService interface {
	GetAll() []datamodels.Activity
	Add(activity datamodels.Activity)
}

type activityService struct {
	repo repositories.ActivityRepository
}

func (s *activityService) GetAll() []datamodels.Activity {
	return s.repo.SelectAll()
}

func (s *activityService) Add(activity datamodels.Activity) {
	s.repo.Append(activity)
}

func NewActivityService(repository repositories.ActivityRepository) ActivityService {
	return &activityService{repo: repository}
}
