package repositories

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/datasource"
	"github.com/jinzhu/gorm"
)

type WinnerRepository interface {
	Contains(winner datamodels.Winner) bool
	Append(winner datamodels.Winner) bool
	GetAll() []datamodels.Winner
}

type winnerSQLRepository struct {
	source *gorm.DB
}

func (r *winnerSQLRepository) Contains(winner datamodels.Winner) bool {
	resWinner := datamodels.Winner{}
	return !r.source.Where("UID=? AND A_ID=?", winner.UID, winner.AID).First(&resWinner).RecordNotFound()
}

func (r *winnerSQLRepository) Append(winner datamodels.Winner) bool {
	if r.Contains(winner) {
		return false
	}
	r.source.Create(winner)
	return true
}

func (r *winnerSQLRepository) GetAll() (winners []datamodels.Winner) {
	r.source.Find(winners)
	return
}

// NewWinnerRepository return a WinnerRepository object
func NewWinnerRepository() WinnerRepository {
	db := datasource.DB
	if !db.HasTable(&datamodels.Winner{}) {
		db.CreateTable(&datamodels.Winner{})
	}
	return &winnerSQLRepository{source: db}
}
