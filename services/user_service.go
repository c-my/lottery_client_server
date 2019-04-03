package services

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/repositories"
)

// UserService handles some CRUID operations
type UserService interface {
	GetAll() []datamodels.User
	GetByID(id uint) datamodels.User
	GetRandomly() datamodels.User
	Add(datamodels.User)
}

type userService struct {
	repo repositories.UserRepository
}

// NewUserService returns a UserService object
func NewUserService(repository repositories.UserRepository) UserService {
	return &userService{repo: repository}
}

func (s *userService) GetAll() []datamodels.User {
	return s.repo.SelectAll()
}

func (s *userService) GetByID(uid uint) datamodels.User {
	user, found := s.repo.SelectByID(uid)
	if found {
		return user
	}
	return datamodels.User{}
}

func (s *userService) GetRandomly() datamodels.User {
	return s.repo.RandomSelect()
}

func (s *userService) Add(user datamodels.User) {
	s.repo.Append(user)
}
