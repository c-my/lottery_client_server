package services

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/repositories"
)

var WinnerServicer = winnerService{
	repo: repositories.NewWinnerRepository(),
}

type WinnerService interface {
	AlreadyWin(uid, aid string) bool
	AddWinner(uid, aid string) bool
	GetAllWinners() []datamodels.Winner
}

type winnerService struct {
	repo repositories.WinnerRepository
}

func (s *winnerService) AlreadyWin(uid, aid string) bool {
	return s.repo.Contains(datamodels.Winner{UID: uid, AID: aid})
}

func (s *winnerService) AddWinner(uid, aid string) bool {
	return s.repo.Append(datamodels.Winner{UID: uid, AID: aid})
}

func (s *winnerService) GetAllWinners() []datamodels.Winner {
	return s.repo.GetAll()
}

// NewWinnerService returns a winnerService object
func NewWinnerService(repository repositories.WinnerRepository) WinnerService {
	return &winnerService{repo: repository}
}
