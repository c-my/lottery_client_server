package repositories

import (
	"github.com/c-my/lottery_client_server/datamodels"
	"github.com/c-my/lottery_client_server/datasource"
	"github.com/jinzhu/gorm"
)

// AwardRepository handles basic operations of award
type AwardRepository interface {
	SelectAll() []datamodels.Award
}

type awardSQLRepository struct {
	source *gorm.DB
}

func (r *awardSQLRepository) SelectAll() (users []datamodels.Award) {
	r.source.Find(&users)
	return
}

// NewAwardSQLRepository returns a award repository
func NewAwardSQLRepository() AwardRepository {
	db := datasource.DB
	if !db.HasTable(&datamodels.Award{}) {
		db.CreateTable(&datamodels.Award{})
	}
	return &awardSQLRepository{source: db}
}
